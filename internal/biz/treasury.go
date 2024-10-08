package biz

import (
	"AltTreasury/internal/constants"
	"context"
	"errors"
	"fmt"
	"os"
	"strings"
	"time"

	"github.com/go-kratos/kratos/v2/log"
)

type WithdrawClaim struct {
	ID               int64
	StaffID          int64
	Amount           float64
	TokenAddress     string
	RecipientAddress string
	Status           string
	CreatedAt        time.Time
	UpdatedAt        time.Time
}

type WithdrawClaimConfirmation struct {
	ID              int64
	WithdrawClaimID int64
	ManagerID       int64
	ActionType      string
	ConfirmedAt     time.Time
}

type TreasuryRepo interface {
	CreateWithdrawClaim(ctx context.Context, claim *WithdrawClaim) (int64, error)
	GetWithdrawClaim(ctx context.Context, id int64) (*WithdrawClaim, error)
	UpdateWithdrawClaim(ctx context.Context, claim *WithdrawClaim) error
	CreateWithdrawClaimConfirmation(ctx context.Context, confirmation *WithdrawClaimConfirmation) error
	GetWithdrawClaimConfirmations(ctx context.Context, claimID int64) ([]*WithdrawClaimConfirmation, error)
	CheckTreasuryBalance(ctx context.Context, tokenAddress string, amount float64) (bool, error)
	TransferTokens(ctx context.Context, tokenAddress, toAddress string, amount float64) error
	ListWithdrawClaimConfirmations(ctx context.Context, staffID, managerID int64, actionType string, page, pageSize int32) ([]*WithdrawClaimConfirmation, int32, error)
	WithTransaction(ctx context.Context, fn func(repo TreasuryRepo) error) error
	ListWithdrawClaims(ctx context.Context, staffID *int64, status *string, createdAfter, createdBefore *time.Time, page, pageSize int) ([]*WithdrawClaim, int, error)
}

type TreasuryUsecase struct {
	repo TreasuryRepo
	log  *log.Helper
}

func NewTreasuryUsecase(repo TreasuryRepo, logger log.Logger) *TreasuryUsecase {
	return &TreasuryUsecase{repo: repo, log: log.NewHelper(logger)}
}

var (
	ErrDuplicateConfirmation = errors.New("duplicate confirmation: manager has already performed this action on the claim")
)

func (uc *TreasuryUsecase) CreateWithdrawClaim(ctx context.Context, staffID int64, amount float64, recipientAddress string) (int64, error) {
	tokenAddress := os.Getenv("TOKEN_ADDRESS")
	if tokenAddress == "" {
		tokenAddress = constants.DefaultTokenAddress
	}

	claim := &WithdrawClaim{
		StaffID:          staffID,
		Amount:           amount,
		TokenAddress:     tokenAddress,
		RecipientAddress: recipientAddress,
		Status:           "pending",
		CreatedAt:        time.Now(),
		UpdatedAt:        time.Now(),
	}
	return uc.repo.CreateWithdrawClaim(ctx, claim)
}

func (uc *TreasuryUsecase) GetWithdrawClaim(ctx context.Context, id int64) (*WithdrawClaim, error) {
	return uc.repo.GetWithdrawClaim(ctx, id)
}

func (uc *TreasuryUsecase) ApproveWithdrawClaim(ctx context.Context, claimID int64, managerID int64) (bool, string, error) {
	var success bool
	var message string
	// 模拟授权验证
	isValidManager := false
	for _, id := range constants.ManagerIDs {
		if id == managerID {
			isValidManager = true
			break
		}
	}
	if !isValidManager {
		return false, "invalid manager ID", errors.New("invalid manager ID")
	}

	// do all things in transaction
	err := uc.repo.WithTransaction(ctx, func(repo TreasuryRepo) error {

		claim, err := repo.GetWithdrawClaim(ctx, claimID)
		if err != nil {
			return err
		}

		if claim.Status != "pending" {
			return errors.New("claim is not in pending status")
		}

		// 在 approve 之前检查treasury余额
		sufficient, err := repo.CheckTreasuryBalance(ctx, claim.TokenAddress, claim.Amount)
		if err != nil {
			return err
		}
		if !sufficient {
			return errors.New("insufficient treasury balance, cannot approve")
		}

		confirmation := &WithdrawClaimConfirmation{
			WithdrawClaimID: claimID,
			ManagerID:       managerID,
			ActionType:      "approve",
			ConfirmedAt:     time.Now(),
		}
		err = repo.CreateWithdrawClaimConfirmation(ctx, confirmation)
		if err != nil {
			if errors.Is(err, ErrDuplicateConfirmation) {
				return errors.New("manager has already approved this claim")
			}
			return err
		}

		confirmations, err := repo.GetWithdrawClaimConfirmations(ctx, claimID)
		if err != nil {
			return err
		}

		approvalCount := 0
		for _, conf := range confirmations {
			if conf.ActionType == "approve" {
				approvalCount++
			}
		}

		if approvalCount >= 2 {

			err = uc.transferTokensWithRetry(ctx, claim.TokenAddress, claim.RecipientAddress, claim.Amount)
			if err != nil {
				return err
			}

			claim.Status = "executed"
			claim.UpdatedAt = time.Now()
			err = repo.UpdateWithdrawClaim(ctx, claim)
			if err != nil {
				return err
			}
		}

		return nil // 或者返回错误
	})
	if err != nil {
		success = false
		message = err.Error()
	} else {
		success = true
		message = "Claim has been approved successfully"
	}
	return success, message, err
}

func (uc *TreasuryUsecase) RejectWithdrawClaim(ctx context.Context, claimID int64, managerID int64) (bool, string, error) {
	var success bool
	var message string
	err := uc.repo.WithTransaction(ctx, func(repo TreasuryRepo) error {
		claim, err := uc.repo.GetWithdrawClaim(ctx, claimID)
		if err != nil {
			return err
		}

		if claim.Status != "pending" {
			return fmt.Errorf("claim is not in pending status")
		}

		confirmation := &WithdrawClaimConfirmation{
			WithdrawClaimID: claimID,
			ManagerID:       managerID,
			ActionType:      "reject",
			ConfirmedAt:     time.Now(),
		}
		err = uc.repo.CreateWithdrawClaimConfirmation(ctx, confirmation)
		if err != nil {
			if errors.Is(err, ErrDuplicateConfirmation) {
				return fmt.Errorf("manager has already rejected this claim")
			}
			return err
		}

		confirmations, err := uc.repo.GetWithdrawClaimConfirmations(ctx, claimID)
		if err != nil {
			return err
		}

		rejectionCount := 0
		for _, conf := range confirmations {
			if conf.ActionType == "reject" {
				rejectionCount++
			}
		}

		if rejectionCount >= 2 {
			claim.Status = "rejected"
			claim.UpdatedAt = time.Now()
			err = uc.repo.UpdateWithdrawClaim(ctx, claim)
			if err != nil {
				return err
			}
		}

		return nil // 或者返回错误
	})

	if err != nil {
		success = false
		message = fmt.Sprintf("Reject withdrawal claim failed: %v", err)
		return success, message, err
	}

	success = true
	message = "Withdrawal claim rejected successfully"
	return success, message, nil
}

func (uc *TreasuryUsecase) ListWithdrawClaimConfirmations(ctx context.Context, staffID, managerID int64, actionType string, page, pageSize int32) ([]*WithdrawClaimConfirmation, int32, error) {
	return uc.repo.ListWithdrawClaimConfirmations(ctx, staffID, managerID, actionType, page, pageSize)
}

func (uc *TreasuryUsecase) transferTokensWithRetry(ctx context.Context, tokenAddress, recipientAddress string, amount float64) error {
	maxRetries := 3
	baseDelay := time.Second

	for attempt := 0; attempt < maxRetries; attempt++ {
		err := uc.repo.TransferTokens(ctx, tokenAddress, recipientAddress, amount)
		if err == nil {
			return nil
		}

		uc.log.Errorf("Transfer attempt %d failed: %v", attempt+1, err)

		if strings.Contains(err.Error(), "insufficient ETH balance for gas") {
			// 如果是 ETH 余额不足，直接返回错误，不再重试
			return fmt.Errorf("insufficient ETH balance for gas: %v", err)
		}

		if attempt < maxRetries-1 {
			delay := baseDelay * time.Duration(1<<uint(attempt))
			uc.log.Infof("Retrying in %v...", delay)
			time.Sleep(delay)
		}
	}

	return fmt.Errorf("failed to transfer tokens after %d attempts", maxRetries)
}

func (uc *TreasuryUsecase) ListWithdrawClaims(ctx context.Context, staffID *int64, status *string, createdAfter, createdBefore *time.Time, page, pageSize int) ([]*WithdrawClaim, int, error) {
	if page <= 0 {
		page = 1
	}
	if pageSize <= 0 {
		pageSize = 10
	}
	return uc.repo.ListWithdrawClaims(ctx, staffID, status, createdAfter, createdBefore, page, pageSize)
}

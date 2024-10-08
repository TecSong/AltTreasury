package data

import (
	"context"
	"math"
	"time"

	"AltTreasury/internal/biz"
	"AltTreasury/internal/constants"
	"AltTreasury/internal/ethereum"
	"math/big"

	"github.com/go-kratos/kratos/v2/log"
	"github.com/go-sql-driver/mysql"
	"gorm.io/gorm"
)

type WithdrawClaim struct {
	ID               int64 `gorm:"primaryKey"`
	StaffID          int64
	Amount           float64
	TokenAddress     string
	RecipientAddress string
	Status           string
	CreatedAt        time.Time
	UpdatedAt        time.Time
}

type WithdrawClaimConfirmation struct {
	ID              int64 `gorm:"primaryKey"`
	WithdrawClaimID int64
	ManagerID       int64
	ActionType      string
	ConfirmedAt     time.Time
}

type treasuryRepo struct {
	data *Data
	log  *log.Helper
	eth  *ethereum.EthereumClient
}

func NewTreasuryRepo(data *Data, logger log.Logger) biz.TreasuryRepo {
	eth, err := ethereum.NewEthereumClient()
	if err != nil {
		log.NewHelper(logger).Errorf("Failed to create Ethereum client: %v", err)
		return nil
	}
	return &treasuryRepo{
		data: data,
		log:  log.NewHelper(logger),
		eth:  eth,
	}
}

func (r *treasuryRepo) CreateWithdrawClaim(ctx context.Context, claim *biz.WithdrawClaim) (int64, error) {
	dbClaim := &WithdrawClaim{
		StaffID:          claim.StaffID,
		Amount:           claim.Amount,
		TokenAddress:     claim.TokenAddress,
		RecipientAddress: claim.RecipientAddress,
		Status:           claim.Status,
		CreatedAt:        claim.CreatedAt,
		UpdatedAt:        claim.UpdatedAt,
	}
	if err := r.data.db.WithContext(ctx).Create(dbClaim).Error; err != nil {
		return 0, err
	}
	return dbClaim.ID, nil
}

func (r *treasuryRepo) GetWithdrawClaim(ctx context.Context, id int64) (*biz.WithdrawClaim, error) {
	var dbClaim WithdrawClaim
	err := r.data.db.WithContext(ctx).First(&dbClaim, id).Error
	if err != nil {
		return nil, err
	}
	return &biz.WithdrawClaim{
		ID:               dbClaim.ID,
		StaffID:          dbClaim.StaffID,
		Amount:           dbClaim.Amount,
		TokenAddress:     dbClaim.TokenAddress,
		RecipientAddress: dbClaim.RecipientAddress,
		Status:           dbClaim.Status,
		CreatedAt:        dbClaim.CreatedAt,
		UpdatedAt:        dbClaim.UpdatedAt,
	}, nil
}

func (r *treasuryRepo) UpdateWithdrawClaim(ctx context.Context, claim *biz.WithdrawClaim) error {
	return r.data.db.WithContext(ctx).Model(&WithdrawClaim{}).Where("id = ?", claim.ID).Updates(map[string]interface{}{
		"status":     claim.Status,
		"updated_at": claim.UpdatedAt,
	}).Error
}

func (r *treasuryRepo) CreateWithdrawClaimConfirmation(ctx context.Context, confirmation *biz.WithdrawClaimConfirmation) error {
	dbConfirmation := &WithdrawClaimConfirmation{
		WithdrawClaimID: confirmation.WithdrawClaimID,
		ManagerID:       confirmation.ManagerID,
		ActionType:      confirmation.ActionType,
		ConfirmedAt:     confirmation.ConfirmedAt,
	}
	err := r.data.db.WithContext(ctx).Create(dbConfirmation).Error
	if err != nil {
		if mysqlErr, ok := err.(*mysql.MySQLError); ok && mysqlErr.Number == 1062 {
			return biz.ErrDuplicateConfirmation
		}
		return err
	}
	return nil
}

func (r *treasuryRepo) GetWithdrawClaimConfirmations(ctx context.Context, claimID int64) ([]*biz.WithdrawClaimConfirmation, error) {
	var dbConfirmations []WithdrawClaimConfirmation
	err := r.data.db.WithContext(ctx).Where("withdraw_claim_id = ?", claimID).Find(&dbConfirmations).Error
	if err != nil {
		return nil, err
	}

	confirmations := make([]*biz.WithdrawClaimConfirmation, len(dbConfirmations))
	for i, dbConf := range dbConfirmations {
		confirmations[i] = &biz.WithdrawClaimConfirmation{
			ID:              dbConf.ID,
			WithdrawClaimID: dbConf.WithdrawClaimID,
			ManagerID:       dbConf.ManagerID,
			ActionType:      dbConf.ActionType,
			ConfirmedAt:     dbConf.ConfirmedAt,
		}
	}
	return confirmations, nil
}

func (r *treasuryRepo) CheckTreasuryBalance(ctx context.Context, tokenAddress string, amount float64) (bool, error) {
	balance, err := r.eth.CheckBalance(ctx, tokenAddress)
	if err != nil {
		return false, err
	}

	// 使用常量中的 TokenDecimals
	decimals := constants.TokenDecimals

	// 将 amount 转换为整数
	amountInt := new(big.Int)
	amountFloat := new(big.Float).SetFloat64(amount)
	decimalsBigFloat := new(big.Float).SetFloat64(math.Pow10(int(decimals)))
	amountFloat.Mul(amountFloat, decimalsBigFloat)
	amountFloat.Int(amountInt)

	// 直接比较整数值
	return balance.Cmp(amountInt) >= 0, nil
}

func (r *treasuryRepo) TransferTokens(ctx context.Context, tokenAddress, toAddress string, amount float64) error {
	// 将 float64 转换为 big.Int，保持精度
	amountInt := new(big.Int)
	amountFloat := new(big.Float).SetFloat64(amount)

	// 使用常量中的 TokenDecimals
	decimals := constants.TokenDecimals

	// 将金额乘以 10^decimals
	multiplier := new(big.Float).SetFloat64(math.Pow10(int(decimals)))
	amountFloat.Mul(amountFloat, multiplier)

	amountFloat.Int(amountInt)

	// 使用转换后的 big.Int 进行转账
	return r.eth.TransferTokens(ctx, tokenAddress, toAddress, amountInt)
}

func (r *treasuryRepo) ListWithdrawClaimConfirmations(ctx context.Context, staffID, managerID int64, actionType string, page, pageSize int32) ([]*biz.WithdrawClaimConfirmation, int32, error) {
	var dbConfirmations []WithdrawClaimConfirmation
	var total int64

	query := r.data.db.WithContext(ctx).Model(&WithdrawClaimConfirmation{}).Joins("JOIN withdrawal_claim ON withdrawal_claim.id = withdraw_claim_confirmation.withdraw_claim_id")

	if staffID != 0 {
		query = query.Where("withdrawal_claim.staff_id = ?", staffID)
	}
	if managerID != 0 {
		query = query.Where("withdraw_claim_confirmation.manager_id = ?", managerID)
	}
	if actionType != "" {
		query = query.Where("withdraw_claim_confirmation.action_type = ?", actionType)
	}

	if err := query.Count(&total).Error; err != nil {
		return nil, 0, err
	}

	offset := (page - 1) * pageSize
	err := query.Order("withdraw_claim_confirmation.confirmed_at DESC").Offset(int(offset)).Limit(int(pageSize)).Find(&dbConfirmations).Error
	if err != nil {
		return nil, 0, err
	}

	confirmations := make([]*biz.WithdrawClaimConfirmation, len(dbConfirmations))
	for i, dbConf := range dbConfirmations {
		confirmations[i] = &biz.WithdrawClaimConfirmation{
			ID:              dbConf.ID,
			WithdrawClaimID: dbConf.WithdrawClaimID,
			ManagerID:       dbConf.ManagerID,
			ActionType:      dbConf.ActionType,
			ConfirmedAt:     dbConf.ConfirmedAt,
		}
	}

	return confirmations, int32(total), nil
}

func (r *treasuryRepo) WithTransaction(ctx context.Context, fn func(repo biz.TreasuryRepo) error) error {
	return r.data.db.WithContext(ctx).Transaction(func(tx *gorm.DB) error {
		txRepo := &treasuryRepo{
			data: &Data{db: tx},
			log:  r.log,
			eth:  r.eth,
		}
		return fn(txRepo)
	})
}

func (r *treasuryRepo) ListWithdrawClaims(ctx context.Context, staffID *int64, status *string, createdAfter, createdBefore *time.Time, page, pageSize int) ([]*biz.WithdrawClaim, int, error) {
	if page <= 0 {
		page = 1
	}
	if pageSize <= 0 {
		pageSize = 10
	}

	var claims []WithdrawClaim
	var total int64

	query := r.data.db.Model(&WithdrawClaim{})

	if staffID != nil {
		query = query.Where("staff_id = ?", *staffID)
	}
	if status != nil {
		query = query.Where("status = ?", *status)
	}
	if createdAfter != nil {
		query = query.Where("created_at >= ?", *createdAfter)
	}
	if createdBefore != nil {
		query = query.Where("created_at <= ?", *createdBefore)
	}

	// Count total records
	if err := query.Count(&total).Error; err != nil {
		return nil, 0, err
	}

	// Get paginated records
	offset := (page - 1) * pageSize
	err := query.Offset(offset).Limit(pageSize).Order("created_at DESC").Find(&claims).Error
	if err != nil {
		return nil, 0, err
	}

	bizClaims := make([]*biz.WithdrawClaim, len(claims))
	for i, claim := range claims {
		bizClaims[i] = &biz.WithdrawClaim{
			ID:               claim.ID,
			StaffID:          claim.StaffID,
			Amount:           claim.Amount,
			TokenAddress:     claim.TokenAddress,
			RecipientAddress: claim.RecipientAddress,
			Status:           claim.Status,
			CreatedAt:        claim.CreatedAt,
			UpdatedAt:        claim.UpdatedAt,
		}
	}

	return bizClaims, int(total), nil
}

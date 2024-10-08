package service

import (
	"context"
	"time"

	pb "AltTreasury/api/treasury/v1"
	"AltTreasury/internal/biz"

	"github.com/go-kratos/kratos/v2/log"
	"google.golang.org/protobuf/types/known/timestamppb"
)

type TreasuryService struct {
	pb.UnimplementedTreasuryServer

	uc  *biz.TreasuryUsecase
	log *log.Helper
}

func NewTreasuryService(uc *biz.TreasuryUsecase, logger log.Logger) *TreasuryService {
	return &TreasuryService{uc: uc, log: log.NewHelper(logger)}
}

func (s *TreasuryService) CreateWithdrawClaim(ctx context.Context, req *pb.CreateWithdrawClaimRequest) (*pb.CreateWithdrawClaimReply, error) {
	id, err := s.uc.CreateWithdrawClaim(ctx, req.StaffId, req.Amount, req.RecipientAddress)
	if err != nil {
		return nil, err
	}
	return &pb.CreateWithdrawClaimReply{ClaimId: id}, nil
}

func (s *TreasuryService) GetWithdrawClaim(ctx context.Context, req *pb.GetWithdrawClaimRequest) (*pb.GetWithdrawClaimReply, error) {
	claim, err := s.uc.GetWithdrawClaim(ctx, req.ClaimId)
	if err != nil {
		return nil, err
	}
	return &pb.GetWithdrawClaimReply{
		ClaimId:          claim.ID,
		StaffId:          claim.StaffID,
		Amount:           claim.Amount,
		TokenAddress:     claim.TokenAddress,
		RecipientAddress: claim.RecipientAddress,
		Status:           claim.Status,
		CreatedAt:        timestamppb.New(claim.CreatedAt),
		UpdatedAt:        timestamppb.New(claim.UpdatedAt),
	}, nil
}

func (s *TreasuryService) ApproveWithdrawClaim(ctx context.Context, req *pb.ApproveWithdrawClaimRequest) (*pb.ApproveWithdrawClaimReply, error) {
	success, msg, err := s.uc.ApproveWithdrawClaim(ctx, req.ClaimId, req.ManagerId)
	if err != nil {
		return nil, err
	}
	return &pb.ApproveWithdrawClaimReply{Success: success, Message: msg}, nil
}

func (s *TreasuryService) RejectWithdrawClaim(ctx context.Context, req *pb.RejectWithdrawClaimRequest) (*pb.RejectWithdrawClaimReply, error) {
	success, msg, err := s.uc.RejectWithdrawClaim(ctx, req.ClaimId, req.ManagerId)
	if err != nil {
		return nil, err
	}
	return &pb.RejectWithdrawClaimReply{Success: success, Message: msg}, nil
}

func (s *TreasuryService) ListWithdrawClaimConfirmations(ctx context.Context, req *pb.ListWithdrawClaimConfirmationsRequest) (*pb.ListWithdrawClaimConfirmationsReply, error) {
	staffID := req.GetStaffId()
	managerID := req.GetManagerId()
	actionType := req.GetActionType()
	page := req.GetPage()
	pageSize := req.GetPageSize()

	if page <= 0 {
		page = 1
	}
	if pageSize <= 0 {
		pageSize = 10
	}

	confirmations, total, err := s.uc.ListWithdrawClaimConfirmations(ctx, staffID, managerID, actionType, page, pageSize)
	if err != nil {
		return nil, err
	}

	reply := &pb.ListWithdrawClaimConfirmationsReply{
		Confirmations: make([]*pb.WithdrawClaimConfirmationInfo, len(confirmations)),
		Total:         total,
	}

	for i, conf := range confirmations {
		reply.Confirmations[i] = &pb.WithdrawClaimConfirmationInfo{
			Id:              conf.ID,
			WithdrawClaimId: conf.WithdrawClaimID,
			ManagerId:       conf.ManagerID,
			ActionType:      conf.ActionType,
			ConfirmedAt:     timestamppb.New(conf.ConfirmedAt),
		}
	}

	return reply, nil
}

func (s *TreasuryService) ListWithdrawClaims(ctx context.Context, req *pb.ListWithdrawClaimsRequest) (*pb.ListWithdrawClaimsReply, error) {
	page := int(req.Page)
	pageSize := int(req.PageSize)

	if page <= 0 {
		page = 1
	}
	if pageSize <= 0 {
		pageSize = 10
	}

	var createdAfter, createdBefore *time.Time
	if req.CreatedAfter != nil {
		t := req.CreatedAfter.AsTime()
		createdAfter = &t
	}
	if req.CreatedBefore != nil {
		t := req.CreatedBefore.AsTime()
		createdBefore = &t
	}

	claims, total, err := s.uc.ListWithdrawClaims(ctx, req.StaffId, req.Status, createdAfter, createdBefore, page, pageSize)
	if err != nil {
		return nil, err
	}

	reply := &pb.ListWithdrawClaimsReply{
		Claims:   make([]*pb.WithdrawClaimInfo, len(claims)),
		Total:    int32(total),
		Page:     int32(page),
		PageSize: int32(pageSize),
	}

	for i, claim := range claims {
		reply.Claims[i] = &pb.WithdrawClaimInfo{
			ClaimId:          claim.ID,
			StaffId:          claim.StaffID,
			Amount:           claim.Amount,
			TokenAddress:     claim.TokenAddress,
			RecipientAddress: claim.RecipientAddress,
			Status:           claim.Status,
			CreatedAt:        timestamppb.New(claim.CreatedAt),
			UpdatedAt:        timestamppb.New(claim.UpdatedAt),
		}
	}

	return reply, nil
}

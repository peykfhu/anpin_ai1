package service

import (
	"context"
	"crypto/rand"
	"encoding/hex"
	"fmt"
	"log/slog"
	"math"
	"strings"
	"time"

	dbent "github.com/Wei-Shaw/sub2api/ent"
	"github.com/Wei-Shaw/sub2api/ent/commissionlog"
	"github.com/Wei-Shaw/sub2api/ent/referralrecord"
	"github.com/Wei-Shaw/sub2api/ent/user"
	infraerrors "github.com/Wei-Shaw/sub2api/internal/pkg/errors"
	"github.com/Wei-Shaw/sub2api/internal/pkg/pagination"
)

// --- Errors ---

var (
	ErrReferralCodeNotFound  = infraerrors.NotFound("REFERRAL_CODE_NOT_FOUND", "referral code not found")
	ErrReferralSelfRefer     = infraerrors.BadRequest("REFERRAL_SELF_REFER", "cannot refer yourself")
	ErrReferralAlreadyBound  = infraerrors.BadRequest("REFERRAL_ALREADY_BOUND", "user already has a referrer")
	ErrReferralInvalidAmount = infraerrors.BadRequest("REFERRAL_INVALID_AMOUNT", "withdraw amount must be positive")
	ErrReferralInsufficient  = infraerrors.BadRequest("REFERRAL_INSUFFICIENT", "insufficient commission balance")
)

// --- Constants ---

const (
	defaultCommissionRate = 0.10 // 10% default commission rate
	referralCodeLength    = 8    // length of generated referral codes
	settingKeyEnabled     = "referral_enabled"
)

// --- Commission Tier ---

// CommissionTier defines a single tier in the tiered commission rate system.
type CommissionTier struct {
	MinMonthlyRecharge float64 `json:"min_monthly_recharge"`
	MaxMonthlyRecharge float64 `json:"max_monthly_recharge"` // 0 means unlimited
	Rate               float64 `json:"rate"`
}

// DefaultCommissionTiers returns the default tiered commission schedule.
// 10%: monthly recharge < 500
// 12%: 500 – 3000
// 15%: 3000+
var DefaultCommissionTiers = []CommissionTier{
	{MinMonthlyRecharge: 0, MaxMonthlyRecharge: 500, Rate: 0.10},
	{MinMonthlyRecharge: 500, MaxMonthlyRecharge: 3000, Rate: 0.12},
	{MinMonthlyRecharge: 3000, MaxMonthlyRecharge: 0, Rate: 0.15},
}

// --- Response Types ---

type ReferralInfo struct {
	ReferralCode    string           `json:"referral_code"`
	CommissionRate  float64          `json:"commission_rate"`
	TotalCommission float64          `json:"total_commission"`
	TotalInvitees   int              `json:"total_invitees"`
	TotalRecharged  float64          `json:"total_recharged"`
	Balance         float64          `json:"balance"`
	CommissionTiers []CommissionTier `json:"commission_tiers"`
}

type InviteRecord struct {
	ID              int       `json:"id"`
	InviteeEmail    string    `json:"invitee_email"`
	InviteeUsername string    `json:"invitee_username"`
	RegisteredAt    time.Time `json:"registered_at"`
	TotalRecharged  float64   `json:"total_recharged"`
	TotalCommission float64   `json:"total_commission"`
	Status          string    `json:"status"`
}

type CommissionRecord struct {
	ID               int       `json:"id"`
	InviteeEmail     string    `json:"invitee_email"`
	InviteeUsername  string    `json:"invitee_username"`
	OrderID          string    `json:"order_id"`
	RechargeAmount   float64   `json:"recharge_amount"`
	CommissionRate   float64   `json:"commission_rate"`
	CommissionAmount float64   `json:"commission_amount"`
	Status           string    `json:"status"`
	CreatedAt        time.Time `json:"created_at"`
}

type WithdrawResult struct {
	Message string  `json:"message"`
	Balance float64 `json:"balance"`
}

// --- Service ---

type ReferralService struct {
	entClient   *dbent.Client
	settingRepo SettingRepository
	userRepo    UserRepository
}

func NewReferralService(entClient *dbent.Client, settingRepo SettingRepository, userRepo UserRepository) *ReferralService {
	return &ReferralService{
		entClient:   entClient,
		settingRepo: settingRepo,
		userRepo:    userRepo,
	}
}

// getCommissionTiers returns the current commission tier schedule.
func (s *ReferralService) getCommissionTiers() []CommissionTier {
	return DefaultCommissionTiers
}

// getCommissionRateByMonthlyRecharge returns the tiered commission rate
// based on the invitee's total recharge amount in the current calendar month.
func (s *ReferralService) getCommissionRateByMonthlyRecharge(monthlyRecharge float64) float64 {
	tiers := s.getCommissionTiers()
	for i := len(tiers) - 1; i >= 0; i-- {
		if monthlyRecharge >= tiers[i].MinMonthlyRecharge {
			return tiers[i].Rate
		}
	}
	return defaultCommissionRate
}

// getInviteeMonthlyRecharge calculates the invitee's total recharge in the current calendar month
// by summing recharge_amount from commission logs (excluding withdrawn entries).
func (s *ReferralService) getInviteeMonthlyRecharge(ctx context.Context, inviteeID int64) float64 {
	now := time.Now()
	monthStart := time.Date(now.Year(), now.Month(), 1, 0, 0, 0, 0, now.Location())

	logs, err := s.entClient.CommissionLog.Query().
		Where(
			commissionlog.InviteeID(inviteeID),
			commissionlog.StatusNEQ("withdrawn"),
			commissionlog.CreatedAtGTE(monthStart),
		).All(ctx)
	if err != nil {
		return 0
	}

	var total float64
	for _, l := range logs {
		total += l.RechargeAmount
	}
	return total
}

// IsEnabled checks if the referral system is enabled
func (s *ReferralService) IsEnabled(ctx context.Context) bool {
	val, err := s.settingRepo.GetValue(ctx, settingKeyEnabled)
	if err != nil || val == "" {
		return true // enabled by default
	}
	return val == "true" || val == "1"
}

// generateCode creates a random alphanumeric referral code
func generateCode() string {
	b := make([]byte, referralCodeLength)
	rand.Read(b)
	code := strings.ToUpper(hex.EncodeToString(b))[:referralCodeLength]
	return code
}

// GetReferralInfo returns the referral dashboard info for a user
func (s *ReferralService) GetReferralInfo(ctx context.Context, userID int64) (*ReferralInfo, error) {
	u, err := s.entClient.User.Get(ctx, userID)
	if err != nil {
		return nil, ErrUserNotFound
	}

	commissionRate := defaultCommissionRate
	commissionTiers := s.getCommissionTiers()

	// Count invitees
	totalInvitees, err := s.entClient.ReferralRecord.Query().
		Where(referralrecord.ReferrerID(userID)).
		Count(ctx)
	if err != nil {
		totalInvitees = 0
	}

	// Sum total commission and total recharged from referral records
	var totalCommission, totalRecharged float64
	records, err := s.entClient.ReferralRecord.Query().
		Where(referralrecord.ReferrerID(userID)).
		All(ctx)
	if err == nil {
		for _, r := range records {
			totalCommission += r.TotalCommission
			totalRecharged += r.TotalRecharged
		}
	}

	// Calculate available balance (total commission that hasn't been withdrawn)
	// We track withdrawals via commission logs with status "withdrawn"
	withdrawnLogs, err := s.entClient.CommissionLog.Query().
		Where(
			commissionlog.ReferrerID(userID),
			commissionlog.Status("withdrawn"),
		).All(ctx)
	var withdrawn float64
	if err == nil {
		for _, wl := range withdrawnLogs {
			withdrawn += wl.CommissionAmount
		}
	}

	balance := math.Max(0, totalCommission-withdrawn)

	return &ReferralInfo{
		ReferralCode:    u.ReferralCode,
		CommissionRate:  commissionRate,
		TotalCommission: totalCommission,
		TotalInvitees:   totalInvitees,
		TotalRecharged:  totalRecharged,
		Balance:         math.Round(balance*100) / 100,
		CommissionTiers: commissionTiers,
	}, nil
}

// GenerateReferralCode generates or regenerates a referral code for a user
func (s *ReferralService) GenerateReferralCode(ctx context.Context, userID int64) (string, error) {
	code := generateCode()
	// Ensure uniqueness
	for i := 0; i < 5; i++ {
		exists, _ := s.entClient.User.Query().Where(user.ReferralCode(code)).Exist(ctx)
		if !exists {
			break
		}
		code = generateCode()
	}

	_, err := s.entClient.User.UpdateOneID(userID).SetReferralCode(code).Save(ctx)
	if err != nil {
		return "", fmt.Errorf("failed to save referral code: %w", err)
	}
	return code, nil
}

// BindReferrer binds an invitee to a referrer via referral code (called during registration)
func (s *ReferralService) BindReferrer(ctx context.Context, inviteeID int64, referralCode string) error {
	if referralCode == "" {
		return nil
	}

	// Find the referrer by code
	referrer, err := s.entClient.User.Query().
		Where(user.ReferralCode(referralCode)).
		Only(ctx)
	if err != nil {
		slog.Warn("referral code not found during binding", "code", referralCode)
		return nil // Don't fail registration if code is invalid
	}

	if referrer.ID == inviteeID {
		return nil // Silently ignore self-referral
	}

	// Check if invitee already has a referrer
	exists, _ := s.entClient.ReferralRecord.Query().
		Where(referralrecord.InviteeID(inviteeID)).
		Exist(ctx)
	if exists {
		return nil // Already bound
	}

	// Create referral record
	_, err = s.entClient.ReferralRecord.Create().
		SetReferrerID(referrer.ID).
		SetInviteeID(inviteeID).
		SetReferralCode(referralCode).
		SetStatus("active").
		Save(ctx)
	if err != nil {
		slog.Error("failed to create referral record", "referrerID", referrer.ID, "inviteeID", inviteeID, "error", err)
		return nil // Don't fail registration
	}

	// Update invitee's referred_by field
	s.entClient.User.UpdateOneID(inviteeID).SetReferredBy(referrer.ID).Save(ctx)

	slog.Info("referral binding created", "referrerID", referrer.ID, "inviteeID", inviteeID, "code", referralCode)
	return nil
}

// ProcessCommission calculates and credits commission when an invitee recharges
func (s *ReferralService) ProcessCommission(ctx context.Context, inviteeID int64, orderID string, rechargeAmount float64) error {
	if !s.IsEnabled(ctx) {
		return nil
	}
	if rechargeAmount <= 0 {
		return nil
	}

	// Find the referral record for this invitee
	record, err := s.entClient.ReferralRecord.Query().
		Where(
			referralrecord.InviteeID(inviteeID),
			referralrecord.Status("active"),
		).
		Only(ctx)
	if err != nil {
		return nil // No referrer, skip
	}

	commissionRate := s.getCommissionRateByMonthlyRecharge(
		s.getInviteeMonthlyRecharge(ctx, inviteeID) + rechargeAmount,
	)
	commissionAmount := rechargeAmount * commissionRate

	// Check for duplicate commission log (idempotency)
	exists, _ := s.entClient.CommissionLog.Query().
		Where(
			commissionlog.ReferrerID(record.ReferrerID),
			commissionlog.InviteeID(inviteeID),
			commissionlog.OrderID(orderID),
		).Exist(ctx)
	if exists {
		return nil // Already processed
	}

	// Create commission log
	_, err = s.entClient.CommissionLog.Create().
		SetReferrerID(record.ReferrerID).
		SetInviteeID(inviteeID).
		SetOrderID(orderID).
		SetRechargeAmount(rechargeAmount).
		SetCommissionRate(commissionRate).
		SetCommissionAmount(commissionAmount).
		SetStatus("completed").
		Save(ctx)
	if err != nil {
		slog.Error("failed to create commission log", "error", err)
		return fmt.Errorf("create commission log: %w", err)
	}

	// Update referral record totals
	s.entClient.ReferralRecord.UpdateOneID(record.ID).
		SetTotalCommission(record.TotalCommission + commissionAmount).
		SetTotalRecharged(record.TotalRecharged + rechargeAmount).
		Save(ctx)

	slog.Info("commission processed",
		"referrerID", record.ReferrerID,
		"inviteeID", inviteeID,
		"orderID", orderID,
		"rechargeAmount", rechargeAmount,
		"commissionRate", commissionRate,
		"commissionAmount", commissionAmount,
	)
	return nil
}

// GetInviteRecords returns paginated invite records for a referrer
func (s *ReferralService) GetInviteRecords(ctx context.Context, userID int64, params pagination.PaginationParams) ([]InviteRecord, *pagination.PaginationResult, error) {
	query := s.entClient.ReferralRecord.Query().
		Where(referralrecord.ReferrerID(userID))

	total, err := query.Clone().Count(ctx)
	if err != nil {
		return nil, nil, fmt.Errorf("count invites: %w", err)
	}

	records, err := query.
		Order(dbent.Desc(referralrecord.FieldCreatedAt)).
		Offset(params.Offset()).
		Limit(params.Limit()).
		All(ctx)
	if err != nil {
		return nil, nil, fmt.Errorf("list invites: %w", err)
	}

	result := make([]InviteRecord, 0, len(records))
	for _, r := range records {
		// Get invitee info
		invitee, err := s.entClient.User.Get(ctx, r.InviteeID)
		inviteeEmail := ""
		inviteeUsername := ""
		var registeredAt time.Time
		if err == nil {
			inviteeEmail = maskEmail(invitee.Email)
			inviteeUsername = invitee.Username
			registeredAt = invitee.CreatedAt
		}

		result = append(result, InviteRecord{
			ID:              r.ID,
			InviteeEmail:    inviteeEmail,
			InviteeUsername: inviteeUsername,
			RegisteredAt:    registeredAt,
			TotalRecharged:  r.TotalRecharged,
			TotalCommission: r.TotalCommission,
			Status:          r.Status,
		})
	}

	pages := int(math.Ceil(float64(total) / float64(params.Limit())))
	paginationResult := &pagination.PaginationResult{
		Total:    int64(total),
		Page:     params.Page,
		PageSize: params.Limit(),
		Pages:    pages,
	}

	return result, paginationResult, nil
}

// GetCommissionRecords returns paginated commission records for a referrer
func (s *ReferralService) GetCommissionRecords(ctx context.Context, userID int64, params pagination.PaginationParams) ([]CommissionRecord, *pagination.PaginationResult, error) {
	query := s.entClient.CommissionLog.Query().
		Where(
			commissionlog.ReferrerID(userID),
			commissionlog.StatusNEQ("withdrawn"),
		)

	total, err := query.Clone().Count(ctx)
	if err != nil {
		return nil, nil, fmt.Errorf("count commissions: %w", err)
	}

	logs, err := query.
		Order(dbent.Desc(commissionlog.FieldCreatedAt)).
		Offset(params.Offset()).
		Limit(params.Limit()).
		All(ctx)
	if err != nil {
		return nil, nil, fmt.Errorf("list commissions: %w", err)
	}

	result := make([]CommissionRecord, 0, len(logs))
	for _, l := range logs {
		invitee, err := s.entClient.User.Get(ctx, l.InviteeID)
		inviteeEmail := ""
		inviteeUsername := ""
		if err == nil {
			inviteeEmail = maskEmail(invitee.Email)
			inviteeUsername = invitee.Username
		}

		result = append(result, CommissionRecord{
			ID:               l.ID,
			InviteeEmail:     inviteeEmail,
			InviteeUsername:  inviteeUsername,
			OrderID:          l.OrderID,
			RechargeAmount:   l.RechargeAmount,
			CommissionRate:   l.CommissionRate,
			CommissionAmount: l.CommissionAmount,
			Status:           l.Status,
			CreatedAt:        l.CreatedAt,
		})
	}

	pages := int(math.Ceil(float64(total) / float64(params.Limit())))
	paginationResult := &pagination.PaginationResult{
		Total:    int64(total),
		Page:     params.Page,
		PageSize: params.Limit(),
		Pages:    pages,
	}

	return result, paginationResult, nil
}

// WithdrawCommission transfers commission balance to user's main balance
func (s *ReferralService) WithdrawCommission(ctx context.Context, userID int64, amount float64) (*WithdrawResult, error) {
	if amount <= 0 {
		return nil, ErrReferralInvalidAmount
	}

	info, err := s.GetReferralInfo(ctx, userID)
	if err != nil {
		return nil, err
	}

	if amount > info.Balance {
		return nil, ErrReferralInsufficient
	}

	// Add to user's main balance
	if err := s.userRepo.UpdateBalance(ctx, userID, amount); err != nil {
		return nil, fmt.Errorf("update balance: %w", err)
	}

	// Record withdrawal as a commission log entry with status "withdrawn"
	_, err = s.entClient.CommissionLog.Create().
		SetReferrerID(userID).
		SetInviteeID(userID). // self reference for withdrawal record
		SetOrderID(fmt.Sprintf("withdraw_%d_%d", userID, time.Now().UnixMilli())).
		SetRechargeAmount(0).
		SetCommissionRate(0).
		SetCommissionAmount(amount).
		SetStatus("withdrawn").
		SetRemark(fmt.Sprintf("Commission withdrawal: %.2f", amount)).
		Save(ctx)
	if err != nil {
		slog.Error("failed to record withdrawal", "userID", userID, "amount", amount, "error", err)
	}

	// Get updated balance
	u, _ := s.userRepo.GetByID(ctx, userID)
	newBalance := 0.0
	if u != nil {
		newBalance = u.Balance
	}

	slog.Info("commission withdrawn", "userID", userID, "amount", amount, "newBalance", newBalance)

	return &WithdrawResult{
		Message: "Withdrawal successful",
		Balance: newBalance,
	}, nil
}

// maskEmail masks an email for privacy (e.g., "test@example.com" -> "te***@example.com")
func maskEmail(email string) string {
	parts := strings.SplitN(email, "@", 2)
	if len(parts) != 2 {
		return email
	}
	local := parts[0]
	if len(local) <= 2 {
		return local + "***@" + parts[1]
	}
	return local[:2] + "***@" + parts[1]
}

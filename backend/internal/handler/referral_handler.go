package handler

import (
	"strconv"

	"github.com/Wei-Shaw/sub2api/internal/pkg/pagination"
	"github.com/Wei-Shaw/sub2api/internal/pkg/response"
	middleware2 "github.com/Wei-Shaw/sub2api/internal/server/middleware"
	"github.com/Wei-Shaw/sub2api/internal/service"

	"github.com/gin-gonic/gin"
)

// ReferralHandler handles referral/commission system requests
type ReferralHandler struct {
	referralService *service.ReferralService
}

// NewReferralHandler creates a new ReferralHandler
func NewReferralHandler(referralService *service.ReferralService) *ReferralHandler {
	return &ReferralHandler{
		referralService: referralService,
	}
}

// GetReferralInfo returns referral dashboard info
// GET /api/v1/user/referral/info
func (h *ReferralHandler) GetReferralInfo(c *gin.Context) {
	subject, ok := middleware2.GetAuthSubjectFromContext(c)
	if !ok {
		response.Unauthorized(c, "User not authenticated")
		return
	}

	// Check if referral system is enabled
	if !h.referralService.IsEnabled(c.Request.Context()) {
		response.Success(c, gin.H{"enabled": false})
		return
	}

	info, err := h.referralService.GetReferralInfo(c.Request.Context(), subject.UserID)
	if err != nil {
		response.ErrorFrom(c, err)
		return
	}

	response.Success(c, gin.H{
		"enabled":          true,
		"referral_code":    info.ReferralCode,
		"commission_rate":  info.CommissionRate,
		"total_commission": info.TotalCommission,
		"total_invitees":   info.TotalInvitees,
		"total_recharged":  info.TotalRecharged,
		"balance":          info.Balance,
		"commission_tiers": info.CommissionTiers,
	})
}

// GenerateReferralCode generates or regenerates a referral code
// POST /api/v1/user/referral/generate-code
func (h *ReferralHandler) GenerateReferralCode(c *gin.Context) {
	subject, ok := middleware2.GetAuthSubjectFromContext(c)
	if !ok {
		response.Unauthorized(c, "User not authenticated")
		return
	}

	code, err := h.referralService.GenerateReferralCode(c.Request.Context(), subject.UserID)
	if err != nil {
		response.ErrorFrom(c, err)
		return
	}

	response.Success(c, gin.H{"referral_code": code})
}

// GetInviteRecords returns paginated invite records
// GET /api/v1/user/referral/invites
func (h *ReferralHandler) GetInviteRecords(c *gin.Context) {
	subject, ok := middleware2.GetAuthSubjectFromContext(c)
	if !ok {
		response.Unauthorized(c, "User not authenticated")
		return
	}

	page, _ := strconv.Atoi(c.DefaultQuery("page", "1"))
	pageSize, _ := strconv.Atoi(c.DefaultQuery("page_size", "10"))

	params := pagination.PaginationParams{
		Page:     page,
		PageSize: pageSize,
	}

	records, paginationResult, err := h.referralService.GetInviteRecords(c.Request.Context(), subject.UserID, params)
	if err != nil {
		response.ErrorFrom(c, err)
		return
	}

	response.Success(c, gin.H{
		"items":     records,
		"total":     paginationResult.Total,
		"page":      paginationResult.Page,
		"page_size": paginationResult.PageSize,
		"pages":     paginationResult.Pages,
	})
}

// GetCommissionRecords returns paginated commission records
// GET /api/v1/user/referral/commissions
func (h *ReferralHandler) GetCommissionRecords(c *gin.Context) {
	subject, ok := middleware2.GetAuthSubjectFromContext(c)
	if !ok {
		response.Unauthorized(c, "User not authenticated")
		return
	}

	page, _ := strconv.Atoi(c.DefaultQuery("page", "1"))
	pageSize, _ := strconv.Atoi(c.DefaultQuery("page_size", "10"))

	params := pagination.PaginationParams{
		Page:     page,
		PageSize: pageSize,
	}

	records, paginationResult, err := h.referralService.GetCommissionRecords(c.Request.Context(), subject.UserID, params)
	if err != nil {
		response.ErrorFrom(c, err)
		return
	}

	response.Success(c, gin.H{
		"items":     records,
		"total":     paginationResult.Total,
		"page":      paginationResult.Page,
		"page_size": paginationResult.PageSize,
		"pages":     paginationResult.Pages,
	})
}

// WithdrawCommission withdraws commission to user's main balance
// POST /api/v1/user/referral/withdraw
func (h *ReferralHandler) WithdrawCommission(c *gin.Context) {
	subject, ok := middleware2.GetAuthSubjectFromContext(c)
	if !ok {
		response.Unauthorized(c, "User not authenticated")
		return
	}

	var req struct {
		Amount float64 `json:"amount" binding:"required,gt=0"`
	}
	if err := c.ShouldBindJSON(&req); err != nil {
		response.BadRequest(c, "Invalid request: "+err.Error())
		return
	}

	result, err := h.referralService.WithdrawCommission(c.Request.Context(), subject.UserID, req.Amount)
	if err != nil {
		response.ErrorFrom(c, err)
		return
	}

	response.Success(c, result)
}

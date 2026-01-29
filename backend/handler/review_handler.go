package handler

import (
	"strconv"

	"linuxdo-review/dto"
	"linuxdo-review/middleware"
	"linuxdo-review/pkg/response"
	"linuxdo-review/service"

	"github.com/gin-gonic/gin"
)

// ReviewHandler 审核处理器
type ReviewHandler struct {
	reviewService *service.ReviewService
	postService   *service.PostService
}

// NewReviewHandler 创建审核处理器
func NewReviewHandler(reviewService *service.ReviewService, postService *service.PostService) *ReviewHandler {
	return &ReviewHandler{
		reviewService: reviewService,
		postService:   postService,
	}
}

// Approve 通过审核并提交邀请码
func (h *ReviewHandler) Approve(c *gin.Context) {
	idStr := c.Param("id")
	id, err := strconv.ParseUint(idStr, 10, 32)
	if err != nil {
		response.BadRequest(c, "无效的帖子ID")
		return
	}

	var req dto.ApproveRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		response.BadRequest(c, "参数错误: "+err.Error())
		return
	}

	reviewerID := middleware.GetUserID(c)
	// ReviewService 内部会处理邮件通知
	if err := h.reviewService.ApproveWithNotification(uint(id), reviewerID, req.InviteCode); err != nil {
		response.Error(c, err.Error())
		return
	}

	response.SuccessMessage(c, "审核通过，邀请码已发送给申请者")
}

// Reject 拒绝申请
func (h *ReviewHandler) Reject(c *gin.Context) {
	idStr := c.Param("id")
	id, err := strconv.ParseUint(idStr, 10, 32)
	if err != nil {
		response.BadRequest(c, "无效的帖子ID")
		return
	}

	var req dto.RejectRequest
	// 允许不提供拒绝原因
	_ = c.ShouldBindJSON(&req)

	// ReviewService 内部会处理邮件通知
	if err := h.reviewService.RejectWithNotification(uint(id), req.Reason); err != nil {
		response.Error(c, err.Error())
		return
	}

	response.SuccessMessage(c, "已拒绝")
}

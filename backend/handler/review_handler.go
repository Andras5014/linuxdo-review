package handler

import (
	"strconv"
	"strings"

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
	// 检查锁定状态并通过审核
	if err := h.reviewService.CheckLockAndApprove(uint(id), reviewerID, req.InviteCode); err != nil {
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

	userID := middleware.GetUserID(c)
	// 检查锁定状态并拒绝
	if err := h.reviewService.CheckLockAndReject(uint(id), userID, req.Reason); err != nil {
		response.Error(c, err.Error())
		return
	}

	response.SuccessMessage(c, "已拒绝")
}

// GetNext 获取下一个待审核的帖子
func (h *ReviewHandler) GetNext(c *gin.Context) {
	userID := middleware.GetUserID(c)

	// 获取要跳过的帖子ID列表（逗号分隔）
	skipIDsStr := c.Query("skip_ids")
	var skipIDs []uint
	if skipIDsStr != "" {
		for _, idStr := range strings.Split(skipIDsStr, ",") {
			if id, err := strconv.ParseUint(strings.TrimSpace(idStr), 10, 32); err == nil {
				skipIDs = append(skipIDs, uint(id))
			}
		}
	}

	post, err := h.reviewService.GetNextForReview(userID, skipIDs)
	if err != nil {
		// 没有更多待审核的帖子
		response.Success(c, gin.H{
			"post":  nil,
			"total": 0,
		})
		return
	}

	// 获取总数
	total, _ := h.reviewService.GetReviewCount()

	response.Success(c, gin.H{
		"post":  post,
		"total": total,
	})
}

// Skip 跳过当前帖子
func (h *ReviewHandler) Skip(c *gin.Context) {
	idStr := c.Param("id")
	id, err := strconv.ParseUint(idStr, 10, 32)
	if err != nil {
		response.BadRequest(c, "无效的帖子ID")
		return
	}

	userID := middleware.GetUserID(c)
	if err := h.reviewService.SkipPost(uint(id), userID); err != nil {
		response.Error(c, err.Error())
		return
	}

	response.SuccessMessage(c, "已跳过")
}

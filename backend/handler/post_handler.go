package handler

import (
	"strconv"

	"linuxdo-review/dto"
	"linuxdo-review/middleware"
	"linuxdo-review/models"
	"linuxdo-review/pkg/response"
	"linuxdo-review/service"

	"github.com/gin-gonic/gin"
)

// PostHandler 帖子处理器
type PostHandler struct {
	postService   *service.PostService
	reviewService *service.ReviewService
}

// NewPostHandler 创建帖子处理器
func NewPostHandler(postService *service.PostService, reviewService *service.ReviewService) *PostHandler {
	return &PostHandler{
		postService:   postService,
		reviewService: reviewService,
	}
}

// Create 创建帖子
func (h *PostHandler) Create(c *gin.Context) {
	var req dto.CreatePostRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		response.BadRequest(c, "参数错误: "+err.Error())
		return
	}

	userID := middleware.GetUserID(c)
	post, err := h.postService.Create(userID, &req)
	if err != nil {
		response.Error(c, err.Error())
		return
	}

	response.Success(c, dto.ToPostResponse(post))
}

// List 获取帖子列表(申请列表)
func (h *PostHandler) List(c *gin.Context) {
	var req dto.PostListRequest
	if err := c.ShouldBindQuery(&req); err != nil {
		response.BadRequest(c, "参数错误")
		return
	}

	// status 为 nil 时返回所有帖子
	// status 为具体值时返回对应状态的帖子
	posts, total, err := h.postService.ListWithFilter(req.Status, req.GetPage(), req.GetPageSize())
	if err != nil {
		response.Error(c, "获取帖子列表失败")
		return
	}

	// 转换为响应格式
	postResponses := make([]*dto.PostResponse, len(posts))
	postIDs := make([]uint, len(posts))
	for i, post := range posts {
		postResponses[i] = dto.ToPostResponse(post)
		postIDs[i] = post.ID
	}

	// 如果用户已登录,获取用户的投票情况
	userID := middleware.GetUserID(c)
	if userID > 0 {
		votes, _ := h.postService.GetUserVotesForPosts(userID, postIDs)
		for i, resp := range postResponses {
			if vote, ok := votes[resp.ID]; ok {
				postResponses[i].MyVote = vote
			}
		}
	}

	response.Success(c, dto.PaginationResponse{
		List:     postResponses,
		Total:    total,
		Page:     req.GetPage(),
		PageSize: req.GetPageSize(),
	})
}

// ListForReview 获取二级审核列表(认证用户专属)
func (h *PostHandler) ListForReview(c *gin.Context) {
	var pagination dto.PaginationRequest
	if err := c.ShouldBindQuery(&pagination); err != nil {
		response.BadRequest(c, "参数错误")
		return
	}

	posts, total, err := h.postService.ListForSecondReview(pagination.GetPage(), pagination.GetPageSize())
	if err != nil {
		response.Error(c, "获取帖子列表失败")
		return
	}

	postResponses := make([]*dto.PostResponse, len(posts))
	for i, post := range posts {
		postResponses[i] = dto.ToPostResponse(post)
	}

	response.Success(c, dto.PaginationResponse{
		List:     postResponses,
		Total:    total,
		Page:     pagination.GetPage(),
		PageSize: pagination.GetPageSize(),
	})
}

// Get 获取帖子详情
func (h *PostHandler) Get(c *gin.Context) {
	idStr := c.Param("id")
	id, err := strconv.ParseUint(idStr, 10, 32)
	if err != nil {
		response.BadRequest(c, "无效的帖子ID")
		return
	}

	post, err := h.postService.GetByID(uint(id))
	if err != nil {
		response.NotFound(c, "帖子不存在")
		return
	}

	resp := dto.ToPostResponse(post)

	// 如果用户已登录,获取用户的投票情况
	userID := middleware.GetUserID(c)
	if userID > 0 {
		vote, _ := h.postService.GetUserVoteForPost(uint(id), userID)
		resp.MyVote = vote
	}

	response.Success(c, resp)
}

// Vote 投票
func (h *PostHandler) Vote(c *gin.Context) {
	idStr := c.Param("id")
	id, err := strconv.ParseUint(idStr, 10, 32)
	if err != nil {
		response.BadRequest(c, "无效的帖子ID")
		return
	}

	var req dto.VoteRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		response.BadRequest(c, "参数错误: "+err.Error())
		return
	}

	userID := middleware.GetUserID(c)
	voteResp, err := h.postService.VoteWithResponse(uint(id), userID, models.VoteType(req.VoteType))
	if err != nil {
		response.Error(c, err.Error())
		return
	}

	response.Success(c, voteResp)
}

// MyPosts 获取我的帖子列表
func (h *PostHandler) MyPosts(c *gin.Context) {
	var pagination dto.PaginationRequest
	if err := c.ShouldBindQuery(&pagination); err != nil {
		response.BadRequest(c, "参数错误")
		return
	}

	userID := middleware.GetUserID(c)
	posts, total, err := h.postService.ListByUserID(userID, pagination.GetPage(), pagination.GetPageSize())
	if err != nil {
		response.Error(c, "获取帖子列表失败")
		return
	}

	postResponses := make([]*dto.PostResponse, len(posts))
	for i, post := range posts {
		postResponses[i] = dto.ToPostResponse(post)
	}

	response.Success(c, dto.PaginationResponse{
		List:     postResponses,
		Total:    total,
		Page:     pagination.GetPage(),
		PageSize: pagination.GetPageSize(),
	})
}

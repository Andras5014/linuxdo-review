package handler

import (
	"strconv"

	"linuxdo-review/dto"
	"linuxdo-review/models"
	"linuxdo-review/pkg/response"
	"linuxdo-review/service"

	"github.com/gin-gonic/gin"
)

// AdminHandler 管理后台处理器
type AdminHandler struct {
	adminService *service.AdminService
}

// NewAdminHandler 创建管理后台处理器
func NewAdminHandler(adminService *service.AdminService) *AdminHandler {
	return &AdminHandler{
		adminService: adminService,
	}
}

// ListUsers 用户列表
func (h *AdminHandler) ListUsers(c *gin.Context) {
	var pagination dto.PaginationRequest
	if err := c.ShouldBindQuery(&pagination); err != nil {
		response.BadRequest(c, "参数错误")
		return
	}

	users, total, err := h.adminService.ListUsers(pagination.GetPage(), pagination.GetPageSize())
	if err != nil {
		response.Error(c, "获取用户列表失败")
		return
	}

	userResponses := dto.ToUserResponseList(users)

	response.Success(c, dto.NewPaginationResponse(
		userResponses,
		total,
		pagination.GetPage(),
		pagination.GetPageSize(),
	))
}

// GetUser 获取用户详情
func (h *AdminHandler) GetUser(c *gin.Context) {
	idStr := c.Param("id")
	id, err := strconv.ParseUint(idStr, 10, 32)
	if err != nil {
		response.BadRequest(c, "无效的用户ID")
		return
	}

	user, err := h.adminService.GetUserByID(uint(id))
	if err != nil {
		response.NotFound(c, "用户不存在")
		return
	}

	response.Success(c, dto.ToUserResponse(user))
}

// UpdateUserRole 修改用户角色
func (h *AdminHandler) UpdateUserRole(c *gin.Context) {
	idStr := c.Param("id")
	id, err := strconv.ParseUint(idStr, 10, 32)
	if err != nil {
		response.BadRequest(c, "无效的用户ID")
		return
	}

	var req dto.UpdateUserRoleRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		response.BadRequest(c, "参数错误: "+err.Error())
		return
	}

	if err := h.adminService.UpdateUserRole(uint(id), models.UserRole(req.Role)); err != nil {
		response.Error(c, err.Error())
		return
	}

	response.SuccessMessage(c, "更新成功")
}

// GetConfigs 获取配置
func (h *AdminHandler) GetConfigs(c *gin.Context) {
	configs, err := h.adminService.GetConfigs()
	if err != nil {
		response.Error(c, "获取配置失败")
		return
	}

	response.Success(c, dto.ToConfigResponseList(configs))
}

// UpdateConfig 更新配置
func (h *AdminHandler) UpdateConfig(c *gin.Context) {
	var req dto.UpdateConfigRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		response.BadRequest(c, "参数错误: "+err.Error())
		return
	}

	if err := h.adminService.UpdateConfig(req.Key, req.Value, ""); err != nil {
		response.Error(c, err.Error())
		return
	}

	response.SuccessMessage(c, "更新成功")
}

// BatchUpdateConfigs 批量更新配置
func (h *AdminHandler) BatchUpdateConfigs(c *gin.Context) {
	var req dto.BatchUpdateConfigRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		response.BadRequest(c, "参数错误: "+err.Error())
		return
	}

	if err := h.adminService.BatchUpdateConfigs(req.Configs); err != nil {
		response.Error(c, err.Error())
		return
	}

	response.SuccessMessage(c, "批量更新成功")
}

// GetStats 数据统计
func (h *AdminHandler) GetStats(c *gin.Context) {
	stats, err := h.adminService.GetStats()
	if err != nil {
		response.Error(c, "获取统计数据失败")
		return
	}

	response.Success(c, stats)
}

package handlers

import (
	"seldom-platform/database"
	"seldom-platform/models"
	"seldom-platform/utils"
	"strconv"

	"github.com/gin-gonic/gin"
)

// TeamHandler 团队处理器
type TeamHandler struct{}

// NewTeamHandler 创建团队处理器
func NewTeamHandler() *TeamHandler {
	return &TeamHandler{}
}

// CreateTeamRequest 创建团队请求结构
type CreateTeamRequest struct {
	Name        string `json:"name" binding:"required"`
	Description string `json:"description"`
	Owner       uint   `json:"owner" binding:"required"`
}

// UpdateTeamRequest 更新团队请求结构
type UpdateTeamRequest struct {
	Name        string `json:"name"`
	Description string `json:"description"`
	Owner       uint   `json:"owner"`
}

// GetTeams 获取团队列表
// @Summary 获取团队列表
// @Description 获取团队列表，支持分页和搜索
// @Tags 团队管理
// @Produce json
// @Security BearerAuth
// @Param page query int false "页码" default(1)
// @Param size query int false "每页数量" default(10)
// @Param search query string false "搜索关键词"
// @Success 200 {object} utils.PageResponse{data=[]models.Team}
// @Failure 401 {object} utils.Response
// @Router /api/teams [get]
func (h *TeamHandler) GetTeams(c *gin.Context) {
	db := database.GetDB()

	// 获取分页参数
	page, _ := strconv.Atoi(c.DefaultQuery("page", "1"))
	size, _ := strconv.Atoi(c.DefaultQuery("size", "10"))
	search := c.Query("search")

	if page < 1 {
		page = 1
	}
	if size < 1 || size > 100 {
		size = 10
	}

	offset := (page - 1) * size

	// 构建查询
	query := db.Model(&models.Team{})
	if search != "" {
		query = query.Where("name LIKE ? OR description LIKE ?", "%"+search+"%", "%"+search+"%")
	}

	// 获取总数
	var total int64
	query.Count(&total)

	// 获取数据
	var teams []models.Team
	if err := query.Offset(offset).Limit(size).Find(&teams).Error; err != nil {
		utils.InternalServerError(c, "Failed to fetch teams")
		return
	}

	utils.PageSuccess(c, teams, total, page, size)
}

// GetTeam 获取团队详情
// @Summary 获取团队详情
// @Description 根据ID获取团队详情
// @Tags 团队管理
// @Produce json
// @Security BearerAuth
// @Param id path int true "团队ID"
// @Success 200 {object} utils.Response{data=models.Team}
// @Failure 404 {object} utils.Response
// @Router /api/teams/{id} [get]
func (h *TeamHandler) GetTeam(c *gin.Context) {
	id := c.Param("id")
	db := database.GetDB()

	var team models.Team
	if err := db.First(&team, id).Error; err != nil {
		utils.NotFound(c, "Team not found")
		return
	}

	utils.Success(c, team)
}

// CreateTeam 创建团队
// @Summary 创建团队
// @Description 创建新的团队
// @Tags 团队管理
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param team body CreateTeamRequest true "团队信息"
// @Success 200 {object} utils.Response{data=models.Team}
// @Failure 400 {object} utils.Response
// @Router /api/teams [post]
func (h *TeamHandler) CreateTeam(c *gin.Context) {
	var req CreateTeamRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		utils.BadRequest(c, "Invalid request format")
		return
	}

	db := database.GetDB()

	// 检查团队名是否已存在
	var existingTeam models.Team
	if err := db.Where("name = ?", req.Name).First(&existingTeam).Error; err == nil {
		utils.BadRequest(c, "Team name already exists")
		return
	}

	// 创建团队
	team := models.Team{
		Name:  req.Name,
		Email: req.Description, // 将Description映射到Email字段
	}

	if err := db.Create(&team).Error; err != nil {
		utils.InternalServerError(c, "Failed to create team")
		return
	}

	utils.SuccessWithMessage(c, "Team created successfully", team)
}

// UpdateTeam 更新团队
// @Summary 更新团队
// @Description 更新团队信息
// @Tags 团队管理
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param id path int true "团队ID"
// @Param team body UpdateTeamRequest true "团队信息"
// @Success 200 {object} utils.Response{data=models.Team}
// @Failure 400 {object} utils.Response
// @Failure 404 {object} utils.Response
// @Router /api/teams/{id} [put]
func (h *TeamHandler) UpdateTeam(c *gin.Context) {
	id := c.Param("id")
	var req UpdateTeamRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		utils.BadRequest(c, "Invalid request format")
		return
	}

	db := database.GetDB()
	var team models.Team

	if err := db.First(&team, id).Error; err != nil {
		utils.NotFound(c, "Team not found")
		return
	}

	// 更新团队信息
	if req.Name != "" {
		team.Name = req.Name
	}
	if req.Description != "" {
		team.Email = req.Description
	}

	if err := db.Save(&team).Error; err != nil {
		utils.InternalServerError(c, "Failed to update team")
		return
	}

	utils.SuccessWithMessage(c, "Team updated successfully", team)
}

// DeleteTeam 删除团队
// @Summary 删除团队
// @Description 删除团队
// @Tags 团队管理
// @Produce json
// @Security BearerAuth
// @Param id path int true "团队ID"
// @Success 200 {object} utils.Response
// @Failure 404 {object} utils.Response
// @Router /api/teams/{id} [delete]
func (h *TeamHandler) DeleteTeam(c *gin.Context) {
	id := c.Param("id")
	db := database.GetDB()

	var team models.Team
	if err := db.First(&team, id).Error; err != nil {
		utils.NotFound(c, "Team not found")
		return
	}

	if err := db.Delete(&team).Error; err != nil {
		utils.InternalServerError(c, "Failed to delete team")
		return
	}

	utils.SuccessWithMessage(c, "Team deleted successfully", nil)
}
package handlers

import (
	"seldom-platform/database"
	"seldom-platform/models"
	"seldom-platform/utils"
	"strconv"

	"github.com/gin-gonic/gin"
)

// ProjectHandler 项目处理器
type ProjectHandler struct{}

// NewProjectHandler 创建项目处理器
func NewProjectHandler() *ProjectHandler {
	return &ProjectHandler{}
}

// CreateProjectRequest 创建项目请求结构
type CreateProjectRequest struct {
	Name        string `json:"name" binding:"required"`
	Description string `json:"description"`
	Image       string `json:"image"`
	Host        string `json:"host"`
	Port        int    `json:"port"`
	Owner       uint   `json:"owner"`
}

// UpdateProjectRequest 更新项目请求结构
type UpdateProjectRequest struct {
	Name        string `json:"name"`
	Description string `json:"description"`
	Image       string `json:"image"`
	Host        string `json:"host"`
	Port        int    `json:"port"`
	Owner       uint   `json:"owner"`
}

// GetProjects 获取项目列表
// @Summary 获取项目列表
// @Description 获取项目列表，支持分页
// @Tags 项目管理
// @Produce json
// @Security BearerAuth
// @Param page query int false "页码" default(1)
// @Param size query int false "每页数量" default(10)
// @Param search query string false "搜索关键词"
// @Success 200 {object} utils.PageResponse{data=[]models.Project}
// @Failure 401 {object} utils.Response
// @Router /api/projects [get]
func (h *ProjectHandler) GetProjects(c *gin.Context) {
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
	query := db.Model(&models.Project{})
	if search != "" {
		query = query.Where("name LIKE ? OR description LIKE ?", "%"+search+"%", "%"+search+"%")
	}

	// 获取总数
	var total int64
	query.Count(&total)

	// 获取数据
	var projects []models.Project
	if err := query.Offset(offset).Limit(size).Find(&projects).Error; err != nil {
		utils.InternalServerError(c, "Failed to fetch projects")
		return
	}

	utils.PageSuccess(c, projects, total, page, size)
}

// GetProject 获取项目详情
// @Summary 获取项目详情
// @Description 根据ID获取项目详情
// @Tags 项目管理
// @Produce json
// @Security BearerAuth
// @Param id path int true "项目ID"
// @Success 200 {object} utils.Response{data=models.Project}
// @Failure 404 {object} utils.Response
// @Router /api/projects/{id} [get]
func (h *ProjectHandler) GetProject(c *gin.Context) {
	id := c.Param("id")
	db := database.GetDB()

	var project models.Project
	if err := db.First(&project, id).Error; err != nil {
		utils.NotFound(c, "Project not found")
		return
	}

	utils.Success(c, project)
}

// CreateProject 创建项目
// @Summary 创建项目
// @Description 创建新项目
// @Tags 项目管理
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param project body CreateProjectRequest true "项目信息"
// @Success 200 {object} utils.Response{data=models.Project}
// @Failure 400 {object} utils.Response
// @Router /api/projects [post]
func (h *ProjectHandler) CreateProject(c *gin.Context) {
	var req CreateProjectRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		utils.BadRequest(c, "Invalid request format")
		return
	}

	db := database.GetDB()

	// 检查项目名是否已存在
	var existingProject models.Project
	if err := db.Where("name = ?", req.Name).First(&existingProject).Error; err == nil {
		utils.BadRequest(c, "Project name already exists")
		return
	}

	// 创建项目
	project := models.Project{
		Name:       req.Name,
		Address:    req.Host, // 将Host映射到Address字段
		CaseDir:    "test_dir",
		CoverName:  req.Image,
		PathName:   req.Image,
	}

	if err := db.Create(&project).Error; err != nil {
		utils.InternalServerError(c, "Failed to create project")
		return
	}

	utils.SuccessWithMessage(c, "Project created successfully", project)
}

// UpdateProject 更新项目
// @Summary 更新项目
// @Description 更新项目信息
// @Tags 项目管理
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param id path int true "项目ID"
// @Param project body UpdateProjectRequest true "项目信息"
// @Success 200 {object} utils.Response{data=models.Project}
// @Failure 400 {object} utils.Response
// @Failure 404 {object} utils.Response
// @Router /api/projects/{id} [put]
func (h *ProjectHandler) UpdateProject(c *gin.Context) {
	id := c.Param("id")
	var req UpdateProjectRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		utils.BadRequest(c, "Invalid request format")
		return
	}

	db := database.GetDB()
	var project models.Project

	if err := db.First(&project, id).Error; err != nil {
		utils.NotFound(c, "Project not found")
		return
	}

	// 更新项目信息
	if req.Name != "" {
		project.Name = req.Name
	}
	if req.Image != "" {
		project.CoverName = req.Image
		project.PathName = req.Image
	}
	if req.Host != "" {
		project.Address = req.Host
	}

	if err := db.Save(&project).Error; err != nil {
		utils.InternalServerError(c, "Failed to update project")
		return
	}

	utils.SuccessWithMessage(c, "Project updated successfully", project)
}

// DeleteProject 删除项目
// @Summary 删除项目
// @Description 删除项目
// @Tags 项目管理
// @Produce json
// @Security BearerAuth
// @Param id path int true "项目ID"
// @Success 200 {object} utils.Response
// @Failure 404 {object} utils.Response
// @Router /api/projects/{id} [delete]
func (h *ProjectHandler) DeleteProject(c *gin.Context) {
	id := c.Param("id")
	db := database.GetDB()

	var project models.Project
	if err := db.First(&project, id).Error; err != nil {
		utils.NotFound(c, "Project not found")
		return
	}

	if err := db.Delete(&project).Error; err != nil {
		utils.InternalServerError(c, "Failed to delete project")
		return
	}

	utils.SuccessWithMessage(c, "Project deleted successfully", nil)
}
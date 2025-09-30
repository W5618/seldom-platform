package handlers

import (
	"seldom-platform/database"
	"seldom-platform/models"
	"seldom-platform/utils"
	"strconv"

	"github.com/gin-gonic/gin"
)

// EnvHandler 环境处理器
type EnvHandler struct{}

// NewEnvHandler 创建环境处理器
func NewEnvHandler() *EnvHandler {
	return &EnvHandler{}
}

// CreateEnvRequest 创建环境请求结构
type CreateEnvRequest struct {
	Name        string `json:"name" binding:"required"`
	Host        string `json:"host" binding:"required"`
	Port        int    `json:"port"`
	Protocol    string `json:"protocol"`
	Description string `json:"description"`
	Project     uint   `json:"project" binding:"required"`
}

// UpdateEnvRequest 更新环境请求结构
type UpdateEnvRequest struct {
	Name        string `json:"name"`
	Host        string `json:"host"`
	Port        int    `json:"port"`
	Protocol    string `json:"protocol"`
	Description string `json:"description"`
	Project     uint   `json:"project"`
}

// GetEnvs 获取环境列表
// @Summary 获取环境列表
// @Description 获取环境列表，支持分页和筛选
// @Tags 环境管理
// @Produce json
// @Security BearerAuth
// @Param page query int false "页码" default(1)
// @Param size query int false "每页数量" default(10)
// @Param project query int false "项目ID"
// @Success 200 {object} utils.PageResponse{data=[]models.Env}
// @Failure 401 {object} utils.Response
// @Router /api/envs [get]
func (h *EnvHandler) GetEnvs(c *gin.Context) {
	db := database.GetDB()

	// 获取分页参数
	page, _ := strconv.Atoi(c.DefaultQuery("page", "1"))
	size, _ := strconv.Atoi(c.DefaultQuery("size", "10"))
	projectID := c.Query("project")

	if page < 1 {
		page = 1
	}
	if size < 1 || size > 100 {
		size = 10
	}

	offset := (page - 1) * size

	// 构建查询
	query := db.Model(&models.Env{})
	if projectID != "" {
		query = query.Where("project = ?", projectID)
	}

	// 获取总数
	var total int64
	query.Count(&total)

	// 获取数据
	var envs []models.Env
	if err := query.Offset(offset).Limit(size).Find(&envs).Error; err != nil {
		utils.InternalServerError(c, "Failed to fetch environments")
		return
	}

	utils.PageSuccess(c, envs, total, page, size)
}

// GetEnv 获取环境详情
// @Summary 获取环境详情
// @Description 根据ID获取环境详情
// @Tags 环境管理
// @Produce json
// @Security BearerAuth
// @Param id path int true "环境ID"
// @Success 200 {object} utils.Response{data=models.Env}
// @Failure 404 {object} utils.Response
// @Router /api/envs/{id} [get]
func (h *EnvHandler) GetEnv(c *gin.Context) {
	id := c.Param("id")
	db := database.GetDB()

	var env models.Env
	if err := db.First(&env, id).Error; err != nil {
		utils.NotFound(c, "Environment not found")
		return
	}

	utils.Success(c, env)
}

// CreateEnv 创建环境
// @Summary 创建环境
// @Description 创建新的环境
// @Tags 环境管理
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param env body CreateEnvRequest true "环境信息"
// @Success 200 {object} utils.Response{data=models.Env}
// @Failure 400 {object} utils.Response
// @Router /api/envs [post]
func (h *EnvHandler) CreateEnv(c *gin.Context) {
	var req CreateEnvRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		utils.BadRequest(c, "Invalid request format")
		return
	}

	db := database.GetDB()

	// 检查环境名在同一项目下是否已存在
	var existingEnv models.Env
	if err := db.Where("name = ?", req.Name).First(&existingEnv).Error; err == nil {
		utils.BadRequest(c, "Environment name already exists")
		return
	}

	// 创建环境
	env := models.Env{
		Name:         req.Name,
		TestType:     "http",
		Env:          req.Protocol + "://" + req.Host,
		BaseURL:      req.Protocol + "://" + req.Host,
		Browser:      "chrome",
	}

	if err := db.Create(&env).Error; err != nil {
		utils.InternalServerError(c, "Failed to create environment")
		return
	}

	utils.SuccessWithMessage(c, "Environment created successfully", env)
}

// UpdateEnv 更新环境
// @Summary 更新环境
// @Description 更新环境信息
// @Tags 环境管理
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param id path int true "环境ID"
// @Param env body UpdateEnvRequest true "环境信息"
// @Success 200 {object} utils.Response{data=models.Env}
// @Failure 400 {object} utils.Response
// @Failure 404 {object} utils.Response
// @Router /api/envs/{id} [put]
func (h *EnvHandler) UpdateEnv(c *gin.Context) {
	id := c.Param("id")
	var req UpdateEnvRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		utils.BadRequest(c, "Invalid request format")
		return
	}

	db := database.GetDB()
	var env models.Env

	if err := db.First(&env, id).Error; err != nil {
		utils.NotFound(c, "Environment not found")
		return
	}

	// 更新环境信息
	if req.Name != "" {
		env.Name = req.Name
	}
	if req.Host != "" || req.Protocol != "" {
		// 更新BaseURL和Env字段
		protocol := req.Protocol
		if protocol == "" {
			// 从现有的BaseURL中提取协议
			if env.BaseURL != "" {
				if env.BaseURL[:5] == "https" {
					protocol = "https"
				} else {
					protocol = "http"
				}
			} else {
				protocol = "http"
			}
		}
		host := req.Host
		if host == "" {
			// 从现有的BaseURL中提取主机
			if env.BaseURL != "" {
				if protocol == "https" {
					host = env.BaseURL[8:] // 去掉 "https://"
				} else {
					host = env.BaseURL[7:] // 去掉 "http://"
				}
			}
		}
		env.BaseURL = protocol + "://" + host
		env.Env = protocol + "://" + host
	}

	if err := db.Save(&env).Error; err != nil {
		utils.InternalServerError(c, "Failed to update environment")
		return
	}

	utils.SuccessWithMessage(c, "Environment updated successfully", env)
}

// DeleteEnv 删除环境
// @Summary 删除环境
// @Description 删除环境
// @Tags 环境管理
// @Produce json
// @Security BearerAuth
// @Param id path int true "环境ID"
// @Success 200 {object} utils.Response
// @Failure 404 {object} utils.Response
// @Router /api/envs/{id} [delete]
func (h *EnvHandler) DeleteEnv(c *gin.Context) {
	id := c.Param("id")
	db := database.GetDB()

	var env models.Env
	if err := db.First(&env, id).Error; err != nil {
		utils.NotFound(c, "Environment not found")
		return
	}

	if err := db.Delete(&env).Error; err != nil {
		utils.InternalServerError(c, "Failed to delete environment")
		return
	}

	utils.SuccessWithMessage(c, "Environment deleted successfully", nil)
}
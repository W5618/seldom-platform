package handlers

import (
	"seldom-platform/database"
	"seldom-platform/models"
	"seldom-platform/utils"
	"strconv"

	"github.com/gin-gonic/gin"
)

// CaseHandler 测试用例处理器
type CaseHandler struct{}

// NewCaseHandler 创建测试用例处理器
func NewCaseHandler() *CaseHandler {
	return &CaseHandler{}
}

// CreateCaseRequest 创建测试用例请求结构
type CreateCaseRequest struct {
	Name        string `json:"name" binding:"required"`
	Info        string `json:"info"`
	Project     uint   `json:"project" binding:"required"`
	Module      string `json:"module"`
	Author      uint   `json:"author"`
	Include     string `json:"include"`
	Request     string `json:"request"`
	Tag         string `json:"tag"`
	Relation    int    `json:"relation"`
	Priority    int    `json:"priority"`
	Status      int    `json:"status"`
}

// UpdateCaseRequest 更新测试用例请求结构
type UpdateCaseRequest struct {
	Name     string `json:"name"`
	Info     string `json:"info"`
	Project  uint   `json:"project"`
	Module   string `json:"module"`
	Author   uint   `json:"author"`
	Include  string `json:"include"`
	Request  string `json:"request"`
	Tag      string `json:"tag"`
	Relation int    `json:"relation"`
	Priority int    `json:"priority"`
	Status   int    `json:"status"`
}

// GetCases 获取测试用例列表
// @Summary 获取测试用例列表
// @Description 获取测试用例列表，支持分页和筛选
// @Tags 测试用例管理
// @Produce json
// @Security BearerAuth
// @Param page query int false "页码" default(1)
// @Param size query int false "每页数量" default(10)
// @Param project query int false "项目ID"
// @Param search query string false "搜索关键词"
// @Success 200 {object} utils.PageResponse{data=[]models.TestCase}
// @Failure 401 {object} utils.Response
// @Router /api/cases [get]
func (h *CaseHandler) GetCases(c *gin.Context) {
	db := database.GetDB()

	// 获取分页参数
	page, _ := strconv.Atoi(c.DefaultQuery("page", "1"))
	size, _ := strconv.Atoi(c.DefaultQuery("size", "10"))
	projectID := c.Query("project")
	search := c.Query("search")

	if page < 1 {
		page = 1
	}
	if size < 1 || size > 100 {
		size = 10
	}

	offset := (page - 1) * size

	// 构建查询
	query := db.Model(&models.TestCase{})
	
	if projectID != "" {
		query = query.Where("project = ?", projectID)
	}
	
	if search != "" {
		query = query.Where("name LIKE ? OR info LIKE ?", "%"+search+"%", "%"+search+"%")
	}

	// 获取总数
	var total int64
	query.Count(&total)

	// 获取数据
	var cases []models.TestCase
	if err := query.Offset(offset).Limit(size).Find(&cases).Error; err != nil {
		utils.InternalServerError(c, "Failed to fetch test cases")
		return
	}

	utils.PageSuccess(c, cases, total, page, size)
}

// GetCase 获取测试用例详情
// @Summary 获取测试用例详情
// @Description 根据ID获取测试用例详情
// @Tags 测试用例管理
// @Produce json
// @Security BearerAuth
// @Param id path int true "测试用例ID"
// @Success 200 {object} utils.Response{data=models.TestCase}
// @Failure 404 {object} utils.Response
// @Router /api/cases/{id} [get]
func (h *CaseHandler) GetCase(c *gin.Context) {
	id := c.Param("id")
	db := database.GetDB()

	var testCase models.TestCase
	if err := db.First(&testCase, id).Error; err != nil {
		utils.NotFound(c, "Test case not found")
		return
	}

	utils.Success(c, testCase)
}

// CreateCase 创建测试用例
// @Summary 创建测试用例
// @Description 创建新的测试用例
// @Tags 测试用例管理
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param case body CreateCaseRequest true "测试用例信息"
// @Success 200 {object} utils.Response{data=models.TestCase}
// @Failure 400 {object} utils.Response
// @Router /api/cases [post]
func (h *CaseHandler) CreateCase(c *gin.Context) {
	var req CreateCaseRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		utils.BadRequest(c, "Invalid request format")
		return
	}

	db := database.GetDB()

	// 创建测试用例
	testCase := models.TestCase{
		ProjectID:  req.Project,
		FileName:   req.Name,
		ClassName:  req.Module,
		ClassDoc:   req.Info,
		CaseName:   req.Name,
		CaseDoc:    req.Info,
		Label:      req.Tag,
		Status:     req.Status,
		CaseHash:   utils.GenerateMD5(req.Name + req.Info),
	}

	if err := db.Create(&testCase).Error; err != nil {
		utils.InternalServerError(c, "Failed to create test case")
		return
	}

	utils.SuccessWithMessage(c, "Test case created successfully", testCase)
}

// UpdateCase 更新测试用例
// @Summary 更新测试用例
// @Description 更新测试用例信息
// @Tags 测试用例管理
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param id path int true "测试用例ID"
// @Param case body UpdateCaseRequest true "测试用例信息"
// @Success 200 {object} utils.Response{data=models.TestCase}
// @Failure 400 {object} utils.Response
// @Failure 404 {object} utils.Response
// @Router /api/cases/{id} [put]
func (h *CaseHandler) UpdateCase(c *gin.Context) {
	id := c.Param("id")
	var req UpdateCaseRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		utils.BadRequest(c, "Invalid request format")
		return
	}

	db := database.GetDB()
	var testCase models.TestCase

	if err := db.First(&testCase, id).Error; err != nil {
		utils.NotFound(c, "Test case not found")
		return
	}

	// 更新测试用例信息
	if req.Name != "" {
		testCase.FileName = req.Name
		testCase.CaseName = req.Name
	}
	if req.Info != "" {
		testCase.ClassDoc = req.Info
		testCase.CaseDoc = req.Info
	}
	if req.Project != 0 {
		testCase.ProjectID = req.Project
	}
	if req.Module != "" {
		testCase.ClassName = req.Module
	}
	if req.Tag != "" {
		testCase.Label = req.Tag
	}
	testCase.Status = req.Status

	// 更新hash值
	if req.Name != "" || req.Info != "" {
		testCase.CaseHash = utils.GenerateMD5(testCase.CaseName + testCase.CaseDoc)
	}

	if err := db.Save(&testCase).Error; err != nil {
		utils.InternalServerError(c, "Failed to update test case")
		return
	}

	utils.SuccessWithMessage(c, "Test case updated successfully", testCase)
}

// DeleteCase 删除测试用例
// @Summary 删除测试用例
// @Description 删除测试用例
// @Tags 测试用例管理
// @Produce json
// @Security BearerAuth
// @Param id path int true "测试用例ID"
// @Success 200 {object} utils.Response
// @Failure 404 {object} utils.Response
// @Router /api/cases/{id} [delete]
func (h *CaseHandler) DeleteCase(c *gin.Context) {
	id := c.Param("id")
	db := database.GetDB()

	var testCase models.TestCase
	if err := db.First(&testCase, id).Error; err != nil {
		utils.NotFound(c, "Test case not found")
		return
	}

	if err := db.Delete(&testCase).Error; err != nil {
		utils.InternalServerError(c, "Failed to delete test case")
		return
	}

	utils.SuccessWithMessage(c, "Test case deleted successfully", nil)
}

// CopyCase 复制测试用例
// @Summary 复制测试用例
// @Description 复制现有测试用例
// @Tags 测试用例管理
// @Produce json
// @Security BearerAuth
// @Param id path int true "测试用例ID"
// @Success 200 {object} utils.Response{data=models.TestCase}
// @Failure 404 {object} utils.Response
// @Router /api/cases/{id}/copy [post]
func (h *CaseHandler) CopyCase(c *gin.Context) {
	id := c.Param("id")
	db := database.GetDB()

	var originalCase models.TestCase
	if err := db.First(&originalCase, id).Error; err != nil {
		utils.NotFound(c, "Test case not found")
		return
	}

	// 创建副本
	newCase := models.TestCase{
		ProjectID:  originalCase.ProjectID,
		FileName:   originalCase.FileName + " (Copy)",
		ClassName:  originalCase.ClassName,
		ClassDoc:   originalCase.ClassDoc,
		CaseName:   originalCase.CaseName + " (Copy)",
		CaseDoc:    originalCase.CaseDoc,
		Label:      originalCase.Label,
		Status:     0, // 新副本状态设为未执行
		CaseHash:   utils.GenerateMD5(originalCase.CaseName + " (Copy)" + originalCase.CaseDoc),
	}

	if err := db.Create(&newCase).Error; err != nil {
		utils.InternalServerError(c, "Failed to copy test case")
		return
	}

	utils.SuccessWithMessage(c, "Test case copied successfully", newCase)
}
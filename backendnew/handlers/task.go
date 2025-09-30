package handlers

import (
	"seldom-platform/database"
	"seldom-platform/models"
	"seldom-platform/services"
	"seldom-platform/utils"
	"strconv"

	"github.com/gin-gonic/gin"
)

// TaskHandler 任务处理器
type TaskHandler struct{}

// NewTaskHandler 创建任务处理器
func NewTaskHandler() *TaskHandler {
	return &TaskHandler{}
}

// CreateTaskRequest 创建任务请求结构
type CreateTaskRequest struct {
	Name           string `json:"name" binding:"required"`
	Project        uint   `json:"project" binding:"required"`
	Env            uint   `json:"env"`
	CronTime       string `json:"cron_time"`
	CronExpression string `json:"cron_expression"`
	IsScheduled    bool   `json:"is_scheduled"`
	Type           int    `json:"type"`
	Status         int    `json:"status"`
	CaseList       string `json:"case_list"`
	Email          string `json:"email"`
	DingTalk       string `json:"ding_talk"`
	WebHook        string `json:"web_hook"`
	Performer      uint   `json:"performer"`
}

// UpdateTaskRequest 更新任务请求结构
type UpdateTaskRequest struct {
	Name        string `json:"name"`
	Project     uint   `json:"project"`
	Env         uint   `json:"env"`
	CronTime    string `json:"cron_time"`
	Type        int    `json:"type"`
	Status      int    `json:"status"`
	CaseList    string `json:"case_list"`
	Email       string `json:"email"`
	DingTalk    string `json:"ding_talk"`
	WebHook     string `json:"web_hook"`
	Performer   uint   `json:"performer"`
}

// GetTasks 获取任务列表
// @Summary 获取任务列表
// @Description 获取任务列表，支持分页和筛选
// @Tags 任务管理
// @Produce json
// @Security BearerAuth
// @Param page query int false "页码" default(1)
// @Param size query int false "每页数量" default(10)
// @Param project query int false "项目ID"
// @Success 200 {object} utils.PageResponse{data=[]models.TestTask}
// @Failure 401 {object} utils.Response
// @Router /api/tasks [get]
func (h *TaskHandler) GetTasks(c *gin.Context) {
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
	query := db.Model(&models.TestTask{})
	if projectID != "" {
		query = query.Where("project = ?", projectID)
	}

	// 获取总数
	var total int64
	query.Count(&total)

	// 获取数据
	var tasks []models.TestTask
	if err := query.Offset(offset).Limit(size).Find(&tasks).Error; err != nil {
		utils.InternalServerError(c, "Failed to fetch tasks")
		return
	}

	utils.PageSuccess(c, tasks, total, page, size)
}

// GetTask 获取任务详情
// @Summary 获取任务详情
// @Description 根据ID获取任务详情
// @Tags 任务管理
// @Produce json
// @Security BearerAuth
// @Param id path int true "任务ID"
// @Success 200 {object} utils.Response{data=models.TestTask}
// @Failure 404 {object} utils.Response
// @Router /api/tasks/{id} [get]
func (h *TaskHandler) GetTask(c *gin.Context) {
	id := c.Param("id")
	db := database.GetDB()

	var task models.TestTask
	if err := db.First(&task, id).Error; err != nil {
		utils.NotFound(c, "Task not found")
		return
	}

	utils.Success(c, task)
}

// CreateTask 创建任务
// @Summary 创建任务
// @Description 创建新的任务
// @Tags 任务管理
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param task body CreateTaskRequest true "任务信息"
// @Success 200 {object} utils.Response{data=models.TestTask}
// @Failure 400 {object} utils.Response
// @Router /api/tasks [post]
func (h *TaskHandler) CreateTask(c *gin.Context) {
	var req CreateTaskRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		utils.BadRequest(c, "Invalid request format")
		return
	}

	db := database.GetDB()

	// 创建任务
	task := models.TestTask{
		Name:           req.Name,
		ProjectID:      req.Project,
		EnvID:          &req.Env,
		Timed:          req.CronTime,
		CronExpression: req.CronExpression,
		IsScheduled:    req.IsScheduled,
		Status:         req.Status,
		Email:          req.Email,
	}

	if err := db.Create(&task).Error; err != nil {
		utils.InternalServerError(c, "Failed to create task")
		return
	}

	// 如果任务有cron表达式且设置为定时任务，添加到调度器
	if req.CronExpression != "" && req.IsScheduled && req.Status == 1 {
		if services.GlobalScheduler != nil {
			if err := services.GlobalScheduler.AddTask(task.ID); err != nil {
				// 记录错误但不影响任务创建
				// 可以考虑添加日志记录
			}
		}
	}

	utils.SuccessWithMessage(c, "Task created successfully", task)
}

// UpdateTask 更新任务
// @Summary 更新任务
// @Description 更新任务信息
// @Tags 任务管理
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param id path int true "任务ID"
// @Param task body UpdateTaskRequest true "任务信息"
// @Success 200 {object} utils.Response{data=models.TestTask}
// @Failure 400 {object} utils.Response
// @Failure 404 {object} utils.Response
// @Router /api/tasks/{id} [put]
func (h *TaskHandler) UpdateTask(c *gin.Context) {
	id := c.Param("id")
	var req UpdateTaskRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		utils.BadRequest(c, "Invalid request format")
		return
	}

	db := database.GetDB()
	var task models.TestTask

	if err := db.First(&task, id).Error; err != nil {
		utils.NotFound(c, "Task not found")
		return
	}

	// 更新任务信息
	if req.Name != "" {
		task.Name = req.Name
	}
	if req.Project != 0 {
		task.ProjectID = req.Project
	}
	if req.Env != 0 {
		task.EnvID = &req.Env
	}
	if req.CronTime != "" {
		task.Timed = req.CronTime
	}
	task.Status = req.Status
	if req.Email != "" {
		task.Email = req.Email
	}

	if err := db.Save(&task).Error; err != nil {
		utils.InternalServerError(c, "Failed to update task")
		return
	}

	utils.SuccessWithMessage(c, "Task updated successfully", task)
}

// DeleteTask 删除任务
// @Summary 删除任务
// @Description 删除任务
// @Tags 任务管理
// @Produce json
// @Security BearerAuth
// @Param id path int true "任务ID"
// @Success 200 {object} utils.Response
// @Failure 404 {object} utils.Response
// @Router /api/tasks/{id} [delete]
func (h *TaskHandler) DeleteTask(c *gin.Context) {
	id := c.Param("id")
	db := database.GetDB()

	var task models.TestTask
	if err := db.First(&task, id).Error; err != nil {
		utils.NotFound(c, "Task not found")
		return
	}

	if err := db.Delete(&task).Error; err != nil {
		utils.InternalServerError(c, "Failed to delete task")
		return
	}

	// 从调度器中移除任务
	if services.GlobalScheduler != nil {
		taskID, _ := strconv.ParseUint(id, 10, 32)
		services.GlobalScheduler.RemoveTask(uint(taskID))
	}

	utils.SuccessWithMessage(c, "Task deleted successfully", nil)
}

// RunTask 执行任务
// @Summary 执行任务
// @Description 手动执行任务
// @Tags 任务管理
// @Produce json
// @Security BearerAuth
// @Param id path int true "任务ID"
// @Success 200 {object} utils.Response
// @Failure 404 {object} utils.Response
// @Router /api/tasks/{id}/run [post]
func (h *TaskHandler) RunTask(c *gin.Context) {
	id := c.Param("id")
	db := database.GetDB()

	var task models.TestTask
	if err := db.First(&task, id).Error; err != nil {
		utils.NotFound(c, "Task not found")
		return
	}

	// 使用TaskService执行任务
	taskService := services.NewTaskService()
	go func() {
		// 异步执行任务
		result, err := taskService.ExecuteTask(task.ID)
		if err != nil {
			utils.LogError("Task execution failed", err)
			return
		}
		utils.LogInfo("Task execution completed", map[string]interface{}{
			"task_id": task.ID,
			"status":  result.Status,
		})
	}()
	
	utils.SuccessWithMessage(c, "Task execution started", gin.H{
		"task_id": task.ID,
		"status":  "running",
	})
}

// GetTaskReports 获取任务报告列表
// @Summary 获取任务报告列表
// @Description 获取指定任务的报告列表
// @Tags 任务管理
// @Produce json
// @Security BearerAuth
// @Param id path int true "任务ID"
// @Param page query int false "页码" default(1)
// @Param size query int false "每页数量" default(10)
// @Success 200 {object} utils.PageResponse{data=[]models.TaskReport}
// @Failure 404 {object} utils.Response
// @Router /api/tasks/{id}/reports [get]
func (h *TaskHandler) GetTaskReports(c *gin.Context) {
	taskID := c.Param("id")
	db := database.GetDB()

	// 检查任务是否存在
	var task models.TestTask
	if err := db.First(&task, taskID).Error; err != nil {
		utils.NotFound(c, "Task not found")
		return
	}

	// 获取分页参数
	page, _ := strconv.Atoi(c.DefaultQuery("page", "1"))
	size, _ := strconv.Atoi(c.DefaultQuery("size", "10"))

	if page < 1 {
		page = 1
	}
	if size < 1 || size > 100 {
		size = 10
	}

	offset := (page - 1) * size

	// 获取报告列表
	var total int64
	db.Model(&models.TaskReport{}).Where("task = ?", taskID).Count(&total)

	var reports []models.TaskReport
	if err := db.Where("task = ?", taskID).Offset(offset).Limit(size).Order("create_time DESC").Find(&reports).Error; err != nil {
		utils.InternalServerError(c, "Failed to fetch task reports")
		return
	}

	utils.PageSuccess(c, reports, total, page, size)
}
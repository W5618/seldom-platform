package services

import (
	"fmt"
	"time"

	"github.com/robfig/cron/v3"
	"seldom-platform/database"
	"seldom-platform/models"
	"seldom-platform/utils"
)

// SchedulerService 调度服务
type SchedulerService struct {
	cron   *cron.Cron
	logger *utils.Logger
	taskService *TaskService
}

// NewSchedulerService 创建调度服务实例
func NewSchedulerService() *SchedulerService {
	return &SchedulerService{
		cron:   cron.New(cron.WithSeconds()),
		logger: utils.GetLogger(),
		taskService: NewTaskService(),
	}
}

// Start 启动调度服务
func (s *SchedulerService) Start() error {
	// 加载所有启用的定时任务
	if err := s.loadScheduledTasks(); err != nil {
		return fmt.Errorf("加载定时任务失败: %v", err)
	}

	// 启动cron调度器
	s.cron.Start()
	
	s.logger.LogInfo("SCHEDULER", "调度服务已启动", nil)
	return nil
}

// Stop 停止调度服务
func (s *SchedulerService) Stop() {
	s.cron.Stop()
	s.logger.LogInfo("SCHEDULER", "调度服务已停止", nil)
}

// loadScheduledTasks 加载定时任务
func (s *SchedulerService) loadScheduledTasks() error {
	db := database.GetDB()
	
	var tasks []models.TestTask
	// 修复查询条件：is_scheduled为布尔值，status为整数
	if err := db.Where("is_scheduled = ? AND is_delete = ?", true, false).Find(&tasks).Error; err != nil {
		return err
	}

	for _, task := range tasks {
		if err := s.addScheduledTask(task); err != nil {
			s.logger.LogError("SCHEDULER", fmt.Sprintf("添加定时任务失败: %v", err), map[string]interface{}{
				"task_id": task.ID,
				"task_name": task.Name,
			})
		}
	}

	s.logger.LogInfo("SCHEDULER", fmt.Sprintf("已加载 %d 个定时任务", len(tasks)), nil)
	return nil
}

// addScheduledTask 添加定时任务
func (s *SchedulerService) addScheduledTask(task models.TestTask) error {
	if task.CronExpression == "" {
		return fmt.Errorf("任务 %d 缺少cron表达式", task.ID)
	}

	// 验证cron表达式
	if !utils.IsValidCronExpression(task.CronExpression) {
		return fmt.Errorf("任务 %d 的cron表达式无效: %s", task.ID, task.CronExpression)
	}

	// 添加到cron调度器
	_, err := s.cron.AddFunc(task.CronExpression, func() {
		s.executeScheduledTask(task.ID)
	})

	if err != nil {
		return fmt.Errorf("添加cron任务失败: %v", err)
	}

	s.logger.LogInfo("SCHEDULER", fmt.Sprintf("已添加定时任务: %s", task.Name), map[string]interface{}{
		"task_id": task.ID,
		"cron_expression": task.CronExpression,
	})

	return nil
}

// executeScheduledTask 执行定时任务
func (s *SchedulerService) executeScheduledTask(taskID uint) {
	s.logger.LogInfo("SCHEDULER", fmt.Sprintf("开始执行定时任务: %d", taskID), map[string]interface{}{
		"task_id": taskID,
	})

	// 检查任务是否已在运行
	status, err := s.taskService.GetTaskStatus(taskID)
	if err != nil {
		s.logger.LogError("SCHEDULER", fmt.Sprintf("获取任务状态失败: %v", err), map[string]interface{}{
			"task_id": taskID,
		})
		return
	}

	if status == "running" {
		s.logger.LogInfo("SCHEDULER", fmt.Sprintf("任务 %d 已在运行中，跳过本次执行", taskID), map[string]interface{}{
			"task_id": taskID,
		})
		return
	}

	// 异步执行任务
	go func() {
		result, err := s.taskService.ExecuteTask(taskID)
		if err != nil {
			s.logger.LogError("SCHEDULER", fmt.Sprintf("定时任务执行失败: %v", err), map[string]interface{}{
				"task_id": taskID,
			})
		} else {
			s.logger.LogInfo("SCHEDULER", fmt.Sprintf("定时任务执行完成: %d", taskID), map[string]interface{}{
				"task_id": taskID,
				"status": result.Status,
				"duration": result.Duration.String(),
			})
		}
	}()
}

// AddTask 添加新的定时任务
func (s *SchedulerService) AddTask(taskID uint) error {
	db := database.GetDB()
	
	var task models.TestTask
	if err := db.First(&task, taskID).Error; err != nil {
		return fmt.Errorf("任务不存在: %v", err)
	}

	if !task.IsScheduled {
		return fmt.Errorf("任务未启用定时调度")
	}

	return s.addScheduledTask(task)
}

// RemoveTask 移除定时任务
func (s *SchedulerService) RemoveTask(taskID uint) error {
	// 注意：cron/v3 不支持直接通过ID移除任务
	// 这里需要重新加载所有任务
	s.cron.Stop()
	s.cron = cron.New(cron.WithSeconds())
	
	if err := s.loadScheduledTasks(); err != nil {
		return err
	}
	
	s.cron.Start()
	
	s.logger.LogInfo("SCHEDULER", fmt.Sprintf("已移除定时任务: %d", taskID), map[string]interface{}{
		"task_id": taskID,
	})
	
	return nil
}

// UpdateTask 更新定时任务
func (s *SchedulerService) UpdateTask(taskID uint) error {
	// 重新加载任务（简单实现）
	return s.RemoveTask(taskID)
}

// GetScheduledTasks 获取所有定时任务
func (s *SchedulerService) GetScheduledTasks() ([]models.TestTask, error) {
	db := database.GetDB()
	
	var tasks []models.TestTask
	if err := db.Where("is_scheduled = ?", true).Find(&tasks).Error; err != nil {
		return nil, err
	}

	return tasks, nil
}

// GetNextRunTime 获取任务下次执行时间
func (s *SchedulerService) GetNextRunTime(cronExpression string) (time.Time, error) {
	parser := cron.NewParser(cron.Second | cron.Minute | cron.Hour | cron.Dom | cron.Month | cron.Dow)
	schedule, err := parser.Parse(cronExpression)
	if err != nil {
		return time.Time{}, err
	}

	return schedule.Next(time.Now()), nil
}

// ValidateCronExpression 验证cron表达式
func (s *SchedulerService) ValidateCronExpression(expression string) error {
	parser := cron.NewParser(cron.Second | cron.Minute | cron.Hour | cron.Dom | cron.Month | cron.Dow)
	_, err := parser.Parse(expression)
	return err
}

// GetRunningTasks 获取正在运行的任务
func (s *SchedulerService) GetRunningTasks() ([]models.TestTask, error) {
	db := database.GetDB()
	
	var tasks []models.TestTask
	if err := db.Where("status = ?", "running").Find(&tasks).Error; err != nil {
		return nil, err
	}

	return tasks, nil
}

// GetTaskHistory 获取任务执行历史
func (s *SchedulerService) GetTaskHistory(taskID uint, limit int) ([]models.TaskReport, error) {
	db := database.GetDB()
	
	var reports []models.TaskReport
	query := db.Where("task_id = ?", taskID).Order("created_at DESC")
	
	if limit > 0 {
		query = query.Limit(limit)
	}
	
	if err := query.Find(&reports).Error; err != nil {
		return nil, err
	}

	return reports, nil
}

// GetTaskStatistics 获取任务统计信息
func (s *SchedulerService) GetTaskStatistics(taskID uint, days int) (map[string]interface{}, error) {
	db := database.GetDB()
	
	// 计算时间范围
	endTime := time.Now()
	startTime := endTime.AddDate(0, 0, -days)
	
	// 获取执行次数
	var totalRuns int64
	db.Model(&models.TaskReport{}).Where("task_id = ? AND created_at BETWEEN ? AND ?", taskID, startTime, endTime).Count(&totalRuns)
	
	// 获取成功次数
	var successRuns int64
	db.Model(&models.TaskReport{}).Where("task_id = ? AND status = ? AND created_at BETWEEN ? AND ?", taskID, "success", startTime, endTime).Count(&successRuns)
	
	// 获取失败次数
	var failedRuns int64
	db.Model(&models.TaskReport{}).Where("task_id = ? AND status = ? AND created_at BETWEEN ? AND ?", taskID, "failed", startTime, endTime).Count(&failedRuns)
	
	// 计算成功率
	var successRate float64
	if totalRuns > 0 {
		successRate = float64(successRuns) / float64(totalRuns) * 100
	}
	
	// 获取平均执行时间
	var avgDuration float64
	var reports []models.TaskReport
	db.Where("task_id = ? AND created_at BETWEEN ? AND ?", taskID, startTime, endTime).Find(&reports)
	
	if len(reports) > 0 {
		var totalDuration time.Duration
		for _, report := range reports {
			// TaskReport中的RunTime是字符串格式，需要解析
			if duration, err := time.ParseDuration(report.RunTime + "s"); err == nil {
				totalDuration += duration
			}
		}
		avgDuration = totalDuration.Seconds() / float64(len(reports))
	}
	
	return map[string]interface{}{
		"total_runs":    totalRuns,
		"success_runs":  successRuns,
		"failed_runs":   failedRuns,
		"success_rate":  successRate,
		"avg_duration":  avgDuration,
		"period_days":   days,
	}, nil
}
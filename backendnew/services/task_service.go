package services

import (
	"encoding/json"
	"fmt"
	"os/exec"
	"time"

	"seldom-platform/database"
	"seldom-platform/models"
	"seldom-platform/utils"
)

// TaskService 任务服务
type TaskService struct {
	logger *utils.Logger
}

// NewTaskService 创建任务服务实例
func NewTaskService() *TaskService {
	return &TaskService{
		logger: utils.GetLogger(),
	}
}

// TaskExecutionResult 任务执行结果
type TaskExecutionResult struct {
	TaskID    uint                   `json:"task_id"`
	Status    string                 `json:"status"`
	StartTime time.Time              `json:"start_time"`
	EndTime   time.Time              `json:"end_time"`
	Duration  time.Duration          `json:"duration"`
	Results   []CaseExecutionResult  `json:"results"`
	Summary   TaskExecutionSummary   `json:"summary"`
	Error     string                 `json:"error,omitempty"`
}

// CaseExecutionResult 用例执行结果
type CaseExecutionResult struct {
	CaseID      uint      `json:"case_id"`
	CaseName    string    `json:"case_name"`
	Status      string    `json:"status"`
	StartTime   time.Time `json:"start_time"`
	EndTime     time.Time `json:"end_time"`
	Duration    time.Duration `json:"duration"`
	ErrorMsg    string    `json:"error_msg,omitempty"`
	Screenshots []string  `json:"screenshots,omitempty"`
	Logs        []string  `json:"logs,omitempty"`
}

// TaskExecutionSummary 任务执行摘要
type TaskExecutionSummary struct {
	TotalCases  int `json:"total_cases"`
	PassedCases int `json:"passed_cases"`
	FailedCases int `json:"failed_cases"`
	SkippedCases int `json:"skipped_cases"`
	PassRate    float64 `json:"pass_rate"`
}

// ExecuteTask 执行任务
func (s *TaskService) ExecuteTask(taskID uint) (*TaskExecutionResult, error) {
	db := database.GetDB()
	
	// 获取任务信息
	var task models.TestTask
	if err := db.First(&task, taskID).Error; err != nil {
		return nil, fmt.Errorf("任务不存在: %v", err)
	}

	// 更新任务状态为运行中
	task.Status = 1 // 1表示执行中
	db.Save(&task)

	result := &TaskExecutionResult{
		TaskID:    taskID,
		Status:    "running",
		StartTime: time.Now(),
		Results:   make([]CaseExecutionResult, 0),
	}

	// 记录开始执行
	s.logger.LogInfo("TASK_EXECUTION", fmt.Sprintf("开始执行任务: %d", taskID), map[string]interface{}{
		"task_id": taskID,
		"task_name": task.Name,
	})

	// 获取任务关联的测试用例
	var relevances []models.TaskCaseRelevance
	if err := db.Where("task_id = ?", taskID).Find(&relevances).Error; err != nil {
		result.Status = "failed"
		result.Error = fmt.Sprintf("获取任务用例失败: %v", err)
		s.updateTaskStatus(&task, "failed", result.Error)
		return result, err
	}

	// 执行每个测试用例
	for _, relevance := range relevances {
		// 根据CaseHash查找测试用例
		var testCase models.TestCase
		if err := db.Where("case_hash = ?", relevance.CaseHash).First(&testCase).Error; err != nil {
			// 如果找不到用例，创建一个失败的结果
			caseResult := CaseExecutionResult{
				CaseID:    0, // 用例不存在
				CaseName:  relevance.CaseHash,
				Status:    "failed",
				StartTime: time.Now(),
				EndTime:   time.Now(),
				Duration:  0,
				ErrorMsg:  fmt.Sprintf("用例不存在: %v", err),
			}
			result.Results = append(result.Results, caseResult)
			continue
		}
		
		caseResult := s.executeSingleCase(testCase.ID)
		result.Results = append(result.Results, caseResult)
		
		// 保存用例执行结果
		s.saveCaseResult(taskID, caseResult)
	}

	// 计算执行结果
	result.EndTime = time.Now()
	result.Duration = result.EndTime.Sub(result.StartTime)
	result.Summary = s.calculateSummary(result.Results)
	
	// 确定任务最终状态
	if result.Summary.FailedCases > 0 {
		result.Status = "failed"
	} else {
		result.Status = "success"
	}

	// 更新任务状态
	s.updateTaskStatus(&task, result.Status, "")

	// 保存任务报告
	s.saveTaskReport(taskID, result)

	// 记录执行完成
	s.logger.LogInfo("TASK_EXECUTION", fmt.Sprintf("任务执行完成: %d", taskID), map[string]interface{}{
		"task_id": taskID,
		"status": result.Status,
		"duration": result.Duration.String(),
		"total_cases": result.Summary.TotalCases,
		"passed_cases": result.Summary.PassedCases,
		"failed_cases": result.Summary.FailedCases,
	})

	return result, nil
}

// executeSingleCase 执行单个测试用例
func (s *TaskService) executeSingleCase(caseID uint) CaseExecutionResult {
	db := database.GetDB()
	
	result := CaseExecutionResult{
		CaseID:    caseID,
		StartTime: time.Now(),
		Status:    "running",
	}

	// 获取用例信息
	var testCase models.TestCase
	if err := db.First(&testCase, caseID).Error; err != nil {
		result.Status = "failed"
		result.ErrorMsg = fmt.Sprintf("获取用例失败: %v", err)
		result.EndTime = time.Now()
		result.Duration = result.EndTime.Sub(result.StartTime)
		return result
	}

	result.CaseName = testCase.CaseName

	// 解析用例数据
	var caseData map[string]interface{}
	if err := json.Unmarshal([]byte(testCase.CaseDoc), &caseData); err != nil {
		result.Status = "failed"
		result.ErrorMsg = fmt.Sprintf("解析用例数据失败: %v", err)
		result.EndTime = time.Now()
		result.Duration = result.EndTime.Sub(result.StartTime)
		return result
	}

	// 执行用例（这里简化处理，实际应该根据用例类型执行不同的逻辑）
	if err := s.runTestCase(caseData); err != nil {
		result.Status = "failed"
		result.ErrorMsg = err.Error()
	} else {
		result.Status = "passed"
	}

	result.EndTime = time.Now()
	result.Duration = result.EndTime.Sub(result.StartTime)

	return result
}

// runTestCase 运行测试用例
func (s *TaskService) runTestCase(caseData map[string]interface{}) error {
	// 这里是简化的实现，实际应该根据用例类型执行不同的测试逻辑
	// 例如：HTTP接口测试、UI自动化测试等
	
	// 模拟执行时间
	time.Sleep(time.Millisecond * 100)
	
	// 模拟随机成功/失败（实际应该根据真实测试结果）
	// 这里总是返回成功，实际实现中应该执行真正的测试逻辑
	return nil
}

// executeSeldomTest 执行Seldom测试
func (s *TaskService) executeSeldomTest(scriptPath string, env map[string]string) error {
	// 构建命令
	cmd := exec.Command("python", "-m", "seldom", scriptPath)
	
	// 设置环境变量
	if env != nil {
		for key, value := range env {
			cmd.Env = append(cmd.Env, fmt.Sprintf("%s=%s", key, value))
		}
	}

	// 执行命令
	output, err := cmd.CombinedOutput()
	if err != nil {
		return fmt.Errorf("执行Seldom测试失败: %v, 输出: %s", err, string(output))
	}

	return nil
}

// saveCaseResult 保存用例执行结果
func (s *TaskService) saveCaseResult(taskID uint, result CaseExecutionResult) {
	db := database.GetDB()
	// 保存用例执行结果到数据库
	caseResult := models.CaseResult{
		CaseID:     result.CaseID,
		Name:       result.CaseName,
		Report:     fmt.Sprintf("Status: %s, Duration: %v", result.Status, result.Duration),
		RunTime:    result.Duration.Seconds(),
	}

	if err := db.Create(&caseResult).Error; err != nil {
		s.logger.LogError("SAVE_CASE_RESULT", fmt.Sprintf("保存用例结果失败: %v", err), map[string]interface{}{
			"task_id": taskID,
			"case_id": result.CaseID,
		})
	}
}

// saveTaskReport 保存任务报告
func (s *TaskService) saveTaskReport(taskID uint, result *TaskExecutionResult) {
	db := database.GetDB()
	// 创建任务报告
	report := models.TaskReport{
		TaskID:  taskID,
		Name:    fmt.Sprintf("Task %d Execution Report", taskID),
		Report:  fmt.Sprintf("Task executed with status: %s", result.Status),
		Passed:  result.Summary.PassedCases,
		Error:   result.Summary.FailedCases,
		Failure: result.Summary.FailedCases,
		Skipped: result.Summary.SkippedCases,
		Tests:   result.Summary.TotalCases,
		RunTime: fmt.Sprintf("%.2fs", result.Duration.Seconds()),
	}

	if err := db.Create(&report).Error; err != nil {
		s.logger.LogError("SAVE_TASK_REPORT", fmt.Sprintf("保存任务报告失败: %v", err), map[string]interface{}{
			"task_id": taskID,
		})
	}
}

// updateTaskStatus 更新任务状态
func (s *TaskService) updateTaskStatus(task *models.TestTask, status, errorMsg string) {
	db := database.GetDB()
	
	// 将字符串状态转换为整数
	var statusInt int
	switch status {
	case "running":
		statusInt = 1
	case "success":
		statusInt = 2
	case "failed":
		statusInt = 2 // 失败也算已执行
	default:
		statusInt = 0 // 未执行
	}
	
	task.Status = statusInt
	// 注意：TestTask模型中没有EndTime和ErrorMsg字段，这里移除相关代码
	
	db.Save(task)
}

// calculateSummary 计算执行摘要
func (s *TaskService) calculateSummary(results []CaseExecutionResult) TaskExecutionSummary {
	summary := TaskExecutionSummary{
		TotalCases: len(results),
	}

	for _, result := range results {
		switch result.Status {
		case "passed":
			summary.PassedCases++
		case "failed":
			summary.FailedCases++
		case "skipped":
			summary.SkippedCases++
		}
	}

	if summary.TotalCases > 0 {
		summary.PassRate = float64(summary.PassedCases) / float64(summary.TotalCases) * 100
	}

	return summary
}

// GetTaskStatus 获取任务状态
func (s *TaskService) GetTaskStatus(taskID uint) (string, error) {
	db := database.GetDB()
	
	var task models.TestTask
	if err := db.First(&task, taskID).Error; err != nil {
		return "", fmt.Errorf("任务不存在: %v", err)
	}

	utils.LogInfo("Getting task status", map[string]interface{}{
		"task_id":   task.ID,
		"task_name": task.Name,
	})

	return fmt.Sprintf("%d", task.Status), nil
}

// StopTask 停止任务执行
func (s *TaskService) StopTask(taskID uint) error {
	db := database.GetDB()
	
	var task models.TestTask
	if err := db.First(&task, taskID).Error; err != nil {
		return fmt.Errorf("任务不存在: %v", err)
	}

	if task.Status != 1 { // 1表示运行中
		return fmt.Errorf("任务未在运行中")
	}

	// 更新任务状态为已停止
	task.Status = 2 // 2表示已停止
	task.UpdateTime = time.Now() // 使用UpdateTime而不是EndTime
	
	if err := db.Save(&task).Error; err != nil {
		return fmt.Errorf("更新任务状态失败: %v", err)
	}

	s.logger.LogInfo("TASK_STOP", fmt.Sprintf("任务已停止: %d", taskID), map[string]interface{}{
		"task_id": taskID,
	})

	return nil
}
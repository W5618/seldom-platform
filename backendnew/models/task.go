package models

import (
	"time"

	"github.com/jinzhu/gorm"
)

// TestTask 测试任务
type TestTask struct {
	ID             uint      `gorm:"primary_key" json:"id"`
	ProjectID      uint      `gorm:"not null" json:"project_id"`                                       // 项目ID
	Project        Project   `gorm:"foreignkey:ProjectID;constraint:OnDelete:CASCADE" json:"project"` // 项目关联
	Name           string    `gorm:"size:200;not null;default:''" json:"name"`                         // 任务名
	Status         int       `gorm:"default:0" json:"status"`                                          // 状态 0未执行、1执行中、2已执行
	EnvID          *uint     `json:"env_id"`                                                           // 环境ID
	TeamID         *uint     `json:"team_id"`                                                          // 团队ID
	Email          string    `gorm:"size:100" json:"email"`                                            // 发送告警邮箱
	Timed          string    `gorm:"size:500;default:''" json:"timed"`                                 // 定时任务
	IsScheduled    bool      `gorm:"default:false" json:"is_scheduled"`                                // 是否启用定时调度
	CronExpression string    `gorm:"size:200;default:''" json:"cron_expression"`                       // Cron表达式
	ExecuteCount   int       `gorm:"default:0" json:"execute_count"`                                   // 执行次数
	IsDelete       bool      `gorm:"default:false" json:"is_delete"`                                   // 删除
	CreateTime     time.Time `gorm:"autoCreateTime" json:"create_time"`                                // 创建时间
	UpdateTime     time.Time `gorm:"autoUpdateTime" json:"update_time"`                                // 更新时间
}

// TableName 指定表名
func (TestTask) TableName() string {
	return "app_task_testtask"
}

// TaskCaseRelevance 任务用例关联表
type TaskCaseRelevance struct {
	ID         uint      `gorm:"primary_key" json:"id"`
	TaskID     uint      `gorm:"not null" json:"task_id"`                                           // 任务ID
	Task       TestTask  `gorm:"foreignkey:TaskID;constraint:OnDelete:CASCADE" json:"task"`        // 任务关联
	CaseHash   string    `gorm:"size:200;not null" json:"case_hash"`                                // 用例hash
	CreateTime time.Time `gorm:"autoCreateTime" json:"create_time"`                                 // 创建时间
}

// TableName 指定表名
func (TaskCaseRelevance) TableName() string {
	return "app_task_taskcaserelevance"
}

// TaskReport 任务报告
type TaskReport struct {
	ID         uint      `gorm:"primary_key" json:"id"`
	TaskID     uint      `gorm:"not null" json:"task_id"`                                           // 任务ID
	Task       TestTask  `gorm:"foreignkey:TaskID;constraint:OnDelete:CASCADE" json:"task"`        // 任务关联
	Name       string    `gorm:"size:500;not null;default:''" json:"name"`                          // 名称
	Report     string    `gorm:"type:text;default:''" json:"report"`                                // 报告内容
	Passed     int       `gorm:"default:0" json:"passed"`                                           // 通过用例
	Error      int       `gorm:"default:0" json:"error"`                                            // 错误用例
	Failure    int       `gorm:"default:0" json:"failure"`                                          // 失败用例
	Skipped    int       `gorm:"default:0" json:"skipped"`                                          // 跳过用例
	Tests      int       `gorm:"default:0" json:"tests"`                                            // 总用例数
	RunTime    string    `gorm:"size:100;default:'0'" json:"run_time"`                              // 运行时长
	CreateTime time.Time `gorm:"autoCreateTime" json:"create_time"`                                 // 创建时间
}

// TableName 指定表名
func (TaskReport) TableName() string {
	return "app_task_taskreport"
}

// ReportDetails 报告详情
type ReportDetails struct {
	ID             uint       `gorm:"primary_key" json:"id"`
	ResultID       uint       `gorm:"not null" json:"result_id"`                                        // 报告ID
	Result         TaskReport `gorm:"foreignkey:ResultID;constraint:OnDelete:CASCADE" json:"result"`    // 报告关联
	Name           string     `gorm:"size:500;not null;default:''" json:"name"`                          // 名称
	ClassName      string     `gorm:"size:200;not null;default:''" json:"class_name"`                    // 类名
	Status         string     `gorm:"size:20;not null;default:''" json:"status"`                         // 状态
	Time           string     `gorm:"size:100;not null;default:''" json:"time"`                          // 时间
	FailureMessage string     `gorm:"type:text;default:''" json:"failure_message"`                       // 失败信息
	ErrorOut       string     `gorm:"type:text;default:''" json:"error_out"`                             // 用例错误
	SkippedMessage string     `gorm:"type:text;default:''" json:"skipped_message"`                       // 跳过信息
	CreateTime     time.Time  `gorm:"autoCreateTime" json:"create_time"`                                 // 创建时间
}

// TableName 指定表名
func (ReportDetails) TableName() string {
	return "app_task_reportdetails"
}

// BeforeCreate GORM钩子，创建前执行
func (t *TestTask) BeforeCreate(scope *gorm.Scope) error {
	now := time.Now()
	scope.SetColumn("CreateTime", now)
	scope.SetColumn("UpdateTime", now)
	return nil
}

// BeforeUpdate GORM钩子，更新前执行
func (t *TestTask) BeforeUpdate(scope *gorm.Scope) error {
	scope.SetColumn("UpdateTime", time.Now())
	return nil
}

// BeforeCreate GORM钩子，创建前执行
func (t *TaskCaseRelevance) BeforeCreate(scope *gorm.Scope) error {
	scope.SetColumn("CreateTime", time.Now())
	return nil
}

// BeforeCreate GORM钩子，创建前执行
func (t *TaskReport) BeforeCreate(scope *gorm.Scope) error {
	scope.SetColumn("CreateTime", time.Now())
	return nil
}

// BeforeCreate GORM钩子，创建前执行
func (r *ReportDetails) BeforeCreate(scope *gorm.Scope) error {
	scope.SetColumn("CreateTime", time.Now())
	return nil
}
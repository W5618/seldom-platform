package models

import (
	"time"

	"github.com/jinzhu/gorm"
)

// TestCaseTemp 测试用例备份表
type TestCaseTemp struct {
	ID         uint      `gorm:"primary_key" json:"id"`
	ProjectID  uint      `gorm:"not null" json:"project_id"`                                        // 项目ID
	Project    Project   `gorm:"foreignkey:ProjectID;constraint:OnDelete:CASCADE" json:"project"`  // 项目关联
	FileName   string    `gorm:"size:500;not null;default:''" json:"file_name"`                     // 文件名
	ClassName  string    `gorm:"size:200;not null;default:''" json:"class_name"`                    // 类名
	ClassDoc   string    `gorm:"type:text;default:''" json:"class_doc"`                             // 类描述
	CaseName   string    `gorm:"size:200;not null;default:''" json:"case_name"`                     // 方法名
	CaseDoc    string    `gorm:"type:text;default:''" json:"case_doc"`                              // 方法描述
	Label      string    `gorm:"type:text;default:''" json:"label"`                                 // 用例标签
	CaseHash   string    `gorm:"size:200;not null;default:''" json:"case_hash"`                     // 用例hash
	CreateTime time.Time `gorm:"autoCreateTime" json:"create_time"`                                 // 创建时间
}

// TableName 指定表名
func (TestCaseTemp) TableName() string {
	return "app_case_testcasetemp"
}

// TestCase 测试类&用例
type TestCase struct {
	ID         uint      `gorm:"primary_key" json:"id"`
	ProjectID  uint      `gorm:"not null" json:"project_id"`                                        // 项目ID
	Project    Project   `gorm:"foreignkey:ProjectID;constraint:OnDelete:CASCADE" json:"project"`  // 项目关联
	FileName   string    `gorm:"size:500;not null;default:''" json:"file_name"`                     // 文件名
	ClassName  string    `gorm:"size:200;not null;default:''" json:"class_name"`                    // 类名
	ClassDoc   string    `gorm:"type:text;default:''" json:"class_doc"`                             // 类描述
	CaseName   string    `gorm:"size:200;not null;default:''" json:"case_name"`                     // 方法名
	CaseDoc    string    `gorm:"type:text;default:''" json:"case_doc"`                              // 方法描述
	Label      string    `gorm:"type:text;default:''" json:"label"`                                 // 用例标签
	Status     int       `gorm:"default:0" json:"status"`                                           // 状态 0未执行、1执行中、2已执行
	CaseHash   string    `gorm:"size:200;not null;default:''" json:"case_hash"`                     // 用例hash
	CreateTime time.Time `gorm:"autoCreateTime" json:"create_time"`                                 // 创建时间
	UpdateTime time.Time `gorm:"autoUpdateTime" json:"update_time"`                                 // 更新时间
}

// TableName 指定表名
func (TestCase) TableName() string {
	return "app_case_testcase"
}

// CaseResult 测试用例保存结果
type CaseResult struct {
	ID         uint      `gorm:"primary_key" json:"id"`
	CaseID     uint      `gorm:"not null" json:"case_id"`                                           // 用例ID
	Case       TestCase  `gorm:"foreignkey:CaseID;constraint:OnDelete:CASCADE" json:"case"`        // 用例关联
	Name       string    `gorm:"size:100;not null;default:''" json:"name"`                          // 名称
	Report     string    `gorm:"type:text;default:''" json:"report"`                                // 报告内容
	Passed     int       `gorm:"default:0" json:"passed"`                                           // 通过用例
	Error      int       `gorm:"default:0" json:"error"`                                            // 错误用例
	Failure    int       `gorm:"default:0" json:"failure"`                                          // 失败用例
	Skipped    int       `gorm:"default:0" json:"skipped"`                                          // 跳过用例
	Tests      int       `gorm:"default:0" json:"tests"`                                            // 总用例数
	SystemOut  string    `gorm:"type:text;default:''" json:"system_out"`                            // 日志
	RunTime    float64   `gorm:"default:0" json:"run_time"`                                         // 运行时长
	CreateTime time.Time `gorm:"autoCreateTime" json:"create_time"`                                 // 创建时间
}

// TableName 指定表名
func (CaseResult) TableName() string {
	return "app_case_caseresult"
}

// BeforeCreate GORM钩子，创建前执行
func (t *TestCaseTemp) BeforeCreate(scope *gorm.Scope) error {
	scope.SetColumn("CreateTime", time.Now())
	return nil
}

// BeforeCreate GORM钩子，创建前执行
func (t *TestCase) BeforeCreate(scope *gorm.Scope) error {
	now := time.Now()
	scope.SetColumn("CreateTime", now)
	scope.SetColumn("UpdateTime", now)
	return nil
}

// BeforeUpdate GORM钩子，更新前执行
func (t *TestCase) BeforeUpdate(scope *gorm.Scope) error {
	scope.SetColumn("UpdateTime", time.Now())
	return nil
}

// BeforeCreate GORM钩子，创建前执行
func (c *CaseResult) BeforeCreate(scope *gorm.Scope) error {
	scope.SetColumn("CreateTime", time.Now())
	return nil
}
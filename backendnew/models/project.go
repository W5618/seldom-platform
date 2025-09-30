package models

import (
	"time"

	"github.com/jinzhu/gorm"
)

// Project 项目表
type Project struct {
	ID         uint      `gorm:"primary_key" json:"id"`
	Name       string    `gorm:"size:50;not null" json:"name" binding:"required"`                    // 名称
	Address    string    `gorm:"size:200;not null" json:"address" binding:"required"`               // 项目地址
	CaseDir    string    `gorm:"size:200;default:'test_dir'" json:"case_dir"`                       // 用例目录
	IsDelete   bool      `gorm:"default:false" json:"is_delete"`                                    // 删除
	CreateTime time.Time `gorm:"autoCreateTime" json:"create_time"`                                 // 创建时间
	UpdateTime time.Time `gorm:"autoUpdateTime" json:"update_time"`                                 // 更新时间
	CoverName  string    `gorm:"size:64;default:''" json:"cover_name"`                              // 封面名称
	PathName   string    `gorm:"size:64;default:''" json:"path_name"`                               // 封面路径名称
	TestNum    int       `gorm:"default:0" json:"test_num"`                                         // 测试文件数
	IsClone    int       `gorm:"default:0" json:"is_clone"`                                         // 克隆
	RunVersion string    `gorm:"size:200;default:''" json:"run_version"`                            // 当前运行版本（蓝绿运行）
}

// TableName 指定表名
func (Project) TableName() string {
	return "app_project_project"
}

// Env 环境管理
type Env struct {
	ID           uint      `gorm:"primary_key" json:"id"`
	Name         string    `gorm:"size:50;not null" json:"name" binding:"required"`        // 名称
	TestType     string    `gorm:"size:20;default:'http'" json:"test_type"`                // 测试类型
	Env          string    `gorm:"size:50;default:''" json:"env"`                          // 环境值
	Rerun        int       `gorm:"default:0" json:"rerun"`                                 // 重跑次数
	IsClearCache bool      `gorm:"default:false" json:"is_clear_cache"`                    // 是否清除缓存
	Browser      string    `gorm:"size:20;default:''" json:"browser"`                      // 浏览器
	BaseURL      string    `gorm:"size:200;default:''" json:"base_url"`                    // URL
	Remote       string    `gorm:"size:200;default:''" json:"remote"`                      // remote
	AppServer    string    `gorm:"size:100;default:''" json:"app_server"`                  // APP服务
	AppInfo      string    `gorm:"size:1000;default:'{}'" json:"app_info"`                 // APP信息
	IsDelete     bool      `gorm:"default:false" json:"is_delete"`                         // 删除
	CreateTime   time.Time `gorm:"autoCreateTime" json:"create_time"`                      // 创建时间
	UpdateTime   time.Time `gorm:"autoUpdateTime" json:"update_time"`                      // 更新时间
}

// TableName 指定表名
func (Env) TableName() string {
	return "app_project_env"
}

// BeforeCreate GORM钩子，创建前执行
func (p *Project) BeforeCreate(scope *gorm.Scope) error {
	now := time.Now()
	scope.SetColumn("CreateTime", now)
	scope.SetColumn("UpdateTime", now)
	return nil
}

// BeforeUpdate GORM钩子，更新前执行
func (p *Project) BeforeUpdate(scope *gorm.Scope) error {
	scope.SetColumn("UpdateTime", time.Now())
	return nil
}

// BeforeCreate GORM钩子，创建前执行
func (e *Env) BeforeCreate(scope *gorm.Scope) error {
	now := time.Now()
	scope.SetColumn("CreateTime", now)
	scope.SetColumn("UpdateTime", now)
	return nil
}

// BeforeUpdate GORM钩子，更新前执行
func (e *Env) BeforeUpdate(scope *gorm.Scope) error {
	scope.SetColumn("UpdateTime", time.Now())
	return nil
}
package models

import (
	"time"

	"github.com/jinzhu/gorm"
)

// Team 团队表
type Team struct {
	ID         uint      `gorm:"primary_key" json:"id"`
	Name       string    `gorm:"size:200;not null" json:"name" binding:"required"`   // 团队名
	Email      string    `gorm:"type:text;default:''" json:"email"`                  // 团队邮箱
	IsDelete   bool      `gorm:"default:false" json:"is_delete"`                     // 删除
	CreateTime time.Time `gorm:"autoCreateTime" json:"create_time"`                  // 创建时间
	UpdateTime time.Time `gorm:"autoUpdateTime" json:"update_time"`                  // 更新时间
}

// TableName 指定表名
func (Team) TableName() string {
	return "app_team_team"
}

// BeforeCreate GORM钩子，创建前执行
func (t *Team) BeforeCreate(scope *gorm.Scope) error {
	now := time.Now()
	scope.SetColumn("CreateTime", now)
	scope.SetColumn("UpdateTime", now)
	return nil
}

// BeforeUpdate GORM钩子，更新前执行
func (t *Team) BeforeUpdate(scope *gorm.Scope) error {
	scope.SetColumn("UpdateTime", time.Now())
	return nil
}
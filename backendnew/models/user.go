package models

import (
	"time"

	"github.com/jinzhu/gorm"
	"golang.org/x/crypto/bcrypt"
)

// User 用户表（对应Django的User模型）
type User struct {
	ID          uint      `gorm:"primary_key" json:"id"`
	Username    string    `gorm:"size:150;unique;not null" json:"username" binding:"required"`    // 用户名
	Email       string    `gorm:"size:254" json:"email"`                                           // 邮箱
	FirstName   string    `gorm:"size:150" json:"first_name"`                                      // 名
	LastName    string    `gorm:"size:150" json:"last_name"`                                       // 姓
	Password    string    `gorm:"size:128;not null" json:"-"`                                      // 密码（不返回给前端）
	IsStaff     bool      `gorm:"default:false" json:"is_staff"`                                   // 是否为员工
	IsActive    bool      `gorm:"default:true" json:"is_active"`                                   // 是否激活
	IsSuperuser bool      `gorm:"default:false" json:"is_superuser"`                               // 是否为超级用户
	DateJoined  time.Time `gorm:"autoCreateTime" json:"date_joined"`                               // 加入时间
	LastLogin   *time.Time `json:"last_login"`                                                     // 最后登录时间
}

// TableName 指定表名
func (User) TableName() string {
	return "auth_user"
}

// SetPassword 设置密码（加密）
func (u *User) SetPassword(password string) error {
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		return err
	}
	u.Password = string(hashedPassword)
	return nil
}

// CheckPassword 验证密码
func (u *User) CheckPassword(password string) bool {
	err := bcrypt.CompareHashAndPassword([]byte(u.Password), []byte(password))
	return err == nil
}

// BeforeCreate GORM钩子，创建前执行
func (u *User) BeforeCreate(scope *gorm.Scope) error {
	scope.SetColumn("DateJoined", time.Now())
	return nil
}

// GetFullName 获取全名
func (u *User) GetFullName() string {
	if u.FirstName != "" && u.LastName != "" {
		return u.FirstName + " " + u.LastName
	}
	return u.Username
}
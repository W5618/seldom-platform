package handlers

import (
	"seldom-platform/config"
	"seldom-platform/database"
	"seldom-platform/models"
	"seldom-platform/utils"
	"time"

	"github.com/gin-gonic/gin"
)

// AuthHandler 认证处理器
type AuthHandler struct {
	config *config.Config
}

// NewAuthHandler 创建认证处理器
func NewAuthHandler(cfg *config.Config) *AuthHandler {
	return &AuthHandler{config: cfg}
}

// LoginRequest 登录请求结构
type LoginRequest struct {
	Username string `json:"username" binding:"required"`
	Password string `json:"password" binding:"required"`
}

// RegisterRequest 注册请求结构
type RegisterRequest struct {
	Username  string `json:"username" binding:"required"`
	Password  string `json:"password" binding:"required"`
	Email     string `json:"email"`
	FirstName string `json:"first_name"`
	LastName  string `json:"last_name"`
}

// LoginResponse 登录响应结构
type LoginResponse struct {
	Token string      `json:"token"`
	User  models.User `json:"user"`
}

// Login 用户登录
// @Summary 用户登录
// @Description 用户登录接口
// @Tags 认证
// @Accept json
// @Produce json
// @Param login body LoginRequest true "登录信息"
// @Success 200 {object} utils.Response{data=LoginResponse}
// @Failure 400 {object} utils.Response
// @Failure 401 {object} utils.Response
// @Router /api/auth/login [post]
func (h *AuthHandler) Login(c *gin.Context) {
	var req LoginRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		utils.BadRequest(c, "Invalid request format")
		return
	}

	db := database.GetDB()
	var user models.User

	// 查找用户
	if err := db.Where("username = ?", req.Username).First(&user).Error; err != nil {
		utils.Unauthorized(c, "Invalid username or password")
		return
	}

	// 验证密码
	if !user.CheckPassword(req.Password) {
		utils.Unauthorized(c, "Invalid username or password")
		return
	}

	// 检查用户是否激活
	if !user.IsActive {
		utils.Unauthorized(c, "User account is disabled")
		return
	}

	// 更新最后登录时间
	now := time.Now()
	user.LastLogin = &now
	db.Save(&user)

	// 生成JWT token
	token, err := utils.GenerateJWT(user.ID, user.Username, h.config.JWT.Secret, h.config.JWT.Expire)
	if err != nil {
		utils.InternalServerError(c, "Failed to generate token")
		return
	}

	utils.Success(c, LoginResponse{
		Token: token,
		User:  user,
	})
}

// Register 用户注册
// @Summary 用户注册
// @Description 用户注册接口
// @Tags 认证
// @Accept json
// @Produce json
// @Param register body RegisterRequest true "注册信息"
// @Success 200 {object} utils.Response{data=models.User}
// @Failure 400 {object} utils.Response
// @Router /api/auth/register [post]
func (h *AuthHandler) Register(c *gin.Context) {
	var req RegisterRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		utils.BadRequest(c, "Invalid request format")
		return
	}

	db := database.GetDB()

	// 检查用户名是否已存在
	var existingUser models.User
	if err := db.Where("username = ?", req.Username).First(&existingUser).Error; err == nil {
		utils.BadRequest(c, "Username already exists")
		return
	}

	// 创建新用户
	user := models.User{
		Username:  req.Username,
		Email:     req.Email,
		FirstName: req.FirstName,
		LastName:  req.LastName,
		IsActive:  true,
	}

	// 设置密码
	if err := user.SetPassword(req.Password); err != nil {
		utils.InternalServerError(c, "Failed to encrypt password")
		return
	}

	// 保存用户
	if err := db.Create(&user).Error; err != nil {
		utils.InternalServerError(c, "Failed to create user")
		return
	}

	utils.SuccessWithMessage(c, "User registered successfully", user)
}

// GetProfile 获取用户信息
// @Summary 获取用户信息
// @Description 获取当前登录用户的信息
// @Tags 认证
// @Produce json
// @Security BearerAuth
// @Success 200 {object} utils.Response{data=models.User}
// @Failure 401 {object} utils.Response
// @Router /api/auth/profile [get]
func (h *AuthHandler) GetProfile(c *gin.Context) {
	userID, exists := c.Get("user_id")
	if !exists {
		utils.Unauthorized(c, "User not authenticated")
		return
	}

	db := database.GetDB()
	var user models.User

	if err := db.First(&user, userID).Error; err != nil {
		utils.NotFound(c, "User not found")
		return
	}

	utils.Success(c, user)
}

// UpdateProfile 更新用户信息
// @Summary 更新用户信息
// @Description 更新当前登录用户的信息
// @Tags 认证
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param profile body RegisterRequest true "用户信息"
// @Success 200 {object} utils.Response{data=models.User}
// @Failure 400 {object} utils.Response
// @Failure 401 {object} utils.Response
// @Router /api/auth/profile [put]
func (h *AuthHandler) UpdateProfile(c *gin.Context) {
	userID, exists := c.Get("user_id")
	if !exists {
		utils.Unauthorized(c, "User not authenticated")
		return
	}

	var req RegisterRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		utils.BadRequest(c, "Invalid request format")
		return
	}

	db := database.GetDB()
	var user models.User

	if err := db.First(&user, userID).Error; err != nil {
		utils.NotFound(c, "User not found")
		return
	}

	// 更新用户信息
	user.Email = req.Email
	user.FirstName = req.FirstName
	user.LastName = req.LastName

	// 如果提供了新密码，则更新密码
	if req.Password != "" {
		if err := user.SetPassword(req.Password); err != nil {
			utils.InternalServerError(c, "Failed to encrypt password")
			return
		}
	}

	if err := db.Save(&user).Error; err != nil {
		utils.InternalServerError(c, "Failed to update user")
		return
	}

	utils.SuccessWithMessage(c, "Profile updated successfully", user)
}
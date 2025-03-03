package handlers

import (
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/binginx/bqd_chat_log/models"
)

// CreateLogRequest 创建日志请求结构
// @Description 创建日志的请求参数
type CreateLogRequest struct {
	// 聊天文本内容
	Text string `json:"text" binding:"required" example:"这是一条聊天记录"`
	// 验证结果
	ValidationResult bool `json:"validation_result" example:"true"`
}

// LogResponse 日志响应结构
// @Description 日志操作的响应结果
type LogResponse struct {
	// 操作状态码
	Code int `json:"code" example:"200"`
	// 操作消息
	Message string `json:"message" example:"成功"`
	// 日志数据
	Data interface{} `json:"data,omitempty"`
}

// GetLogsRequest 获取日志请求结构
// @Description 获取日志的查询参数
type GetLogsRequest struct {
	// 开始日期（可选，默认为结束日期前7天）
	StartDate string `form:"start_date" example:"2023-01-01"`
	// 结束日期（可选，默认为当前日期）
	EndDate string `form:"end_date" example:"2023-01-07"`
}

// CreateLog godoc
// @Summary 创建聊天日志
// @Description 创建一条新的聊天日志记录
// @Tags logs
// @Accept json
// @Produce json
// @Param request body CreateLogRequest true "日志信息"
// @Success 200 {object} LogResponse
// @Failure 400 {object} LogResponse
// @Failure 500 {object} LogResponse
// @Router /logs [post]
func CreateLog(c *gin.Context) {
	var req CreateLogRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, LogResponse{
			Code:    http.StatusBadRequest,
			Message: "无效的请求参数: " + err.Error(),
		})
		return
	}

	// 创建日志对象
	log := &models.ChatLog{
		Text:             req.Text,
		ValidationResult: req.ValidationResult,
	}

	// 时间戳将在SaveLog函数中自动生成

	// 保存日志
	if err := models.SaveLog(log); err != nil {
		c.JSON(http.StatusInternalServerError, LogResponse{
			Code:    http.StatusInternalServerError,
			Message: "保存日志失败: " + err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, LogResponse{
		Code:    http.StatusOK,
		Message: "日志保存成功",
		Data:    log,
	})
}

// GetLogs godoc
// @Summary 获取聊天日志
// @Description 获取指定日期范围内的聊天日志
// @Tags logs
// @Accept json
// @Produce json
// @Param start_date query string false "开始日期 (YYYY-MM-DD)"
// @Param end_date query string false "结束日期 (YYYY-MM-DD)"
// @Success 200 {object} LogResponse
// @Failure 400 {object} LogResponse
// @Failure 500 {object} LogResponse
// @Router /logs [get]
func GetLogs(c *gin.Context) {
	var req GetLogsRequest
	if err := c.ShouldBindQuery(&req); err != nil {
		c.JSON(http.StatusBadRequest, LogResponse{
			Code:    http.StatusBadRequest,
			Message: "无效的查询参数: " + err.Error(),
		})
		return
	}

	// 解析日期
	var startDate, endDate time.Time
	var err error

	if req.StartDate != "" {
		startDate, err = time.Parse("2006-01-02", req.StartDate)
		if err != nil {
			c.JSON(http.StatusBadRequest, LogResponse{
				Code:    http.StatusBadRequest,
				Message: "无效的开始日期格式，应为YYYY-MM-DD",
			})
			return
		}
	}

	if req.EndDate != "" {
		endDate, err = time.Parse("2006-01-02", req.EndDate)
		if err != nil {
			c.JSON(http.StatusBadRequest, LogResponse{
				Code:    http.StatusBadRequest,
				Message: "无效的结束日期格式，应为YYYY-MM-DD",
			})
			return
		}
	}

	// 获取日志
	logs, err := models.GetLogs(startDate, endDate)
	if err != nil {
		c.JSON(http.StatusInternalServerError, LogResponse{
			Code:    http.StatusInternalServerError,
			Message: "获取日志失败: " + err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, LogResponse{
		Code:    http.StatusOK,
		Message: "获取日志成功",
		Data:    logs,
	})
}
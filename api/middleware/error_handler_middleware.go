package middleware

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/go-playground/validator/v10"
	"github.com/littlecheny/go-backend/domain"
	"github.com/littlecheny/go-backend/logger"
)

// ErrorHandlerMiddleware 统一错误处理中间件
func ErrorHandlerMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		c.Next()

		// 检查是否有错误
		if len(c.Errors) > 0 {
			err := c.Errors.Last()
			
			// 记录错误日志
			logger.Log.WithContext(c.Request.Context()).Errorf("Request error: %v", err.Err)

			// 处理验证错误
			if validationErrs, ok := err.Err.(validator.ValidationErrors); ok {
				errorMessages := make([]string, 0)
				for _, e := range validationErrs {
					errorMessages = append(errorMessages, formatValidationError(e))
				}
				c.JSON(http.StatusBadRequest, domain.ErrorResponse{
					Message: "Validation failed: " + errorMessages[0],
				})
				return
			}

			// 默认返回500错误
			c.JSON(http.StatusInternalServerError, domain.ErrorResponse{
				Message: "Internal server error",
			})
		}
	}
}

// formatValidationError 格式化验证错误消息
func formatValidationError(e validator.FieldError) string {
	switch e.Tag() {
	case "required":
		return e.Field() + " is required"
	case "email":
		return e.Field() + " must be a valid email address"
	case "min":
		return e.Field() + " must be at least " + e.Param() + " characters"
	case "max":
		return e.Field() + " must be at most " + e.Param() + " characters"
	default:
		return e.Field() + " is invalid"
	}
}


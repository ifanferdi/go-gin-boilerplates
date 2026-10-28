package middlewares

import (
	"errors"
	"fmt"
	"runtime/debug"
	"strings"

	"go-gin-boilerplate/internal/config"
	"go-gin-boilerplate/internal/constants"
	"go-gin-boilerplate/internal/helpers"

	"github.com/gin-gonic/gin"
	"github.com/go-playground/validator/v10"
)

// ErrorHandler menangkap error yang dilempar ke c.Error(err) atau panic (ekuivalen dengan error-handler.ts di POS-nextjs-express).
func ErrorHandler(cfg *config.Config) gin.HandlerFunc {
	return func(c *gin.Context) {
		defer func() {
			if rec := recover(); rec != nil {
				var err error
				switch v := rec.(type) {
				case error:
					err = v
				default:
					err = fmt.Errorf("%v", v)
				}
				respondError(c, cfg, err, string(debug.Stack()))
				c.Abort()
			}
		}()

		c.Next()

		if len(c.Errors) > 0 && !c.Writer.Written() {
			err := c.Errors.Last().Err
			respondError(c, cfg, err, "")
		}
	}
}

func respondError(c *gin.Context, cfg *config.Config, err error, stackTrace string) {
	statusCode := constants.StatusInternalServerError
	message := err.Error()
	var details interface{}

	// 1. Cek apakah error berasal dari go-playground/validator (ekuivalen ZodError di POS-nextjs-express)
	var validationErrs validator.ValidationErrors
	if errors.As(err, &validationErrs) && len(validationErrs) > 0 {
		firstErr := validationErrs[0]
		field := toCamelCase(firstErr.Field())
		message = fmt.Sprintf("Column %s: %s", field, formatValidationMessage(firstErr))
		statusCode = constants.StatusBadRequest
	} else if appErr, ok := helpers.IsAppError(err); ok {
		// 2. Cek apakah error adalah custom AppError
		statusCode = appErr.StatusCode
		message = appErr.Message
		details = appErr.Details
	}

	response := gin.H{
		"message": message,
	}

	if details != nil {
		response["details"] = details
	}

	if strings.EqualFold(cfg.App.Env, "development") && stackTrace != "" {
		response["stack"] = stackTrace
	}

	c.JSON(statusCode, response)
}

func formatValidationMessage(fe validator.FieldError) string {
	switch fe.Tag() {
	case "required":
		return fmt.Sprintf("%s cannot be empty.", fe.Field())
	case "min":
		return fmt.Sprintf("must be at least %s characters/value.", fe.Param())
	case "max":
		return fmt.Sprintf("must be at most %s characters/value.", fe.Param())
	case "oneof":
		return fmt.Sprintf("must be one of [%s].", fe.Param())
	default:
		return fmt.Sprintf("failed on '%s' validation.", fe.Tag())
	}
}

func toCamelCase(s string) string {
	if s == "" {
		return s
	}
	if len(s) == 1 {
		return strings.ToLower(s)
	}
	return strings.ToLower(s[:1]) + s[1:]
}

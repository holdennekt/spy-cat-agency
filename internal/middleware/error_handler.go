package middleware

import (
	"fmt"
	"net/http"
	"reflect"
	"spy-cat-agency/pkg/custerr"
	"strings"

	"github.com/gin-gonic/gin"
	"github.com/go-playground/validator/v10"
)

func messageForTag(fe validator.FieldError) string {
	switch fe.Tag() {
	case "required":
		return fmt.Sprintf("%s is required", fe.Namespace())
	case "min":
		switch fe.Type().Kind() {
		case reflect.String:
			return fmt.Sprintf("%s must be at least %s characters long", fe.Namespace(), fe.Param())
		case reflect.Int:
			return fmt.Sprintf("%s must be at least %s", fe.Namespace(), fe.Param())
		case reflect.Slice:
			return fmt.Sprintf("%s must have at least %s items", fe.Namespace(), fe.Param())
		}
	case "max":
		switch fe.Type().Kind() {
		case reflect.String:
			return fmt.Sprintf("%s must be at most %s characters long", fe.Namespace(), fe.Param())
		case reflect.Int:
			return fmt.Sprintf("%s must be at most %s", fe.Namespace(), fe.Param())
		case reflect.Slice:
			return fmt.Sprintf("%s must have at most %s items", fe.Namespace(), fe.Param())
		}
	}
	return fe.Error()
}

func ErrorHandler() gin.HandlerFunc {
	return func(c *gin.Context) {
		c.Next()

		if len(c.Errors) == 0 {
			return
		}

		err := c.Errors.Last()
		switch err.Type {
		case gin.ErrorTypeBind:
			validationErrors := err.Err.(validator.ValidationErrors)
			errs := make([]string, len(validationErrors))
			for i, fieldErr := range validationErrors {
				errs[i] = messageForTag(fieldErr)
			}
			c.JSON(http.StatusBadRequest, gin.H{"error": strings.Join(errs, ", ")})
		default:
			switch custErr := err.Err.(type) {
			case custerr.BadRequestErr:
				c.JSON(http.StatusBadRequest, gin.H{"error": custErr.Error()})
			case custerr.NotFoundErr:
				c.JSON(http.StatusNotFound, gin.H{"error": custErr.Error()})
			case custerr.ConflictErr:
				c.JSON(http.StatusConflict, gin.H{"error": custErr.Error()})
			default:
				c.JSON(http.StatusInternalServerError, gin.H{"error": "Internal server error"})
			}
		}
	}
}

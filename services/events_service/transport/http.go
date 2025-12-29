package transport

import (
	"context"
	"errors"
	"io"
	"net/http"
	"os"
	"strings"

	"github.com/gin-gonic/gin"
	"github.com/go-kit/kit/endpoint"
	"github.com/piyushsharma67/events_booking/services/events_service/logger"
	"github.com/piyushsharma67/events_booking/services/events_service/service"
	"github.com/piyushsharma67/events_booking/services/events_service/utils"
)

func GinHandler(
	e endpoint.Endpoint,
	newRequest func() interface{},
	logger logger.Logger,
) gin.HandlerFunc {
	return func(c *gin.Context) {
		request := newRequest()
		ctx := c.Request.Context()

		// Attach user_id and role from Gin context to Go Kit context
		if userID := c.GetString("user_id"); userID != "" {
			ctx = context.WithValue(ctx, "user_id", userID)
		}
		if role := c.GetString("role"); role != "" {
			ctx = context.WithValue(ctx, "role", role)
		}


		if err := c.ShouldBindJSON(request); err != nil {
			status := http.StatusBadRequest

			msg := err.Error()
			if errors.Is(err, io.EOF) {
				msg = "request body is required"
			}

			logger.
				WithContext(ctx).
				Warn("invalid request body",
					"status_code", status,
					"error", err.Error(),
					"method", c.Request.Method,
					"path", c.Request.URL.Path,
				)

			c.JSON(status, gin.H{"error": msg})
			return
		}

		logger.
			WithContext(ctx).
			Info("incoming request",
				"method", c.Request.Method,
				"path", c.Request.URL.Path,
			)

		resp, err := e(ctx, request)
		if err != nil {
			if errors.Is(err, context.DeadlineExceeded) {
				status := http.StatusRequestTimeout

				logger.
					WithContext(ctx).
					Error("request timed out",
						"status_code", status,
					)

				c.JSON(status, gin.H{"error": "request timed out"})
				return
			}

			status := http.StatusInternalServerError

			logger.
				WithContext(ctx).
				Error("endpoint failed",
					"status_code", status,
					"error", err.Error(),
				)

			c.JSON(status, gin.H{"error": err.Error()})
			return
		}

		status := http.StatusOK

		logger.
			WithContext(ctx).
			Info("request completed",
				"status_code", status,
				"method", c.Request.Method,
				"path", c.Request.URL.Path,
			)

		c.JSON(status, resp)
	}
}

func ValidateGinHandler(
	svc service.EventService,
	logger logger.Logger,
) gin.HandlerFunc {
	return func(c *gin.Context) {
		ctx := c.Request.Context()
		path := c.Request.URL.Path
		method := c.Request.Method

		authHeader := c.GetHeader("Authorization")
		if authHeader == "" {
			status := http.StatusUnauthorized

			logger.
				WithContext(ctx).
				Warn("missing authorization header",
					"status_code", status,
					"method", method,
					"path", path,
				)

			c.Status(status)
			return
		}

		token := strings.TrimPrefix(authHeader, "Bearer ")

		claims, err := utils.ValidateJWT(token, os.Getenv("JWT_SECRET"))
		if err != nil {
			status := http.StatusUnauthorized

			logger.
				WithContext(ctx).
				Warn("jwt validation failed",
					"status_code", status,
					"method", method,
					"path", path,
					"error", err.Error(),
				)

			c.Status(status)
			return
		}

		// 🔥 Headers consumed by NGINX
		c.Header("X-User-Id", claims.UserID)
		c.Header("X-User-Role", claims.Role)

		status := http.StatusOK

		logger.
			WithContext(ctx).
			Info("jwt validated successfully",
				"status_code", status,
				"user_id", claims.UserID,
				"role", claims.Role,
				"method", method,
				"path", path,
			)

		c.Status(status)
	}
}

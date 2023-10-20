package middleware

import (
	"context"
	"errors"
	"github.com/gin-gonic/gin"
	"time"
)

func TimeOutMiddleware(timeout time.Duration) func(ctx *gin.Context) {
	return func(ctx *gin.Context) {
		// wrap the request context with a timeout
		ctxTimeout, cancel := context.WithTimeout(ctx.Request.Context(), timeout)

		defer func() {
			// check if context timeout was reached
			if errors.Is(ctxTimeout.Err(), context.DeadlineExceeded) {

				// write response and abort the request
				ctx.JSON(304, gin.H{
					"error": "time out",
				})
				ctx.Abort()
			}

			//cancel to clear resources after finished
			cancel()
		}()

		// replace request with context wrapped request
		ctx.Request = ctx.Request.WithContext(ctx)
		ctx.Next()
		return
	}
}

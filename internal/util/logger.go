package util

import (
	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
)

var Logger *zap.Logger

func InitLogger() {
	var err error
	config := zap.NewProductionConfig()
	config.OutputPaths = []string{"stdout"}
	Logger, err = config.Build()
	if err != nil {
		panic(err)
	}
}

func GinLogger() gin.HandlerFunc {
	return func(c *gin.Context) {
		Logger.Info("request",
			zap.String("method", c.Request.Method),
			zap.String("path", c.Request.URL.Path),
			zap.Int("status", c.Writer.Status()),
		)
		c.Next()
	}
}

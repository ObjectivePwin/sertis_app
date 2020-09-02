package main

import (
	"log"
	"time"

	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
)

func main() {
	var logger *zap.Logger
	var err error

	logger, err = zap.NewDevelopment()
	if err != nil {
		log.Fatal("cannot initialize logger: ", err)
	}

	defer logger.Sync()

	zap.ReplaceGlobals(logger)

	router := setupRouter()
	router.Run(":8880")
}

func setupRouter() *gin.Engine {
	gin.SetMode(gin.ReleaseMode)

	router := gin.Default()
	router.Use(cors.New(cors.Config{
		AllowOrigins:     []string{"*"},
		AllowHeaders:     []string{"Origin"},
		ExposeHeaders:    []string{"Content-Length"},
		AllowCredentials: true,
		MaxAge:           12 * time.Hour,
	}))

	zap.L().Info("Start Listen at Port 5550")
	router.POST("/login", postLogin)

	return router
}

func postLogin(c *gin.Context) {

}

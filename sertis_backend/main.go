package main

import (
	"database/sql"
	"log"
	"net/http"
	miniblog "sertis_app/mini_blog"
	"sertis_app/model"
	"time"

	storage "sertis_app/storage"

	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
)

var db *sql.DB
var blog *miniblog.Blog

func main() {
	logger, err := zap.NewDevelopment()
	if err != nil {
		log.Fatal("cannot initialize logger: ", err)
	}
	defer logger.Sync()
	zap.ReplaceGlobals(logger)

	errorCreateDB := setupDB()
	if errorCreateDB != nil {
		log.Fatal("cannot connect DB: ", err)
		return
	}
	setupBlog()

	router := setupRouter()

	router.Run(":8880")
}

func setupBlog() {
	blog = miniblog.NewBlog(db)
}

func setupDB() error {
	var err error
	db, err = storage.CreateDBConnection()
	return err
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
	router.POST("/signup", postSignUp)
	router.POST("/signin", postSignIn)

	return router
}

func postSignUp(c *gin.Context) {
	creds := model.Credentials{}
	if err := c.ShouldBindJSON(&creds); err != nil {
		c.Status(http.StatusBadRequest)
		return
	}

	err := blog.CreateAccount(creds)
	if err != nil {
		c.JSON(http.StatusOK, createResponse("Already Have Account"))
	} else {
		c.JSON(http.StatusOK, createResponse(""))
	}
}

func postSignIn(c *gin.Context) {

}

func createResponse(err string) model.ResponseAPI {
	resAPI := model.ResponseAPI{}
	if err != "" {
		resAPI.Success = false
		resAPI.Error = err
	} else {
		resAPI.Success = true
	}

	return resAPI
}

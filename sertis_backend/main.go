package main

import (
	"database/sql"
	"log"
	"net/http"
	miniblog "sertis_app/mini_blog"
	"sertis_app/model"
	"strconv"
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
	gin.SetMode(gin.DebugMode)

	router := gin.Default()
	router.Use(cors.New(cors.Config{
		AllowOrigins:     []string{"*"},
		AllowHeaders:     []string{"Origin"},
		ExposeHeaders:    []string{"Content-Length"},
		AllowCredentials: true,
		MaxAge:           12 * time.Hour,
	}))

	zap.L().Info("Start Listen at Port 8880")
	router.POST("/signup", postSignUp)
	router.POST("/signin", postSignIn)
	router.POST("/addnewcard", postAddNewCard)
	router.GET("/blog", getBlog)
	router.POST("/updatecard", postUpdateCard)
	router.GET("/deletecard/:cardID", postDeleteCard)
	return router
}

func postSignUp(c *gin.Context) {
	creds := model.Credentials{}
	if err := c.ShouldBindJSON(&creds); err != nil {
		c.JSON(http.StatusBadRequest, createResponse(err.Error()))
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
	creds := model.Credentials{}
	if err := c.ShouldBindJSON(&creds); err != nil {
		c.JSON(http.StatusBadRequest, createResponse(err.Error()))
		return
	}

	accessToken, err := blog.LoginAddCreateAccessToken(creds)
	if err != nil {
		c.JSON(http.StatusUnauthorized, createResponse(err.Error()))
	} else {
		c.JSON(http.StatusOK, createResponseWithAccessToken(accessToken))
	}
}

func postAddNewCard(c *gin.Context) {
	token := getTokenFromHeader(c)
	if token == "" {
		return
	}

	claims, err := blog.VerifyJWTToken(token)
	if err != nil {
		c.JSON(http.StatusUnauthorized, createResponse(err.Error()))
		return
	}

	card := model.Card{}
	if err := c.ShouldBindJSON(&card); err != nil {
		c.Status(http.StatusBadRequest)
		return
	}

	//set user id of card
	card.UserID = claims.ID
	errCreateCard := blog.CreateNewCard(card)
	if errCreateCard != nil {
		c.JSON(http.StatusOK, createResponse(err.Error()))
	} else {
		c.JSON(http.StatusOK, createResponse(""))
	}
}

func getBlog(c *gin.Context) {
	token := getTokenFromHeader(c)
	if token == "" {
		return
	}

	_, err := blog.VerifyJWTToken(token)
	if err != nil {
		c.JSON(http.StatusUnauthorized, createResponse(err.Error()))
		return
	}

	cards := blog.GetAllCard()
	c.JSON(http.StatusOK, createResponseWithCards(cards))
}

func postUpdateCard(c *gin.Context) {
	token := getTokenFromHeader(c)
	if token == "" {
		return
	}

	claims, err := blog.VerifyJWTToken(token)
	if err != nil {
		c.JSON(http.StatusUnauthorized, createResponse(err.Error()))
		return
	}

	card := model.Card{}
	if err := c.ShouldBindJSON(&card); err != nil {
		c.JSON(http.StatusBadRequest, createResponse(err.Error()))
		return
	}

	if card.ID == 0 {
		c.JSON(http.StatusBadRequest, createResponse("Invalid Json"))
		return
	}

	card.UserID = claims.ID
	errUpdateCard := blog.UpdateCard(card)
	if errUpdateCard != nil {
		c.JSON(http.StatusOK, createResponse(errUpdateCard.Error()))
	} else {
		c.JSON(http.StatusOK, createResponse(""))
	}
}

func postDeleteCard(c *gin.Context) {
	token := getTokenFromHeader(c)
	if token == "" {
		return
	}

	claims, err := blog.VerifyJWTToken(token)
	if err != nil {
		c.JSON(http.StatusUnauthorized, createResponse(err.Error()))
		return
	}

	id := c.Param("cardID")
	cardID, errCardID := strconv.ParseInt(id, 10, 32)
	if errCardID != nil {
		c.JSON(http.StatusBadRequest, createResponse("Bad Request"))
		return
	}

	errDeleteCard := blog.DeleteCard(int(cardID), claims.ID)
	if errDeleteCard != nil {
		c.JSON(http.StatusOK, createResponse(errDeleteCard.Error()))
	} else {
		c.JSON(http.StatusOK, createResponse(""))
	}

}

func getTokenFromHeader(c *gin.Context) string {
	bearerScheme := "Bearer "
	authHeader := c.GetHeader("Authorization")
	if authHeader == "" {
		c.JSON(http.StatusUnauthorized, createResponse("Not have authorized"))
		return ""
		// return "", errors.New("Not have authorized")
	}
	token := authHeader[len(bearerScheme):]
	return token
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

func createResponseWithAccessToken(accessToken string) model.ResponseAPI {
	resAPI := createResponse("")
	if accessToken != "" {
		resAPI.AccessToken = accessToken
	}
	return resAPI
}

func createResponseWithCards(cards []model.Card) model.ResponseAPI {
	resAPI := createResponse("")
	resAPI.Cards = cards

	return resAPI
}

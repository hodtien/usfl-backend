package main

import (
	"net/http"
	"time"

	"web/usfl-backend/mongohandler"

	"github.com/labstack/echo"
	"github.com/labstack/echo/middleware"
)

var mgodb mongohandler.Mgo

func main() {
	mgodb.InitialDatabase()
	defer mgodb.Close()

	e := echo.New()
	s := &http.Server{
		Addr:         "localhost:2000",
		ReadTimeout:  20 * time.Minute,
		WriteTimeout: 20 * time.Minute,
	}

	e.Use(middleware.CORS())

	// ---- BOOK ----
	e.POST("/api/book/insertABook", mongohandler.InsertABook)
	e.POST("/api/book/insertManyBooks", mongohandler.InsertManyBooks)
	e.POST("/api/book/insertBooksInCategory", mongohandler.InsertBooksInCategory)

	e.GET("/api/book/allNew", mongohandler.GetAllNewBook)
	e.GET("/api/book/allPopular", mongohandler.GetAllPopularBook)
	e.GET("/api/book/all", mongohandler.GetAllBook)
	e.GET("/api/book/detail", mongohandler.GetDetailABook)

	e.GET("/api/book/search", mongohandler.SearchBook)

	// ---- USER ----
	e.POST("/api/user/signup", mongohandler.UserSignUp)
	e.GET("/api/user/signin", mongohandler.UserSignIn)
	e.GET("/api/user/info", mongohandler.UserInfo)
	e.POST("api/user/update", mongohandler.UpdateUser)

	e.POST("/api/user/updatePassword", mongohandler.UpdatePassword)

	e.POST("/api/user/borrowBook", mongohandler.BorrowBook)
	e.POST("/api/user/updateBorrowBook", mongohandler.UpdateBorrowBook)
	e.GET("/api/user/yourBook", mongohandler.YourBook)


	e.POST("/api/comment/addComment", mongohandler.AddComment)
	e.GET("/api/comment/getComments", mongohandler.GetComments)

	// ---- CATEGORY ----
	e.POST("/api/category/createACategory", mongohandler.CreateACategory)
	e.GET("/api/category/getACategory", mongohandler.GetACategory)
	e.GET("/api/category/getAllCategory", mongohandler.GetAllCategory)



	e.Logger.Fatal(e.StartServer(s))
}

package mongohandler

import (
	"encoding/json"
	"fmt"
	"web/usfl-backend/models"

	"github.com/labstack/echo"
	"github.com/lithammer/shortuuid"
)

// GetAllBook - get all book
func GetAllBook(c echo.Context) error {
	collection := "all@book"

	allbook := Mgodb.FindAll(MongoHost, DBName, collection)

	return c.JSON(200, allbook)
}

// InsertABook - insert a book to db
func InsertABook(c echo.Context) error {
	var book models.Book

	err := c.Bind(&book)
	if err != nil {
		fmt.Println(err)
		return c.JSON(400, map[string]interface{}{"code": "-1", "message": err})
	}
	book.ID = shortuuid.New()

	var bodyBytes []byte

	bodyBytes, err = json.Marshal(book)
	if err != nil {
		fmt.Println(err)
		return c.JSON(400, map[string]interface{}{"code": "-1", "message": err})
	}

	bookMap := make(map[string]interface{})

	err = json.Unmarshal(bodyBytes, &bookMap)
	if err != nil {
		fmt.Println(err)
		return c.JSON(400, map[string]interface{}{"code": "-1", "message": err})
	}

	collection := "all@book"

	Mgodb.SaveMongo(MongoHost, DBName, collection, book.ID, bookMap)

	fmt.Println("Inserted a single book: ", book)

	return c.JSON(200, map[string]interface{}{"code": "-1", "message": "INSERT SUCCEEDED"})
}

// InsertManyBooks -InsertManyBooks
func InsertManyBooks(c echo.Context) error {
	var books = new(struct {
		ManyBooks []models.Book `json:"many_books"`
	})

	err := c.Bind(&books)
	if err != nil {
		fmt.Println(err)
		return c.JSON(400, map[string]interface{}{"code": "-1", "message": err})
	}

	for _, book := range books.ManyBooks {
		book.ID = shortuuid.New()

		var bodyBytes []byte

		bodyBytes, err = json.Marshal(book)
		if err != nil {
			fmt.Println(err)
			return c.JSON(400, map[string]interface{}{"code": "-1", "message": err})
		}

		bookMap := make(map[string]interface{})

		err = json.Unmarshal(bodyBytes, &bookMap)
		if err != nil {
			fmt.Println(err)
			return c.JSON(400, map[string]interface{}{"code": "-1", "message": err})
		}

		collection := "all@book"

		Mgodb.SaveMongo(MongoHost, DBName, collection, book.ID, bookMap)
	}

	return c.JSON(200, map[string]interface{}{"code": "0", "message": "INSERT SUCCEEDED"})
}

// GetDetailABook - Get detail of a book
func GetDetailABook(c echo.Context) error {
	bookID := c.QueryParam("bookID")

	collection := "all@book"

	book, err := Mgodb.FindByID(MongoHost, DBName, collection, bookID)
	if err != nil {
		return c.JSON(400, map[string]interface{}{"code": "-1", "message": err})
	}

	fmt.Println("Get detail a book", book)
	return c.JSON(200, map[string]interface{}{"code": "0", "message": book})
}

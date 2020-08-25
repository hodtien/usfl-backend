package mongohandler

import (
	"encoding/json"
	"fmt"
	"strconv"

	// "strings"
	// "unicode"
	"web/usfl-backend/models"

	"github.com/labstack/echo"
	// "github.com/lithammer/shortuuid"
)

// GetAllBook - get all book
func GetAllBook(c echo.Context) error {
	collection := "all@book"

	allbook := Mgodb.FindAll(MongoHost, DBName, collection)

	ret := deleteSomeFieldInAllBooks(allbook)

	return c.JSON(200, ret)
}

// InsertABook - insert a book to db
func InsertABook(c echo.Context) error {
	var book models.Book

	err := c.Bind(&book)
	if err != nil {
		fmt.Println(err)
		return c.JSON(400, map[string]interface{}{"code": "-1", "message": err})
	}

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

	Mgodb.SaveMongo(MongoHost, DBName, collection, book.Sku, bookMap)

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

		Mgodb.SaveMongo(MongoHost, DBName, collection, book.Sku, bookMap)
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

// SearchBook - SearchBook
func SearchBook(c echo.Context) error {
	search := new(struct {
		Key string `json:"key" validate:"required"`
	})

	err := c.Bind(&search)
	if err != nil {
		return c.JSON(400, map[string]interface{}{"code": "-1", "message": err})
	}

	collection := "all@book"

	// key := strings.ToLowerSpecial(unicode.TurkishCase, search.Key)
	key := search.Key

	data, err := Mgodb.SearchInMongo(MongoHost, DBName, collection, key, "title")
	if err != nil {
		return c.JSON(400, map[string]interface{}{"code": "-1", "message": err})
	}

	return c.JSON(200, map[string]interface{}{"code": "0", "message": data})
}

// deleteSomeFieldInAllBooks - deleteSomeFieldInAllBooks
func deleteSomeFieldInAllBooks(allbook []map[string]interface{}) []map[string]interface{} {
	var ret []map[string]interface{}

	for _, val := range allbook {
		delete(val, "views")
		delete(val, "author")
		delete(val, "publisher")
		delete(val, "page")
		delete(val, "type")
		delete(val, "language")
		delete(val, "tag")
		delete(val, "remain")
		delete(val, "rate")
		delete(val, "numRate")
		delete(val, "intro")

		ret = append(ret, val)
	}

	return ret
}

// GetAllNewBook - GetAllNewBook
func GetAllNewBook(c echo.Context) error {
	collection := "all@book"

	allbook := Mgodb.FindAll(MongoHost, DBName, collection)

	for i := 0; i < len(allbook)-1; i++ {
		for j := i + 1; j < len(allbook); j++ {
			ti := allbook[i]
			tj := allbook[j]

			vi, err := strconv.Atoi(ti["views"].(string))
			if err != nil {
				return c.JSON(400, map[string]interface{}{"code": "-1", "message": err})
			}

			vj, err := strconv.Atoi(tj["views"].(string))
			if err != nil {
				return c.JSON(400, map[string]interface{}{"code": "-1", "message": err})
			}

			if vi > vj {
				tmp := allbook[i]
				allbook[i] = allbook[j]
				allbook[j] = tmp
			}
		}
	}

	data := deleteSomeFieldInAllBooks(allbook)

	return c.JSON(200, map[string]interface{}{"code": "0", "data": data, "total": len(data)})
}

// GetAllPopularBook - GetAllPopularBook
func GetAllPopularBook(c echo.Context) error {
	collection := "all@book"

	allbook := Mgodb.FindAll(MongoHost, DBName, collection)

	for i := 0; i < len(allbook)-1; i++ {
		for j := i + 1; j < len(allbook); j++ {
			ti := allbook[i]
			tj := allbook[j]

			vi, err := strconv.Atoi(ti["views"].(string))
			if err != nil {
				return c.JSON(400, map[string]interface{}{"code": "-1", "message": err})
			}

			vj, err := strconv.Atoi(tj["views"].(string))
			if err != nil {
				return c.JSON(400, map[string]interface{}{"code": "-1", "message": err})
			}

			if vi < vj {
				tmp := allbook[i]
				allbook[i] = allbook[j]
				allbook[j] = tmp
			}
		}
	}

	data := deleteSomeFieldInAllBooks(allbook)

	return c.JSON(200, map[string]interface{}{"code": "0", "data": data, "total": len(data)})
}

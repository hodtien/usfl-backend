package mongohandler

import (
	"encoding/json"
	"fmt"
	"strconv"
	"strings"

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

	for key, book := range ret {
		book, err := convertImage(book)
		if err != nil {
			return c.JSON(400, map[string]interface{}{"code": "-1", "message": err})
		}
		ret[key] = book
	}

	return c.JSON(200, ret)
}

func convertImage(book map[string]interface{}) (map[string]interface{}, error) {
	bookBytes, err := json.Marshal(book)
	if err != nil {
		return nil, err
	}

	bookStruct := new(models.Book)
	err = json.Unmarshal(bookBytes, &bookStruct)
	if err != nil {
		return nil, err
	}

	for key, val := range bookStruct.Images {
		println(val)
		val = strings.ReplaceAll(val, "/src/asserts/img/book-", "")
		val = strings.ReplaceAll(val, ".jpg", "")
		bookStruct.Images[key] = val
	}

	bookBytes, err = json.Marshal(bookStruct)
	if err != nil {
		return nil, err
	}

	err = json.Unmarshal(bookBytes, &book)
	if err != nil {
		return nil, err
	}

	return book, nil
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

	book, err = convertImage(book)
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

	key := strings.ToLower(search.Key)
	key = strings.ReplaceAll(key, "  ", " ")

	data, err := Mgodb.SearchInMongo(MongoHost, DBName, collection, key, "en_title")
	if err != nil {
		return c.JSON(400, map[string]interface{}{"code": "-1", "message": err})
	}

	for key, book := range data {
		book, err = convertImage(book)
		if err != nil {
			return c.JSON(400, map[string]interface{}{"code": "-1", "message": err})
		}
		data[key] = book
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

	for key, book := range data {
		book, err := convertImage(book)
		if err != nil {
			return c.JSON(400, map[string]interface{}{"code": "-1", "message": err})
		}
		data[key] = book
	}

	return c.JSON(200, data)
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

	for key, book := range data {
		book, err := convertImage(book)
		if err != nil {
			return c.JSON(400, map[string]interface{}{"code": "-1", "message": err})
		}
		data[key] = book
	}

	return c.JSON(200, data)
}

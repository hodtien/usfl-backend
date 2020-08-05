package mongohandler

import (
	"encoding/json"
	"fmt"
	"web/usfl-backend/models"

	"github.com/labstack/echo"
)

// GetAllCategory - get all category
func GetAllCategory(c echo.Context) error {
	collection := "all@category"

	allCategory := Mgodb.FindAll(MongoHost, DBName, collection)

	return c.JSON(200, map[string]interface{}{"code": "0", "message": allCategory})
}

// CreateACategory - CreateACategory
func CreateACategory(c echo.Context) error {
	var newCategory models.Category
	err := c.Bind(&newCategory)
	if err != nil {
		fmt.Println(err)
		return c.JSON(400, map[string]interface{}{"code": "-1", "message": err})
	}

	newCategory.ID = "Category@" + newCategory.Name

	var bodyBytes []byte

	bodyBytes, err = json.Marshal(newCategory)
	if err != nil {
		fmt.Println(err)
		return c.JSON(400, map[string]interface{}{"code": "-1", "message": err})
	}

	category := make(map[string]interface{})

	err = json.Unmarshal(bodyBytes, &category)
	if err != nil {
		fmt.Println(err)
		return c.JSON(400, map[string]interface{}{"code": "-1", "message": err})
	}

	collection := "all@category"

	Mgodb.SaveMongo(MongoHost, DBName, collection, newCategory.ID, category)

	fmt.Println("Create a category: ", category)

	return c.JSON(200, map[string]interface{}{"code": "0", "message": category})
}

// InsertBooksInCategory - InsertBooksInCategory
func InsertBooksInCategory(c echo.Context) error {
	books := new(struct {
		CategoryName string   `json:"category_name"`
		BooksID      []string `json:"booksID"`
	})

	err := c.Bind(&books)
	if err != nil {
		return c.JSON(400, map[string]interface{}{"code": "-1", "message": err})
	}

	collection := "Category@" + books.CategoryName
	Category, err := Mgodb.FindByID(MongoHost, DBName, "all@category", collection)
	if err != nil {
		return c.JSON(400, map[string]interface{}{"code": "-1", "message": err})
	}
	if Category == nil {
		return c.JSON(400, map[string]interface{}{"code": "-1", "message": "Category not found!!!"})
	}

	for _, v := range books.BooksID {
		book, err := Mgodb.FindByID(MongoHost, DBName, "all@book", v)
		if err != nil {
			fmt.Println("ERROR in INSERT a book in category:", err)
			continue
		}
		Mgodb.SaveMongo(MongoHost, DBName, collection, v, book)
	}

	return c.JSON(200, map[string]interface{}{"code": "0", "message": "Insert succeeded"})
}

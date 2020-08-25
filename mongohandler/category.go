package mongohandler

import (
	"fmt"

	"github.com/labstack/echo"
)

// GetAllCategory - get all category
func GetAllCategory(c echo.Context) error {
	collection := "all@category"

	allCategory := Mgodb.FindAll(MongoHost, DBName, collection)

	return c.JSON(200, map[string]interface{}{"code": "0", "message": allCategory})
}

// GetACategory - GetACategory
func GetACategory(c echo.Context) error {
	category := c.QueryParam("category")

	collection := "Category@" + category

	bookInCategory := Mgodb.FindAll(MongoHost, DBName, collection)

	return c.JSON(200, map[string]interface{}{"code": "0", "message": bookInCategory})
}

// CreateACategory - CreateACategory
func CreateACategory(c echo.Context) error {
	categoryName := c.QueryParam("name")

	CategoryID := "Category@" + categoryName

	category := make(map[string]interface{})

	category["name"] = categoryName

	collection := "all@category"

	Mgodb.SaveMongo(MongoHost, DBName, collection, CategoryID, category)

	fmt.Println("Create a category successfully: ", categoryName)

	return c.JSON(200, map[string]interface{}{"code": "0", "message": "Create a category successfully!"})
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

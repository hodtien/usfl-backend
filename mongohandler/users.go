package mongohandler

import (
	// "crypto/sha256"
	"encoding/json"
	"fmt"
	"strconv"
	"strings"
	"time"
	"web/usfl-backend/models"

	"github.com/google/uuid"
	"github.com/labstack/echo"
	"github.com/rs/xid"
)

// UserSignUp - create user
func UserSignUp(c echo.Context) error {
	var user models.User
	err := c.Bind(&user)
	if err != nil {
		fmt.Println(err)
		return c.JSON(400, map[string]interface{}{"code": "-1", "message": err})
	}
	collection := "users"

	if _, ok := Mgodb.FindOneMongoByField(MongoHost, DBName, collection, "username", user.Username); ok == true {
		return c.JSON(400, "Username is exists")
	}

	// user.Password = fmt.Sprintf("%v", sha256.Sum256([]byte(user.Password)))

	var bodyBytes []byte

	bodyBytes, err = json.Marshal(user)
	if err != nil {
		fmt.Println(err)
		return c.JSON(400, map[string]interface{}{"code": "-1", "message": err})

	}

	userMap := make(map[string]interface{})

	err = json.Unmarshal(bodyBytes, &userMap)
	if err != nil {
		fmt.Println(err)
		return c.JSON(400, map[string]interface{}{"code": "-1", "message": err})
	}

	uid := uuid.New().String()

	Mgodb.SaveMongo(MongoHost, DBName, collection, uid, userMap)

	return c.JSON(200, uid)
}

// UserSignIn - user log in
func UserSignIn(c echo.Context) error {
	username := c.QueryParam("username")
	password := c.QueryParam("password")

	// password = fmt.Sprintf("%v", sha256.Sum256([]byte(password)))

	collection := "users"

	ret, ok := Mgodb.FindOneMongoByField(MongoHost, DBName, collection, "username", username)
	if !ok {
		return c.JSON(400, map[string]interface{}{"code": "-1", "message": false})
	}

	if ret == nil {
		return c.JSON(400, map[string]interface{}{"code": "-1", "message": "Incorrect Username or Password1"})
	}

	if password != ret["password"] {
		return c.JSON(400, map[string]interface{}{"code": "-1", "message": "Incorrect Username or Password2"})
	}



	return c.JSON(200, ret["_id"])
}

// UserInfo - user info
func UserInfo(c echo.Context) error {
	userID := c.QueryParam("userid")

	collection := "users"

	ret, err := Mgodb.FindByID(MongoHost, DBName, collection, userID)
	if err != nil {
		return c.JSON(400, map[string]interface{}{"code": "-1", "message": err})
	}

	if ret == nil {
		return c.JSON(400, map[string]interface{}{"code": "-1", "message": "false"})
	}

	return c.JSON(200, map[string]interface{}{"code": "0", "message": "true", "data": ret})
}

// UpdateUser - UpdateUser
func UpdateUser(c echo.Context) error {
	userID := c.QueryParam("userid")

	var user models.User
	err := c.Bind(&user)
	if err != nil {
		fmt.Println(err)
		return c.JSON(400, map[string]interface{}{"code": "-1", "message": err})
	}

	collection := "users"

	ret, err := Mgodb.FindByID(MongoHost, DBName, collection, userID)
	if err != nil {
		return c.JSON(400, map[string]interface{}{"code": "-1", "message": err})
	}

	if ret == nil {
		return c.JSON(400, map[string]interface{}{"code": "-1", "message": "false"})
	}

	var bodyBytes []byte

	bodyBytes, err = json.Marshal(user)
	if err != nil {
		fmt.Println(err)
		return c.JSON(400, map[string]interface{}{"code": "-1", "message": err})

	}

	userMap := make(map[string]interface{})

	err = json.Unmarshal(bodyBytes, &userMap)
	if err != nil {
		fmt.Println(err)
		return c.JSON(400, map[string]interface{}{"code": "-1", "message": err})
	}

	Mgodb.SaveMongo(MongoHost, DBName, collection, userID, userMap)

	return c.JSON(200, map[string]interface{}{"code": "0", "message": "true", "data": userMap})
}

// UpdatePassword -v
func UpdatePassword(c echo.Context) error {
	user := new(struct {
		Username    string `json:"username" validate:"required"`
		OldPassword string `json:"old_password" validate:"required"`
		NewPassword string `json:"new_password" validate:"required"`
	})

	err := c.Bind(&user)
	if err != nil {
		fmt.Println(err)
		return c.JSON(400, map[string]interface{}{"code": "-1", "message": err})
	}

	username := user.Username
	oldPassword := user.OldPassword
	newPassword := user.NewPassword

	// oldPassword = fmt.Sprintf("%v", sha256.Sum256([]byte(oldPassword)))
	// newPassword = fmt.Sprintf("%v", sha256.Sum256([]byte(newPassword)))

	collection := "users"

	ret, err := Mgodb.FindByID(MongoHost, DBName, collection, username)
	if err != nil {
		return c.JSON(400, map[string]interface{}{"code": "-1", "message": err})

	}

	if ret == nil {
		return c.JSON(400, map[string]interface{}{"code": "-1", "message": "Incorrect Username"})
	}

	if oldPassword != ret["password"] {
		return c.JSON(400, map[string]interface{}{"code": "-1", "message": "Incorrect Password"})
	}

	// user_map := map[string]interface{}{ "password": newPassword}
	// data, err := json.Marshal(user_map)
	// if err != nil {
	// 	fmt.Println(err)
	// 	return c.JSON(400, map[string]interface{}{"code": "-1", "message": err})
	// }

	err = Mgodb.UpdateMongo(MongoHost, DBName, collection, username, "password", newPassword)
	if err != nil {
		fmt.Println(err)
		return c.JSON(400, map[string]interface{}{"code": "-1", "message": err})
	}

	return c.JSON(200, map[string]interface{}{"code": "0", "message": "Update Password Successfully"})
}

// BorrowBook - BorrowBook
func BorrowBook(c echo.Context) error {
	borrowBook := new(models.BorrowBook)
	err := c.Bind(&borrowBook)
	if err != nil {
		return c.JSON(400, map[string]interface{}{"code": "-1", "message": err})
	}
	borrowBook.BorrowID = xid.New().String()

	user, err := Mgodb.FindByID(MongoHost, DBName, "users", borrowBook.UserID)
	if err != nil {
		return c.JSON(400, map[string]interface{}{"code": "-1", "message": err})
	}
	if user == nil {
		return c.JSON(400, map[string]interface{}{"code": "-1", "message": "User not exist"})
	}

	book, err := Mgodb.FindByID(MongoHost, DBName, "all@book", borrowBook.BookID)
	if err != nil {
		return c.JSON(400, map[string]interface{}{"code": "-1", "message": err})
	}

	bookCount, err := strconv.Atoi(fmt.Sprintf("%v", book["remain"]))
	if err != nil {
		return c.JSON(400, map[string]interface{}{"code": "-1", "message": err})
	}
	if bookCount <= 0 {
		return c.JSON(400, map[string]interface{}{"code": "-1", "message": "OUT OF STOCK!"})
	}

	borrowBook.Time = time.Now().Format("2006-01-02 15:04")
	borrowBook.Status = "Pending"

	dataBytes, err := json.Marshal(borrowBook)
	if err != nil {
		return c.JSON(400, map[string]interface{}{"code": "-1", "message": err})
	}
	borrowBookMap := make(map[string]interface{})
	err = json.Unmarshal(dataBytes, &borrowBookMap)
	if err != nil {
		return c.JSON(400, map[string]interface{}{"code": "-1", "message": err})
	}

	Mgodb.SaveMongo(MongoHost, DBName, borrowBook.UserID + "@Borrow", borrowBook.BorrowID, borrowBookMap)

	return c.JSON(200, map[string]interface{}{"code": "0", "message": "Borrow Book Status: " + borrowBook.Status})
}

// YourBook - YourBook
func YourBook(c echo.Context) error {
	userID := c.QueryParam("userid")

	user, err := Mgodb.FindByID(MongoHost, DBName, "users", userID)
	if err != nil {
		return c.JSON(400, map[string]interface{}{"code": "-1", "message": err})
	}

	username := user["username"].(string)

	data := Mgodb.FindAll(MongoHost, DBName, username + "@Borrow")


	return c.JSON(200, data)
}

// UpdateBorrowBook - UpdateBorrowBook
func UpdateBorrowBook(c echo.Context) error {
	userID := c.QueryParam("userid")
	borrowID := c.QueryParam("borrowID")
	status := c.QueryParam("status")

	borrowData, err := Mgodb.FindByID(MongoHost, DBName, userID + "@Borrow", borrowID)
	if err != nil {
		return c.JSON(400, map[string]interface{}{"code": "-1", "message": err})
	}

	borrowData["status"] = status
	stt := strings.ToLower(status)

	bookID := borrowData["bookID"].(string)

	book, err := Mgodb.FindByID(MongoHost, DBName, "all@book", bookID)
	if err != nil {
		return c.JSON(400, map[string]interface{}{"code": "-1", "message": err})
	}

	bookCount, err := strconv.Atoi(fmt.Sprintf("%v", book["remain"]))
	if err != nil {
		return c.JSON(400, map[string]interface{}{"code": "-1", "message": err})
	}


	if stt == "borrowing" {
		Mgodb.UpdateMongo(MongoHost, DBName, "all@book", borrowData["bookID"].(string), "remain", strconv.Itoa(bookCount - 1))
		if bookCount <= 0 {
			return c.JSON(400, map[string]interface{}{"code": "-1", "message": "OUT OF STOCK!"})
		}
	}
	if stt == "returned" {
		Mgodb.UpdateMongo(MongoHost, DBName, "all@book", borrowData["bookID"].(string), "remain", strconv.Itoa(bookCount + 1))
	} else {
		return c.JSON(400, map[string]interface{}{"code": "-1", "message": "Status Invalid"})
	}

	Mgodb.SaveMongo(MongoHost, DBName, userID + "@Borrow", borrowID, borrowData)

	return c.JSON(200, map[string]interface{}{"code": "0", "message": "Borrow Book Status Updated: " + status})
}

// AddComment - AddComment
func AddComment(c echo.Context) error {
	cmt := new(models.Comment)

	err := c.Bind(&cmt)
	if err != nil {
		return c.JSON(400, map[string]interface{}{"code": "-1", "message": err})
	}

	cmt.Time = time.Now().Format("2006-01-02 15:04")
	timestamp := time.Now().Local().Unix()
	cmt.Timestamp = strconv.FormatInt(timestamp, 10)

	var bodyBytes []byte

	bodyBytes, err = json.Marshal(cmt)
	if err != nil {
		fmt.Println(err)
		return c.JSON(400, map[string]interface{}{"code": "-1", "message": err})
	}

	commentMap := make(map[string]interface{})

	err = json.Unmarshal(bodyBytes, &commentMap)
	if err != nil {
		fmt.Println(err)
		return c.JSON(400, map[string]interface{}{"code": "-1", "message": err})
	}

	collection := "comment@" + cmt.BookID

	Mgodb.SaveMongo(MongoHost, DBName, collection, cmt.Timestamp, commentMap)

	return c.JSON(200, map[string]interface{}{"code": "0", "message": commentMap})

}

// GetComments - GetComments
func GetComments(c echo.Context) error {
	bookID := c.QueryParam("bookID")

	collection := "comment@" + bookID
	data, hasData := Mgodb.GetDataInOrderedByField(MongoHost, DBName, collection, "timestamp")

	if !hasData {
		return c.JSON(400, map[string]interface{}{"code": "-1", "message": "Data Not Found!"})
	}

	return c.JSON(200, map[string]interface{}{"code": "0", "message": data})
}

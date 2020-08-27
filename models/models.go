package models

// Book struct - Create Book struct
type Book struct {
	Title     string   `json:"title"`     // Tên sách
	Sku       string   `json:"sku"`       // Mã SKU
	Views     string   `json:"views"`     // Lượt xem
	Author    string   `json:"author"`    // Tên tác giả
	Publisher string   `json:"publisher"` // Nhà xuất bản
	Page      int      `json:"page"`      // Số trang, kiểu int
	Type      string   `json:"type"`      // Bìa cứng/bìa mềm
	Language  string   `json:"language"`  // Ngôn ngữ
	Tag       []string `json:"tag"`       // Tag quyển sách: Lập trình, tiếng anh, hay lịch sử...
	Remain    int      `json:"remain"`    // Số lượng sách còn trong kho, kiểu int
	Rate      string   `json:"rate"`      // Đánh giá bao nhiêu sao (nhỏ hơn hoặc bằng 5)
	NumRate   int      `json:"numRate"`   // Số lượt đánh giá
	Intro     string   `json:"intro"`     // Lời giới thiệu
	Images    []string `json:"images"`    // Các đường dẫn tới file ảnh
}

// User - user
type User struct {
	Username string `json:"username" validate:"required"`
	Password string `json:"password" validate:"required"`

	Name   string `json:"name"`
	Birth  string `json:"birth"`
	Gender string `json:"gender"`
	Photo  string `json:"photo"`

	Card  string `json:"card"`
	Email string `json:"email"`
	Phone string `json:"phone"`
}

// Comment - comment
type Comment struct {
	BookID    string   `json:"bookID" validate:"required"`
	Username  string   `json:"username" validate:"required"`
	Title     string   `json:"title"`
	Content   string   `json:"content" validate:"required"`
	Rate      int      `json:"rate"`
	Time      string   `json:"time"`
	Timestamp string   `json:"timestamp"`
	Images    []string `json:"images"`
}

// BorrowBook - BorrowBook
type BorrowBook struct {
	UserID string `json:"userid"`
	BorrowID string `json:"borrowID"`
	BookID   string `json:"bookID"`
	Status   string `json:"status"`
	Time      string   `json:"time"`
}

package model

import (
	"fmt"
	"time"

	"gorm.io/driver/mysql"
	"gorm.io/gorm"
)

var db *gorm.DB
var DB = db

type User struct {
	
	ID				uint64			`gorm:"id" gorm:"primaryKey"`
	FirstName		string			`gorm:"type:text;not null"`
	LastName		string			`gorm:"type:text;not null"`
	Email			string			`gorm:"type:text;not null"`
	Password		string			`gorm:"type:text;not null"`
	Photo			string			`gorm:"type:text;null"`
	// Search			[]Search
	CreatedAt		time.Time
	UpdatedAt		time.Time
	DeletedAt    	gorm.DeletedAt `gorm:"index"`
} 

type Search struct {
	ID				uint 			`json:"id" gorm:"primaryKey"`
	Term			string			`json:"term"`
	UserRefer		uint64			`json:"user_id"`
	// UserId			uint64			`json:"user_id"`
	User			User			`gorm:"foreignKey:UserRefer"`
	CreatedAt		time.Time
	DeletedAt		time.Time
}

type SearchResponse struct {
	ID				uint64 			`json:"id"`
	Term			string			`json:"term"`
	User			User
	DeletedAt		time.Time
}

type UserResponse struct {
	ID        uint64 `json:"id,omitempty"`
	Name      string    `json:"name,omitempty"`
	Email     string    `json:"email,omitempty"`
	Photo     string    `json:"photo,omitempty"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}

func FilterUserRecord(user *User) UserResponse {
	return UserResponse{
		ID:        user.ID,
		Name:      user.FirstName,
		Email:     user.Email,
		Photo:     user.Photo,
		CreatedAt: user.CreatedAt,
		UpdatedAt: user.UpdatedAt,
	}
}


func Setup() {

	dsn := "root:@tcp(127.0.0.1:3306)/authx?charset=utf8mb4&parseTime=True&loc=Local"
	var err error
	db, err = gorm.Open(mysql.Open(dsn), &gorm.Config{})
	if err != nil {
		panic(err) 
	}

	err = db.AutoMigrate(&User{}, &Search{})
	if err != nil {
		fmt.Println(err)
	}
}
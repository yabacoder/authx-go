package model

import "fmt"

func SearchTerm(term string) ([]User, error) {
	var users []User

	fmt.Println(term)
	tnx := db.Where("first_name LIKE ? OR last_name LIKE ?",
		term, term).Find(&users)

	if tnx.Error != nil {
		return nil, tnx.Error
	}

	return users, nil
}

func SaveUserSearch(term string, user any) error {
	var search Search
	// search.UserRefer = user.ID
	search.Term = term
	search.User = user.(User)

	tnx := db.Create(&search)

	return tnx.Error
}

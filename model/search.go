package model

func SearchTerm(term string) (User, error) {
	var users User

	tnx := db.Where("first_name LIKE ? AND last_name LIKE ?", 
	term, term).Find(&users)

	if tnx.Error != nil {
		return users, tnx.Error
	}

	return users , nil
}

func SaveUserSearch(term, userId string) error {
	var search []Search

	tnx := db.Create(&search)

	return tnx.Error
}

package model

func GetAllUsers() ([]User, error) {

	var users []User

	tnx := db.Find(&users)

	return users, tnx.Error
}

func GetUser(id any) (User, error) {
	var user User
	tnx := db.First(&user, "id = ?", id)

	return user, tnx.Error
}

func CreateUser(user *User) error {
	tnx := db.Create(&user)

	return tnx.Error
}

func UpdateUser(user *User) error {

	tnx := db.Model(user).Select("FirstName", "LastName", "Email").Updates(user)

	return tnx.Error
}

func FindUserWithPassword(user User) (User, error) {
	tnx := db.Where("email = ? ",
		user.Email).First(&user)

	return user, tnx.Error
}

func UploadUserPhoto(user User, url string) error {
	// id := user.ID
	tnx := db.Where("id = ? ",
		user.ID).First(&user)

	if tnx.Error != nil {
		return tnx.Error
	}
	// fmt.Println(url)
	tnx = db.Model(user).Update("photo", url)
	// fmt.Println(tnx)

	return tnx.Error
}

// func FindUserWithPassword(user *User) (error) {
// 	tnx := db.Where("email = ? AND password = ?",
// 	user.Email, user.Password).First(&user)

// 	return tnx.Error
// }

// func DeleteUser(id uint64) error {
// 	tnx := db.Unscoped().Delete(&User, id)
// 	return tnx.Error
// }

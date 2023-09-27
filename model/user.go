package model


func GetAllUsers() ([]User, error) {
	var users []User

	return users, nil
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

func UpdateUser(user User) error {
	tnx := db.Save(&user)

	return tnx.Error
}

func FindUserWithPassword(user User) (User, error) {
	tnx := db.Where("email = ? ", 
	user.Email).First(&user)

	return user, tnx.Error
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
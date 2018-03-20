package user

import (
	"fmt"
	"log"

	"github.com/jinzhu/gorm"
)

func CreateUser(db *gorm.DB) *User {

	fmt.Print("What is the email?: ")
	var email string
	fmt.Scanln(&email)

	fmt.Print("What is the password? Leave blank for autogen: ")
	var password string
	fmt.Scanln(&password)
	if password == "" {
		password = GeneratePassword()
		fmt.Println(password)
	}
	hashed_pw := hashAndSalt([]byte(password))

	user := User{
		Email:        email,
		PasswordHash: hashed_pw,
	}
	_, err := Create(db, &user)
	if err != nil {
		log.Panic(err)
	}
	return &user

}

func LogIn(db *gorm.DB) (bool, *User) {
	fmt.Print("What is your user email?: ")
	var email string
	fmt.Scanln(&email)

	fmt.Print("What is your password?: ")
	var password string
	fmt.Scanln(&password)

	var user User
	rec_not_found := db.Where("email = ?", email).First(&user).RecordNotFound()
	if rec_not_found == true {
		fmt.Println("Couldnt find your user! try again")
		return false, &User{}
	}
	login_success := comparePasswords(user.PasswordHash, []byte(password))
	return login_success, &user

}

func GeneratePassword() string {
	return createPassword()
}

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
	_, err := createUser(db, &user)
	if err != nil {
		log.Panic(err)
	}
	return &user

}

func CreateDomain(db *gorm.DB) *Domain {
	logged_in, user, pw := LogIn(db)
	if logged_in == false {
		fmt.Println("Wasnt able to log in!")
		return &Domain{}
	}

	fmt.Print("What is the domain name?: ")
	var domain_name string
	fmt.Scanln(&domain_name)

	domain_pw := createPassword()
	fmt.Println("Your password: %s", domain_pw)
	key := user.Email + pw
	enctyped_domain_pw := encrypt([]byte(key[:32]), domain_pw)

	new_domain := Domain{
		FQDN:         domain_name,
		PasswordHash: enctyped_domain_pw,
		UserID:       user.ID,
	}
	err := db.Create(&new_domain).Error
	if err != nil {
		fmt.Println("there was an error creating db!")
		log.Panic(err)
		return &Domain{}

	}
	return &new_domain

}

func ListDomains(db *gorm.DB) []*Domain {
	logged_in, user, pw := LogIn(db)
	if logged_in == false {
		fmt.Println("Wasnt able to log in!")
		return []*Domain{
			&Domain{},
		}
	}
	var domains []*Domain
	db.Model(&user).Related(&domains)
	key := user.Email + pw
	var decrypted_pw string
	for _, domain := range domains {
		fmt.Println(domain.FQDN)
		decrypted_pw = decrypt([]byte(key[:32]), domain.PasswordHash)
		fmt.Println(decrypted_pw)

	}
	return domains

}

func LogIn(db *gorm.DB) (bool, *User, string) {
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
		return false, &User{}, ""
	}
	login_success := comparePasswords(user.PasswordHash, []byte(password))
	return login_success, &user, password

}

func GeneratePassword() string {
	return createPassword()
}

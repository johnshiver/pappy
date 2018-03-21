package user

import (
	"fmt"
	"log"
	"os"

	"github.com/jinzhu/gorm"
	"github.com/olekukonko/tablewriter"
)

func CreateUser(db *gorm.DB) {

	fmt.Print("What is the email?: ")
	var email string
	fmt.Scanln(&email)

	fmt.Print("What is the password? Leave blank for autogen: ")
	var password string
	fmt.Scanln(&password)
	if password == "" {
		password = createPassword()
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

	fmt.Printf("User %s created succesfully!", user.Email)
}

func CreateDomain(db *gorm.DB) {
	logged_in, user, pw := LogIn(db)
	if logged_in == false {
		fmt.Println("Wasnt able to log in!")
		return
	}

	fmt.Print("What is the domain name?: ")
	var domain_name string
	fmt.Scanln(&domain_name)

	fmt.Print("Password? Leave blank for auto-gen: ")
	var domain_pw string
	fmt.Scanln(&domain_pw)
	if domain_pw == "" {
		domain_pw = createPassword()
		fmt.Printf("Your password: %s\n", domain_pw)

	}

	createDomain(db, domain_name, domain_pw, pw, user)
}

func ListDomains(db *gorm.DB) {
	logged_in, user, pw := LogIn(db)
	if logged_in == false {
		fmt.Println("Wasnt able to log in!")
		return
	}
	var domains []*Domain
	db.Model(&user).Related(&domains)

	var decrypted_pw string
	var data [][]string
	key := generateEncryptionKey(user.Email, pw)
	for _, domain := range domains {
		decrypted_pw = decrypt(key, domain.PasswordHash)
		data = append(data, []string{domain.FQDN, decrypted_pw})
	}
	printDomains(data)

}

func LookupDomain(db *gorm.DB) {
	logged_in, user, pw := LogIn(db)
	if logged_in == false {
		fmt.Println("Wasnt able to log in!")
		return
	}

	fmt.Print("What is the domain name?: ")
	var domain_name string
	fmt.Scanln(&domain_name)

	query := "%" + domain_name + "%"

	var domains []*Domain
	db.Where("user_id = ? AND FQDN LIKE ?", user.ID, query).Find(&domains)

	if len(domains) == 0 {
		fmt.Printf("Couldnt find %s!\n", domain_name)
		return
	}

	var decrypted_pw string
	var data [][]string
	key := generateEncryptionKey(user.Email, pw)
	for _, domain := range domains {
		decrypted_pw = decrypt(key, domain.PasswordHash)
		data = append(data, []string{domain.FQDN, decrypted_pw})
	}
	printDomains(data)

}

func printDomains(data [][]string) {
	table := tablewriter.NewWriter(os.Stdout)
	table.SetHeader([]string{"Domain", "Password"})

	for _, v := range data {
		table.Append(v)
	}
	table.Render()
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

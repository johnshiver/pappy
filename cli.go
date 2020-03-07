package main

import (
	"fmt"
	"log"
	"os"
	"strconv"

	"github.com/olekukonko/tablewriter"
	"golang.org/x/crypto/ssh/terminal"
)

func (env *runEnv) CreateUser() {

	fmt.Print("What is the email?: ")
	var email string
	_, _ = fmt.Scanln(&email)

	fmt.Print("What is the password? Leave blank for autogen: ")
	var password string
	_, _ = fmt.Scanln(&password)
	if password == "" {
		password = createPassword()
		fmt.Println(password)
	}
	hashedPw := hashAndSalt([]byte(password))

	user := User{
		Email:        email,
		PasswordHash: hashedPw,
	}
	_, err := createUser(db, &user)
	if err != nil {
		log.Panic(err)
	}

	fmt.Printf("User %s created succesfully!", user.Email)
}

func (env *runEnv) CreatePassword() {
	loggedIn, user, pw := LogIn(db)
	if loggedIn == false {
		fmt.Println("Wasnt able to log in!")
		return
	}

	fmt.Print("What website / service is this password for?: ")
	var domainName string
	fmt.Scanln(&domainName)

	fmt.Print("Password? Leave blank for auto-gen: ")
	var domain_pw string
	fmt.Scanln(&domain_pw)
	if domain_pw == "" {
		fmt.Print("Password character limit? Leave blank for no limit")
		var char_limit string
		fmt.Scanln(&char_limit)
		if len(char_limit) == 0 {
			char_limit = "0"
		}
		c_limit, err := strconv.Atoi(char_limit)
		if err != nil {
			fmt.Println("There was an error reading your character limit! Did you submit an integer?")
			log.Panic(err)
		}
		domain_pw = createPassword()
		if c_limit > 0 {
			if c_limit > len(domain_pw) {
				c_limit = len(domain_pw)
			}
			domain_pw = domain_pw[:c_limit]
		}
		fmt.Printf("Your password: %s\n", domain_pw)

	}

	createDomain(db, domainName, domain_pw, pw, user)
}

func (env *runEnv) ListDomains() {
	loggedIn, user, pw := LogIn(db)
	if loggedIn == false {
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

func printDomains(data [][]string) {
	table := tablewriter.NewWriter(os.Stdout)
	table.SetHeader([]string{"Domain", "Password"})

	for _, v := range data {
		table.Append(v)
	}
	table.Render()
}

func (env *runEnv) LogIn() {
	fmt.Print("What is your username: ")
	var username string
	fmt.Scanln(&username)

	fmt.Println("Enter master password: ")
	password, err := terminal.ReadPassword(0)
	if err != nil {
		log.Panic(err)
	}

	user, err := env.FindByUserName(username)
	if user == nil {
		panic(fmt.Errorf("couldnt find user %s", username))
	}
	loginSuccess := comparePasswords([]byte(user.PasswordHash), password)
	if !loginSuccess {
		panic(fmt.Errorf("failed to login with credentials"))
	}
	key := generateEncryptionKey(user.UserName, string(password))

	env.mtx.Lock()
	defer env.mtx.Unlock()

	env.encryptionKey = key
	env.user = user
}

func GeneratePassword() string {
	return createPassword()
}

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

	var (
		username string
		password string
	)

	fmt.Println("Creating a new user!")
	fmt.Print("What is your username?: ")

	userNameRaw := GetUserTextInput()
	if len(userNameRaw) > 1 || len(userNameRaw) == 0 {
		log.Panic(fmt.Errorf("userName should be a single input, you gave %v", userNameRaw))
	}
	username = userNameRaw[0]

	fmt.Print("What is your master password? Leave blank to autogen: ")
	passwordRaw := GetUserTextInput()
	if len(userNameRaw) > 1 {
		log.Panic(fmt.Errorf("password should be a single input or blank, you gave %v", userNameRaw))
	}
	if len(passwordRaw) == 0 {
		password = generatePassword(-1)
		fmt.Printf("your new master password: %s\n", password)
	} else {
		password = passwordRaw[0]
	}

	hashedPw := hashAndSalt([]byte(password))
	user := User{
		UserName:     username,
		PasswordHash: hashedPw,
	}
	env.PersistUser(user)

	fmt.Printf("%s created succesfully!", user.UserName)
}

func (env *runEnv) CreatePassword() {
	var (
		domainName string
		domainPW   string
		charLimit  string
	)

	fmt.Printf("Creating a new password for user %s\n", env.user.UserName)
	fmt.Print("What website / service is this password for?: ")
	_, err := fmt.Scanln(&domainName)
	if err != nil {
		log.Panic(err)
	}

	fmt.Print("What is the password? Leave blank for auto-gen: ")
	_, err = fmt.Scanln(&domainPW)
	if err != nil {
		log.Panic(err)
	}
	if domainPW == "" {
		fmt.Print("Alright, I'll auto-gen a pw for you")
		fmt.Print("Does your password have a character limit? Leave blank for no limit")
		fmt.Scanln(&charLimit)
		if len(charLimit) == 0 {
			charLimit = "0"
		}
		cLimit, err := strconv.Atoi(charLimit)
		if err != nil {
			fmt.Println("There was an error reading your character limit! Did you submit an integer?")
			log.Panic(err)
		}
		domainPW = generatePassword(cLimit)
		fmt.Printf("Your new password for %s: %s\n", domainName, domainPW)

	}
	env.createPassword()
}

func (env *runEnv) ListPasswords() {
	var (
		decryptedPw string
		data        [][]string
	)
	passwords := env.GetPasswords()
	for _, pw := range passwords {
		decryptedPw = decrypt(env.encryptionKey, pw.PasswordHash)
		data = append(data, []string{pw.Location, decryptedPw})
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
	fmt.Println("Lets get you logged in")
	fmt.Print("What is your username: ")
	var username string
	_, err := fmt.Scanln(&username)
	if err != nil {
		log.Panic(err)
	}

	fmt.Println("Enter master password: ")
	password, err := terminal.ReadPassword(0)
	if err != nil {
		log.Panic(err)
	}

	user, err := env.FindByUserName(username)
	if err != nil {
		log.Panic(err)
	}
	if user == nil {
		log.Panic(fmt.Errorf("couldnt find user %s", username))
	}
	loginSuccess := comparePasswords([]byte(user.PasswordHash), password)
	if !loginSuccess {
		log.Panic(fmt.Errorf("failed to login with credentials"))
	}

	key := generateEncryptionKey(user.UserName, string(password))

	env.mtx.Lock()
	defer env.mtx.Unlock()

	env.encryptionKey = key
	env.user = user
}

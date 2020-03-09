package main

import (
	"fmt"
	"os"
	"strconv"

	"github.com/olekukonko/tablewriter"
	log "github.com/sirupsen/logrus"
	"golang.org/x/crypto/ssh/terminal"
)

func (env *runEnv) CreateUser() {

	var (
		username string
		password string
	)

	fmt.Println("Creating a new user!")
	fmt.Print("What is your username?: ")
	username = env.GetUserTextInput()

	fmt.Println("Time to enter your master password.")
	fmt.Println("It is VERY important to remember this password; your account cannot be recovered if you forget it.")
	fmt.Print("What is your master password? Leave blank to autogen: ")
	password = env.GetUserTextInput()
	fmt.Printf("your password %s\n", password)
	if password == "" {
		password = generatePassword(-1)
		fmt.Printf("your master password is: %v\n", password)
	}

	hashedPw := hashAndSalt([]byte(password))
	user := User{
		UserName:     username,
		PasswordHash: hashedPw,
	}
	env.persistUser(user)
	fmt.Printf("%s created succesfully!", user.UserName)
}

func (env *runEnv) CreatePassword() {
	var (
		pwLocation string
		password   string
		charLimit  string
	)

	fmt.Printf("Creating a new password for user %s\n", env.user.UserName)
	fmt.Print("What website / service is this password for?: ")
	pwLocation = env.GetUserTextInput()

	fmt.Print("What is the password? Leave blank for auto-gen: ")
	password = env.GetUserTextInput()
	if password == "" {
		fmt.Print("Alright, I'll auto-gen a pw for you")
		fmt.Print("Does your password have a character limit? Leave blank for no limit")
		charLimit = env.GetUserTextInput()
		cLimit, err := strconv.Atoi(charLimit)
		if err != nil {
			log.Fatal(fmt.Errorf("converting character limit to integer: %v", err))
		}
		password = generatePassword(cLimit)
		fmt.Printf("Your new password for %s: %s\n", pwLocation, password)
	}

	hashedPw := encrypt(env.encryptionKey, password)
	newPassword := Password{
		Location:     pwLocation,
		PasswordHash: hashedPw,
		UserID:       env.user.ID,
	}
	env.createPassword(&newPassword)
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
	var (
		username string
		password []byte
	)
	fmt.Println("Lets get you logged in")
	fmt.Print("What is your username: ")
	username = env.GetUserTextInput()

	fmt.Println("Enter master password: ")
	password, err := terminal.ReadPassword(0)
	if err != nil {
		log.Fatal(err)
	}

	user, err := env.findByUserName(username)
	if err != nil {
		log.Fatal(err)
	}
	if user == nil {
		log.Fatal(fmt.Errorf("couldnt find user %s", username))
	}
	loginSuccess := comparePasswords([]byte(user.PasswordHash), password)
	if !loginSuccess {
		log.Fatal(fmt.Errorf("failed to login with credentials"))
	}
	key := generateEncryptionKey(user.UserName, string(password))

	env.mtx.Lock()
	defer env.mtx.Unlock()

	env.encryptionKey = key
	env.user = user

	fmt.Println("logged in!")
}

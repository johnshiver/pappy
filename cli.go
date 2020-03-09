package main

import (
	"fmt"
	"strconv"

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
	if password == "" {
		password = generatePassword(-1)
	}
	fmt.Printf("your password %s\n", password)

	hashedPw := hashAndSalt([]byte(password))
	user := User{
		UserName:     username,
		PasswordHash: hashedPw,
	}
	env.userSvc.CreateUser(user)
	fmt.Printf("%s created succesfully!", user.UserName)
}

func (env *runEnv) ListUsers() {
	var userData [][]string
	users := env.userSvc.GetUsers()
	for _, u := range users {
		userData = append(userData, []string{u.UserName})
	}
	printDataTable([]string{"username"}, userData)
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
		fmt.Println("Alright, I'll auto-gen a pw for you")
		fmt.Print("Does your password have a character limit? Leave blank for no limit: ")
		charLimit = env.GetUserTextInput()
		if charLimit == "" {
			charLimit = "-1"
		}
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
	env.pwdSvc.CreatePassword(&newPassword)
}

func (env *runEnv) ListPasswords() {
	var pwdData [][]string
	passwords := env.pwdSvc.GetPasswords(env.user.ID)
	for _, pw := range passwords {
		decryptedPw := decrypt(env.encryptionKey, pw.PasswordHash)
		pwdData = append(pwdData, []string{pw.Location, decryptedPw})
	}
	printDataTable([]string{"Location", "Password"}, pwdData)
}

func (env *runEnv) DeletePassword() {
	fmt.Printf("Deleting a password for user %v\n", env.user.UserName)
	fmt.Print("Which location are you deleting?: ")
	pwdLoc := env.GetUserTextInput()
	env.pwdSvc.DeletePassword(env.user.ID, pwdLoc)
	fmt.Println("password deleted!")
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

	user, err := env.userSvc.FindByUsername(username)
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

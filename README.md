# pappy
A password manager written in go

## What is pappy?
I wrote pappy as an easy way to generate / store passwords. 

Some security precautions have been taken but I could definitely use some advice on how to make
this program better! 

The master user password and the  individual domain passwords are encrypted on disk.
The master password is a one way encryption using bcrypt, and the domain passwords take
the plain text master password (after it has been verified) combined with email to create key to encrypt the
domain passwords.

## Installation
go get
go install

Pappy stores its information in a postgres database.  These env vars must be set:

```
	db_user = os.Getenv("PW_MAN_DB_USER")
	db_password = os.Getenv("PW_MAN_DB_PW")
	db_name = os.Getenv("PW_MAN_DB_NAME")
```

## Usage
password_manager --add user
password_manager --add domain
password_manager --list domain
password_manager --lookup domain
password_manager --generate password

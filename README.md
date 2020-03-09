# pappy
A password manager written in go

## What is pappy?
I wrote pappy as an easy way to generate / store passwords.

Some security precautions have been taken but I could definitely use some advice on how to make
this program better!

The master user password and the  individual domain passwords are encrypted on disk.
The master password is a one way encryption using bcrypt, and the domain passwords take
the plain text master password (after it has been verified) combined with username to create key to encrypt the
domain passwords.

## Installation From Source
```
go build
```

Pappy stores its information in a sqlite database.  See the config package to understand
how to modify the database file location.

## Usage
```
pappy -h
```

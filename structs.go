package main

import (

	_ "github.com/go-sql-driver/mysql"
)

type verifyCode struct {
	Email      string `db:"email"`
	VerifyCode string `db:"verify_code"`
}

type responseData struct {
	Code    int    `json:"code"`
	Message string `json:"msg"`
}

type loginResponseData struct {
	Code    int    `json:"code"`
	Message string `json:"msg"`
	Data    string `json:"data"`
}

type account struct {
	Email    string `db:"email"`
	Password string `db:"password"`
}


type userOfflineData struct {
	Type  string `json:"type"`
	Email string `json:"email"`
}

type user struct {
	Email      string `json:"email"`
	Broadcast  string `json:"broadcast"`
	Lng        string `json:"lng"`
	Lat        string `json:"lat"`
	TimeUpdate string `json:"timeupdate"`
}

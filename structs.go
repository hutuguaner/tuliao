package main


import (
	"database/sql"
	_ "github.com/go-sql-driver/mysql"
)


type verifyCode struct {
	ID         int64  `db:"id"`
	Email      string `db:"email"`
	VerifyCode string `db:"verify_code"`
}



type responseData struct {
	Code    int    `json:"code"`
	Message string `json:"msg"`
}

type loginResponseData struct{
	Code int `json:"code"`
	Message string `json:"msg"`
	Data string `json:"data"`
}

type account struct{
	Email string `db:"email"`
	Password string `db:"password"`	
}

type version struct {
	VersionCode int64          `db:"version_code" json:"versionCode"`
	VersionName sql.NullString `db:"version_name" json:"versionName"`
	ForceUpdate int8           `db:"force_update" json:"forceUpdate"`
	VersionPwd  sql.NullString `db:"version_pwd" json:"versionPwd"`
	VersionDes  sql.NullString `db:"version_des" json:"versionDes"`
	VersionType sql.NullString `db:"version_type" json:"versionType"`
}
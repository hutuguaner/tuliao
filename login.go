package main

import (
	"encoding/json"
	"fmt"
	"net/http"
)

func login(w http.ResponseWriter, r *http.Request) {
	decoder := json.NewDecoder(r.Body)
	var params map[string]string
	decoder.Decode(&params)
	email := params["email"]
	//pwd:=params["password"]

	//判断当前是否又相同 用户名的账号在登录
	isUsed := emailIsUsed(email)

	var loginResponseData loginResponseData
	if isUsed {
		loginResponseData.Code = 1
		loginResponseData.Message = "当前用户名正在被使用中，请更换其他用户名登录"
		loginResponseData.Data = email
	} else {

		err := insertEmailInDB(email)
		if err != nil {
			loginResponseData.Code = 0
			loginResponseData.Message = "登录失败"
			loginResponseData.Data = email
		} else {
			loginResponseData.Code = 0
			loginResponseData.Message = "登录成功"
			loginResponseData.Data = email
		}
	}

	b, err := json.Marshal(loginResponseData)
	if err != nil {
		fmt.Fprintf(w, `{"code":1,"msg":"登录失败，json异常"}`)
		return
	}
	fmt.Fprintf(w, string(b))

}

func insertEmailInDB(email string) error {
	if !hasDbInit {
		initDb()
	}
	insForm, err := myDb.Prepare("insert into user(email)values(?)")
	if err != nil {
		return err
	}
	insForm.Exec(email)
	return nil
}

func emailIsUsed(email string) bool {
	if !hasDbInit {
		initDb()
	}
	row := myDb.QueryRow("select email from user where email=?", email)
	emailQuery := ""
	row.Scan(&emailQuery)
	fmt.Println(emailQuery)
	if emailQuery == email {
		return true
	} else {
		return false
	}

}

func queryPwd(email string) (string, error) {
	if !hasDbInit {
		initDb()
	}
	account := new(account)
	row := myDb.QueryRow("select password from account where email=?", email)
	err := row.Scan(&account.Password)
	if err != nil {
		return "", err
	}
	return account.Password, nil
}

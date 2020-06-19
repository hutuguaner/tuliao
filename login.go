package main

import (
	"encoding/json"
	"fmt"
	"net/http"
	"crypto/md5"
	"encoding/hex"
)

func login(w http.ResponseWriter, r *http.Request) {
	decoder := json.NewDecoder(r.Body)
	var params map[string]string
	decoder.Decode(&params)
	email := params["email"]
	pwd:=params["password"]

	//查看邮箱对应帐号是否存在
	if !isAccountExist(email) {
		fmt.Fprintf(w, `{"code":1,"msg":"帐号不存在"}`)
		return
	}

	//从数据库中 取出 密码 判断密码是否一致
	pwdInDB := queryPwdFromDB(email)

	//密码加密
	h := md5.New()
	h.Write([]byte(pwd))
	pwdMd5 := hex.EncodeToString(h.Sum(nil))


	var loginResponseData loginResponseData
	if pwdMd5!=pwdInDB {
		loginResponseData.Code = 1
		loginResponseData.Message = "密码错误"
		loginResponseData.Data = email
	} else {
		
		loginResponseData.Code = 0
			loginResponseData.Message = "登录成功"
			loginResponseData.Data = email
	}

	b, err := json.Marshal(loginResponseData)
	if err != nil {
		fmt.Fprintf(w, `{"code":1,"msg":"登录失败，json异常"}`)
		return
	}
	fmt.Fprintf(w, string(b))

}



func queryPwdFromDB(email string) string {
	if !hasDbInit {
		initDb()
	}
	row := myDb.QueryRow("select password from account where email=?", email)
	pwdQuery := ""
	row.Scan(&pwdQuery)
	return pwdQuery

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

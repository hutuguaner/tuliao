package main

import (
	"crypto/md5"
	"encoding/hex"
	"encoding/json"
	"fmt"
	"net/http"
)


func login(w http.ResponseWriter,r *http.Request){
	decoder:=json.NewDecoder(r.Body)
	var params map[string]string
	decoder.Decode(&params)
	email:=params["email"]
	pwd:=params["password"]

	//判断用户是否存在
	if !isAccountExist(email) {
		fmt.Fprintf(w,`{"code":1,"msg":"帐号不存在"}`)
		return
	}

	//从数据库中查询 邮箱对应的密码
	pwdMd5InDb,err :=queryPwd(email)
	if err!=nil {
		fmt.Fprintf(w,`{"code":1,"msg":"用户存在，当时验证密码是否正确失败，从数据库获取密码失败"}`)
		return
	}
	//用户传入密码加密
	h:=md5.New()
	h.Write([]byte(pwd))
	pwdMd5FromUser:=hex.EncodeToString(h.Sum(nil))

	//判断密码是否正确
	if pwdMd5InDb!=pwdMd5FromUser {
		fmt.Fprintf(w,`{"code":1,"msg":"密码错误"}`)
		return
	}

	var loginResponseData loginResponseData
	loginResponseData.Code=0
	loginResponseData.Message="登录成功"
	loginResponseData.Data = email
	b,err:=json.Marshal(loginResponseData)
	if err!=nil {
		fmt.Fprintf(w,`{"code":1,"msg":"登录失败，json异常"}`)
		return
	}
	fmt.Fprintf(w,string(b))

}


func queryPwd(email string)(string,error)  {
	if !hasDbInit{
		initDb()
	}
	account:=new(account)
	row:=myDb.QueryRow("select password from account where email=?",email)
	err:=row.Scan(&account.Password)
	if err!=nil {
		return "",err
	}
	return account.Password,nil
}
package main

import (
	"encoding/json"
	"fmt"
	"math/rand"
	"net/http"
	"time"

	_ "github.com/go-sql-driver/mysql"
	"gopkg.in/gomail.v2"
)

func getVerifyCodeByEmail(w http.ResponseWriter, r *http.Request) {

	decoder := json.NewDecoder(r.Body)
	var params map[string]string
	decoder.Decode(&params)
	email := params["email"]
	m := gomail.NewMessage()
	m.SetHeader("From", "553761200@qq.com")
	m.SetHeader("To", email)
	m.SetHeader("Subject", "Tuliao VerifyCode")
	verifyCode := getRandomString(6)

	//
	index, err := insertVerifyCodeIntoDb(email, verifyCode)
	if err != nil || index == -1 {

		fmt.Fprintf(w, `{"code":1,"msg":"生成验证码失败"}`)
		return
	}

	//
	emailBody := "Hello,your code is : " + verifyCode
	m.SetBody("text/html", emailBody)
	d := gomail.NewDialer("smtp.qq.com", 465, "553761200@qq.com", "lyenscpsagblbdej")

	var responseData responseData

	if err := d.DialAndSend(m); err != nil {
		responseData.Code = 1
		responseData.Message = err.Error()
		panic(err)
	} else {
		responseData.Code = 0
		responseData.Message = "获取验证码成功"
	}

	b, err := json.Marshal(responseData)
	if err != nil {
		fmt.Fprintf(w, `{"code":1,"msg":"获取验证码失败"}`)
	} else {
		fmt.Fprintf(w, string(b))
	}

}

//随机生成字符串
func getRandomString(l int) string {
	str := "0123456789abcdefghijklmnopqrstuvwxyz"
	bytes := []byte(str)
	result := []byte{}
	r := rand.New(rand.NewSource(time.Now().UnixNano()))
	for i := 0; i < l; i++ {
		result = append(result, bytes[r.Intn(len(bytes))])
	}
	return string(result)
}

func insertVerifyCodeIntoDb(email string, verifyCode string) (int64, error) {
	if !hasDbInit {
		initDb()
	}
	insForm, err := myDb.Prepare("insert into verify_code(email,verify_code)values(?,?)on duplicate key update verify_code=?")
	if err != nil {

		return -1, err
	}

	insForm.Exec(email, verifyCode, verifyCode)
	return 0, nil
}

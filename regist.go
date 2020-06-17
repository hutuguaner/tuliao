package main

import (
	"crypto/md5"
	"encoding/hex"
	"encoding/json"
	"fmt"
	"net/http"

	_ "github.com/go-sql-driver/mysql"
)

func regist(w http.ResponseWriter, r *http.Request) {
	decoder := json.NewDecoder(r.Body)
	var params map[string]string
	decoder.Decode(&params)
	email := params["email"]
	//verifyCode := params["verifyCode"]
	pwd := params["password"]

	//查看邮箱对应帐号是否存在
	if isAccountExist(email) {
		fmt.Fprintf(w, `{"code":1,"msg":"帐号已存在"}`)
		return
	}

	/* verifyCodeInDb, err := getVerifyCodeByEmailFromDB(email)
	if err != nil {
		fmt.Fprintf(w, `{"code":1,"msg":"验证码获取失败，请重新获取验证码"}`)
		return
	}

	if verifyCode != verifyCodeInDb {
		fmt.Fprintf(w, `{"code":1,"msg":"验证码错误，请输入正确验证码"}`)
		return
	} */

	//密码加密
	h := md5.New()
	h.Write([]byte(pwd))
	pwdMd5 := hex.EncodeToString(h.Sum(nil))

	//生成新用户 ，将用户信息插入数据库
	err2 := insertAccountIntoDb(email, pwdMd5)
	if err2 != nil {
		fmt.Fprintf(w, `{"code":1,"msg":"帐号插入数据库失败"}`)
		return
	}

	fmt.Fprintf(w, `{"code":0,"msg":"注册成功"}`)

}

func isAccountExist(email string) bool {

	if !hasDbInit {
		initDb()
	}
	rows, _ := myDb.Query("select * from account where email=?", email)
	return rows.Next()
}

func insertAccountIntoDb(email string, pwdMd5 string) error {
	if !hasDbInit {
		initDb()
	}
	insForm, err := myDb.Prepare("insert into account(email,password)values(?,?)on duplicate key update password=?")
	if err != nil {
		return err
	}
	insForm.Exec(email, pwdMd5, pwdMd5)
	return nil
}

func getVerifyCodeByEmailFromDB(email string) (string, error) {
	if !hasDbInit {
		initDb()
	}
	verifyCode := new(verifyCode)
	row := myDb.QueryRow("select * from verify_code where email=?", email)

	err := row.Scan(&verifyCode.Email, &verifyCode.VerifyCode)
	if err != nil {
		return "", err
	}
	return verifyCode.VerifyCode, err
}

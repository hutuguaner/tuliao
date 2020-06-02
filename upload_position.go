package main

import (
	"encoding/json"
	"fmt"
	"net/http"
	"time"
)

func uploadPosition(w http.ResponseWriter, r *http.Request) {
	decoder := json.NewDecoder(r.Body)
	var params map[string]string
	decoder.Decode(&params)
	email := params["email"]
	lng := params["lng"]
	lat := params["lat"]

	now := time.Now()
	sec := now.Unix()

	var uploadPositionResponse uploadPositionResponse

	userAvailable:=checkUserIsExistAndAvailable(email)

	if !userAvailable{
		uploadPositionResponse.Code = 2
		uploadPositionResponse.Message = "连接超时，请重新登录"
		uploadPositionResponse.Data =""
	}else{
		err1 := updatePositionDB(lng, lat, sec, email)

		if err1 != nil {
			uploadPositionResponse.Code = 1
			uploadPositionResponse.Message = "上传位置数据失败"
			uploadPositionResponse.Data = err1.Error()
		} else {
			uploadPositionResponse.Code = 0
			uploadPositionResponse.Message = "上传位置数据成功"
			uploadPositionResponse.Data = ""
		}
	}

	
	var jsonData []byte
	jsonData, err := json.Marshal(uploadPositionResponse)
	if err != nil {
		fmt.Println(jsonData)
	}

	fmt.Fprintf(w, string(jsonData))

}

type uploadPositionResponse struct {
	Code    int    `json:"code"`
	Message string `json:"message"`
	Data    string `json:"data"`
}

//更新时间
func updatePositionDB(lng string, lat string, stamp int64, email string) error {

	if !hasDbInit {
		initDb()
	}
	insForm, err := myDb.Prepare("update user set lng=?,lat=?,time_update=? where email=?")
	if err != nil {

		return err
	}

	insForm.Exec(lng, lat, stamp, email)
	if err != nil {
		fmt.Println(err.Error())

	}

	return nil
}

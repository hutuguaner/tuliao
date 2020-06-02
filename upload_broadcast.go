package main

import (
	"encoding/json"
	"fmt"
	"net/http"
	"time"
)

func uploadBroadcast(w http.ResponseWriter, r *http.Request) {
	decoder := json.NewDecoder(r.Body)
	var params map[string]string
	decoder.Decode(&params)
	email := params["email"]
	broadcast := params["broadcast"]

	now := time.Now()
	sec := now.Unix()

	var uploadBroadcastResponse uploadBroadcastResponse

	userAvailable := checkUserIsExistAndAvailable(email)
	if !userAvailable{
		uploadBroadcastResponse.Code = 2
		uploadBroadcastResponse.Message = "连接超时，请重新登录"
		uploadBroadcastResponse.Data = ""
	}else{
		err1 := updateBroadcastDB(broadcast, sec, email)

	
		if err1 != nil {
			uploadBroadcastResponse.Code = 1
			uploadBroadcastResponse.Message = "上传广播失败"
			uploadBroadcastResponse.Data = err1.Error()
		} else {
			uploadBroadcastResponse.Code = 0
			uploadBroadcastResponse.Message = "上传广播成功"
			uploadBroadcastResponse.Data = ""
		}
	}

	
	var jsonData []byte
	jsonData, err := json.Marshal(uploadBroadcastResponse)
	if err != nil {
		fmt.Println(err.Error())
	}

	fmt.Fprintf(w, string(jsonData))

}

type uploadBroadcastResponse struct {
	Code    int    `json:"code"`
	Message string `json:"message"`
	Data    string `json:"data"`
}

//更新时间
func updateBroadcastDB(broadcast string, stamp int64, email string) error {

	if !hasDbInit {
		initDb()
	}
	insForm, err := myDb.Prepare("update user set broadcast=?,time_update=? where email=?")
	if err != nil {

		return err
	}

	insForm.Exec(broadcast, stamp, email)
	if err != nil {
		//fmt.Println(err.Error())

	}

	return nil
}

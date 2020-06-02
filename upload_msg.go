package main

import (
	"encoding/json"
	"fmt"
	"net/http"
	"time"
)

func uploadMsg(w http.ResponseWriter, r *http.Request) {
	decoder := json.NewDecoder(r.Body)
	var params map[string]string
	decoder.Decode(&params)
	from := params["from"]
	content := params["content"]
	to := params["to"]

	now := time.Now()
	sec := now.Unix()

	var insertMsgResponse insertMsgResponse

	userAvailable := checkUserIsExistAndAvailable(from)

	if !userAvailable {
		insertMsgResponse.Code = 2
		insertMsgResponse.Message = "连接超时，请重新登录"
		insertMsgResponse.Data = ""
	} else {
		err1 := insertMsgIntoDB(from, to, content, sec)

		if err1 != nil {
			insertMsgResponse.Code = 1
			insertMsgResponse.Message = "上传消息失败"
			insertMsgResponse.Data = err1.Error()
		} else {
			insertMsgResponse.Code = 0
			insertMsgResponse.Message = "上传消息成功"
			insertMsgResponse.Data = ""
		}
	}

	var jsonData []byte
	jsonData, err := json.Marshal(insertMsgResponse)
	if err != nil {
		fmt.Println(jsonData)
	}

	fmt.Fprintf(w, string(jsonData))

}

type insertMsgResponse struct {
	Code    int    `json:"code"`
	Message string `json:"message"`
	Data    string `json:"data"`
}

//更新时间
func insertMsgIntoDB(from string, to string, content string, time int64) error {

	if !hasDbInit {
		initDb()
	}

	insForm, err := myDb.Prepare("insert msg set from_email=?,content=?,to_email=?,time=?")
	if err != nil {

		return err
	}

	insForm.Exec(from, content, to, time)
	if err != nil {
		fmt.Println(err.Error())

	}

	return nil
}

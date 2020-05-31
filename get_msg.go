package main

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"net/http"
)

func getMsgs(w http.ResponseWriter, r *http.Request) {
	decoder := json.NewDecoder(r.Body)
	var params map[string]string
	decoder.Decode(&params)
	email := params["email"]

	msgs, err := getMsgsFromDB(email)
	var getMsgsResponse getMsgsResponse
	if err != nil {
		getMsgsResponse.Code = 1
		getMsgsResponse.Message = "获取消息失败"
		getMsgsResponse.Data = nil
	} else {
		getMsgsResponse.Code = 0
		getMsgsResponse.Message = "获取消息成功"
		getMsgsResponse.Data = msgs
	}
	var jsonData []byte
	jsonData, err = json.Marshal(getMsgsResponse)
	if err != nil {
		fmt.Println(jsonData)
	}

	fmt.Fprintf(w, string(jsonData))

}

type getMsgsResponse struct {
	Code    int    `json:"code"`
	Message string `json:"message"`
	Data    []msg  `json:"data"`
}
type msg struct {
	From    string `json:"from"`
	To      string `json:"to"`
	Content string `json:"content"`
	Time    string `json:"time"`
}

func getMsgsFromDB(email string) ([]msg, error) {
	if !hasDbInit {
		initDb()
	}
	msgs := []msg{}
	results, err := myDb.Query("select * from msg where to_email =? ", email)
	if err != nil {
		return msgs, err
	}

	for results.Next() {
		var m msg
		var from, to, content, time sql.NullString
		err = results.Scan(&from, &content, &to, &time)

		if err != nil {
			fmt.Println(err.Error())
			return msgs, err
		}
		m.From = from.String
		m.To = to.String
		m.Content = content.String
		m.Time = time.String

		msgs = append(msgs, m)
	}
	//查询完了，将查询过的消息删掉
	insForm, err := myDb.Prepare("delete from msg where to_email=?")
	if err != nil {
		return msgs, err
	}
	insForm.Exec(email)
	return msgs, nil
}

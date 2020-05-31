package main

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"net/http"
	"time"
)

func getUsers(w http.ResponseWriter, r *http.Request) {
	decoder := json.NewDecoder(r.Body)
	var params map[string]string
	decoder.Decode(&params)
	email := params["email"]

	now := time.Now()
	sec := now.Unix()

	err1 := updateTimeDB(sec, email)
	users, err2 := getUserFromDB()
	var getUsersResponse getUsersResponse
	if err1 != nil || err2 != nil {
		getUsersResponse.Code = 1
		getUsersResponse.Message = "获取用户失败"
		getUsersResponse.Data = nil
	} else {
		getUsersResponse.Code = 0
		getUsersResponse.Message = "获取用户成功"
		getUsersResponse.Data = users
	}
	var jsonData []byte
	jsonData, err := json.Marshal(getUsersResponse)
	if err != nil {
		fmt.Println(jsonData)
	}

	fmt.Fprintf(w, string(jsonData))

}

type getUsersResponse struct {
	Code    int    `json:"code"`
	Message string `json:"message"`
	Data    []user `json:"data"`
}

//更新时间
func updateTimeDB(stamp int64, email string) error {
	fmt.Println(email)
	if !hasDbInit {
		initDb()
	}
	insForm, err := myDb.Prepare("update user set time_update=? where email=?")
	if err != nil {

		return err
	}
	result, _ := insForm.Exec(stamp, email)
	id, _ := result.RowsAffected()
	fmt.Println(id)
	return nil
}

func getUserFromDB() ([]user, error) {
	if !hasDbInit {
		initDb()
	}
	users := []user{}
	results, err := myDb.Query("select * from user")
	if err != nil {
		return users, err
	}

	for results.Next() {
		var u user
		var email, broadcast, lng, lat, time sql.NullString
		err = results.Scan(&email, &broadcast, &lng, &lat, &time)

		if err != nil {
			fmt.Println(err.Error())
			return users, err
		}
		u.Email = email.String
		u.Broadcast = broadcast.String
		u.Lng = lng.String
		u.Lat = lat.String
		u.TimeUpdate = time.String
		users = append(users, u)
	}

	return users, nil
}

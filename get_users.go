package main

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"net/http"
	"strconv"
	"time"
)

func getUsers(w http.ResponseWriter, r *http.Request) {
	decoder := json.NewDecoder(r.Body)
	var params map[string]string
	decoder.Decode(&params)
	//email := params["email"]


	var getUsersResponse getUsersResponse

	users, err2 := getUserFromDB()

		if err2 != nil {
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

	}

	fmt.Fprintf(w, string(jsonData))

}


//将超时 用户 删除
func deleteUserFromDB(email string) bool {
	if !hasDbInit {
		initDb()
	}
	insForm, err := myDb.Prepare("delete from user where email=?")
	if err != nil {
		return false
	} else {
		insForm.Exec(email)
		return true
	}
}

type getUsersResponse struct {
	Code    int    `json:"code"`
	Message string `json:"message"`
	Data    []user `json:"data"`
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
		var email, broadcast, lng, lat, timeUpdate sql.NullString
		err = results.Scan(&email, &broadcast, &lng, &lat, &timeUpdate)

		if err != nil {
			fmt.Println(err.Error())
			return users, err
		}
		u.Email = email.String
		u.Broadcast = broadcast.String
		u.Lng = lng.String
		u.Lat = lat.String
		u.TimeUpdate = timeUpdate.String

		timeUpdateInt64, err := strconv.ParseInt(u.TimeUpdate, 10, 64)
		if err != nil {
			panic(err)
		} else {
			now := time.Now()
			dis, _ := time.ParseDuration(connectTimeOut)
			disBefore := now.Add(dis).Unix()
			if timeUpdateInt64 < disBefore {
				//超时
				deleteUserFromDB(u.Email)

			} else {
				//正常
				users = append(users, u)
			}
		}

	}

	return users, nil
}

var connectTimeOut = "-1m"

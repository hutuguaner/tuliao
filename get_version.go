package main

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"net/http"
)

func getVersion(w http.ResponseWriter, r *http.Request) {
	

	version,err := getVersionFromDB()
	var getVersionResponse getVersionResponse
	if err!=nil {
		getVersionResponse.Code = 1
		getVersionResponse.Message = "获取版本信息失败"
		getVersionResponse.Data= version
	}else{
		getVersionResponse.Code = 0
		getVersionResponse.Message = "获取版本信息成功"
		getVersionResponse.Data = version
	}
	
	
	var jsonData []byte
	jsonData, err = json.Marshal(getVersionResponse)
	if err != nil {
		fmt.Println(jsonData)
	}

	fmt.Fprintf(w, string(jsonData))

}

type getVersionResponse struct {
	Code    int    `json:"code"`
	Message string `json:"message"`
	Data    version  `json:"data"`
}
type version struct {
	VersionCode    string `json:"versionCode"`
	VersionName      string `json:"versionName"`
	ForceUpdate string `json:"forceUpdate"`
	VersionDesc    string `json:"versionDesc"`
}

func getVersionFromDB() (version,error) {
	if !hasDbInit {
		initDb()
	}
	var v version
	var versionCode,versionName,forceUpdate,versionDesc sql.NullString
	row := myDb.QueryRow("select * from version where version_code=(select max(version_code)from version)")
	
	err:=row.Scan(&versionCode,&versionName,&forceUpdate,&versionDesc)
	if err!=nil {
		fmt.Println(err.Error())
		return v,err
	}
	v.VersionCode = versionCode.String
	v.VersionName = versionName.String
	v.ForceUpdate = forceUpdate.String
	v.VersionDesc = versionDesc.String
	

	return v,nil
}

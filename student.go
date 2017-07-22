package main

import (
	"fmt"
	"net/http"

	"github.com/julienschmidt/httprouter"
)

/*
{
  "studentID":"09882342",
  "name":"张三",
  "phone":"13900000000",
  "studentOpenID":"adsf2324sdfa",
  "email":"aaa@163.com",
  "school":"实验小学",
  "grade":"三年级",
  "capibility":"优秀",
  "address":"福田小区",
  "motherName":"张二",
  "motherPhone":"13800000000",
  "motherOpenID":"adsf2324sdfa",
  "fatherName":"张二",
  "fatherPhone":"13800000000",
  "fatherOpenID":"adsf2324sdfa",
  "lastUpdateTime":"2017-07-11"
}
*/

/*Student is the GO struct for student JSON */
type Student struct {
	StudentID      string `json:"studentID"`
	Name           string `json:"name"`
	Phone          string `json:"phone"`
	StudentOpenID  string `json:"studentOpenID"`
	Email          string `json:"email"`
	School         string `json:"school"`
	Grade          string `json:"grade"`
	Capibility     string `json:"capibility"`
	Address        string `json:"address"`
	MotherName     string `json:"motherName"`
	MotherPhone    string `json:"motherPhone"`
	MotherOpenID   string `json:"motherOpenID"`
	FatherName     string `json:"fatherName"`
	FatherPhone    string `json:"fatherPhone"`
	FatherOpenID   string `json:"fatherOpenID"`
	LastUpdateTime string `json:"lastUpdateTime"`
}

func GetStudent(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
	fmt.Fprintf(w, "hello, %s!\n", ps.ByName("name"))
}

func CreateStudent(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
	fmt.Fprintf(w, "hello, %s!\n", ps.ByName("name"))
}

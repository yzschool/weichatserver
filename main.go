package main

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"

	"github.com/julienschmidt/httprouter"
)

type Class struct {
	ClassID        string   `json:"classID"`
	Name           string   `json:"name"`
	ClassStartTime string   `json:"classStartTime"`
	ClassEndTime   string   `json:"classEndTime"`
	ClassPeroid    string   `json:"classPeroid"`
	ClassTime      string   `json:"classTime"`
	ClassLevel     string   `json:"classLevel"`
	City           string   `json:"city"`
	District       string   `json:"district"`
	Building       string   `json:"building"`
	Latitude       float64  `json:"latitude"`
	Longitude      float64  `json:"longitude"`
	CreatedBy      string   `json:"createdBy"`
	ContactTel     string   `json:"contactTel"`
	CreatedTime    string   `json:"createdTime"`
	LastUpdateTime string   `json:"lastUpdateTime"`
	Price          string   `json:"price"`
	Students       []string `json:"students"`
}

type Student struct {
	StudentID      string   `json:"studentID"`
	Name           string   `json:"name"`
	Phone          string   `json:"phone"`
	Email          string   `json:"email"`
	School         string   `json:"school"`
	Grade          string   `json:"grade"`
	Capibility     string   `json:"capibility"`
	Address        string   `json:"address"`
	Parent         string   `json:"parent"`
	ParentPhone    string   `json:"parentPhone"`
	Class          []string `json:"class"`
	LastUpdateTime string   `json:"lastUpdateTime"`
}

func Index(w http.ResponseWriter, r *http.Request, _ httprouter.Params) {
	fmt.Fprint(w, "Welcome!\n")
}

func GetClass(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
	class := Class{ClassID: "US123", Name: "English 1", ClassStartTime: "2017-07-15", ClassEndTime: "2017-07-30", City: "shenzhen", District: "futian", Building: "xiangmu", CreatedBy: "xiao li"}
	b, err := json.Marshal(class)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	w.Header().Set("Content-Type", "application/json")
	w.Write(b)
}

func CreateClass(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
	//fmt.Fprintf(w, "hello, %s!\n", ps.ByName("name"))

	decoder := json.NewDecoder(r.Body)
	var new_class Class
	err := decoder.Decode(&new_class)

	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	fmt.Println(w, new_class)
	//w.Write("success")
	/*
		decoder := json.NewDecoder(r.Body)
		var t test_struct
		err := decoder.Decode(&t)
		if err != nil {
			panic(err)
		}
		defer r.Body.Close()
		log.Println(t.Test)
	*/
}

func GetStudent(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
	fmt.Fprintf(w, "hello, %s!\n", ps.ByName("name"))
}

func CreateStudent(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
	fmt.Fprintf(w, "hello, %s!\n", ps.ByName("name"))
}

func main() {
	router := httprouter.New()
	router.GET("/", Index)

	router.GET("/weichat/class", GetClass)
	router.POST("/weichat/class", CreateClass)
	router.GET("/weichat/student", GetStudent)
	router.POST("/weichat/student", CreateStudent)

	log.Fatal(http.ListenAndServe("0.0.0.0:8080", router))
}

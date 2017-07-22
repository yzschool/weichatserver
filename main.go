package main

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"

	mgo "gopkg.in/mgo.v2"
	"gopkg.in/mgo.v2/bson"

	"github.com/julienschmidt/httprouter"
)

type Class struct {
	ClassID        string   `json:"classID"`
	ClassName      string   `json:"className"`
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
	Creater        string   `json:"creater"`
	CreaterOpenID  string   `json:"createrOpenID"`
	CreaterTel     string   `json:"createrTel"`
	CreatTime      string   `json:"creatTime"`
	Grade          string   `json:"grade"`
	TeacherTel     string   `json:"teacherTel"`
	Price          string   `json:"price"`
	StudentsID     []string `json:"studentsID"`
}

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

func Index(w http.ResponseWriter, r *http.Request, _ httprouter.Params) {
	fmt.Fprint(w, "Welcome!\n")
}

func ErrorWithJSON(w http.ResponseWriter, message string, code int) {
	w.Header().Set("Content-Type", "application/json; charset=utf-8")
	w.WriteHeader(code)
	fmt.Fprintf(w, "{message: %q}", message)
}

func ResponseWithJSON(w http.ResponseWriter, json []byte, code int) {
	w.Header().Set("Content-Type", "application/json; charset=utf-8")
	w.WriteHeader(code)
	w.Write(json)
}

func GetClass(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
	session := mgoSession.Copy()
	defer session.Close()

	c := session.DB("yzschool").C("class")

	var classes []Class
	err := c.Find(bson.M{}).All(&classes)
	if err != nil {
		ErrorWithJSON(w, "Database error", http.StatusInternalServerError)
		log.Println("Failed get all classes: ", err)
		return
	}

	respBody, err := json.MarshalIndent(classes, "", "  ")
	if err != nil {
		log.Fatal(err)
	}

	ResponseWithJSON(w, respBody, http.StatusOK)
}

func CreateClass(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
	session := mgoSession.Copy()
	defer session.Close()

	var class Class
	decoder := json.NewDecoder(r.Body)
	err := decoder.Decode(&class)
	if err != nil {
		ErrorWithJSON(w, "Incorrect body", http.StatusBadRequest)
		return
	}

	c := session.DB("yzschool").C("class")

	err = c.Insert(class)
	if err != nil {
		if mgo.IsDup(err) {
			ErrorWithJSON(w, "Book with this ISBN already exists", http.StatusBadRequest)
			return
		}

		ErrorWithJSON(w, "Database error", http.StatusInternalServerError)
		log.Println("Failed insert book: ", err)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)

}

func GetStudent(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
	fmt.Fprintf(w, "hello, %s!\n", ps.ByName("name"))
}

func CreateStudent(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
	fmt.Fprintf(w, "hello, %s!\n", ps.ByName("name"))
}

var mgoSession *mgo.Session

func main() {

	session, err := mgo.Dial("localhost:27017")
	if err != nil {
		fmt.Println("connect to mongodb failure", err)
		panic(err)
	}
	defer session.Close()

	session.SetMode(mgo.Monotonic, true)
	mgoSession = session
	//ensureIndex(session)

	router := httprouter.New()
	router.GET("/", Index)

	router.GET("/weichat/class", GetClass)
	router.POST("/weichat/class", CreateClass)
	router.GET("/weichat/student", GetStudent)
	router.POST("/weichat/student", CreateStudent)

	log.Fatal(http.ListenAndServe("0.0.0.0:8080", router))
}

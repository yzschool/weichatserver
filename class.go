package main

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"

	"gopkg.in/mgo.v2"
	"gopkg.in/mgo.v2/bson"

	"github.com/julienschmidt/httprouter"
)

/*
example JSON class
{
  "classID":"1202434354545",
  "className":"yz english AA",
  "classStartTime":"2017-07-11",
  "classEndTime":"2017-09-11",
  "classPeroid":"每周三晚上",
  "classTime":"19:30-20:30",
  "classLevel":"中等",
  "city":"深圳",
  "district":"福田",
  "building":"香蜜湖小区",
  "latitude":114.026694,
  "longitude":22.549416,
  "creater":"xiao lee",
  "createrOpenID":"adsf2324sdfa",
  "createrTel":"13900000000",
  "creatTime":"2017-07-01",
  "grade":"三年纪",
  "teacherTel":"13900000000",
  "price":"1600",
  "studentsID":["09882342", "09882342", "09882342"]
}
*/

/*Class is GO struct for Class JSON */
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

func GetClassByName(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
	session := mgoSession.Copy()
	defer session.Close()

	id := ps.ByName("clasid")

	c := session.DB("yzschool").C("class")
	fmt.Println("GetClassByName classid is ", id)

	var class Class
	err := c.Find(bson.M{"classID": id}).One(&class)
	if err != nil {
		ErrorWithJSON(w, "Database error", http.StatusInternalServerError)
		log.Println("Failed find book: ", err)
		return
	}

	if class.ClassID == "" {
		ErrorWithJSON(w, "Class not found", http.StatusNotFound)
		return
	}

	respBody, err := json.MarshalIndent(class, "", "  ")
	if err != nil {
		log.Fatal(err)
	}

	ResponseWithJSON(w, respBody, http.StatusOK)
}

func UpdateClass(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
	session := mgoSession.Copy()
	defer session.Close()

	id := ps.ByName("classid")

	var class Class
	decoder := json.NewDecoder(r.Body)
	err := decoder.Decode(&class)
	if err != nil {
		ErrorWithJSON(w, "Incorrect body", http.StatusBadRequest)
		return
	}

	c := session.DB("yzschool").C("class")

	err = c.Update(bson.M{"classID": id}, &class)
	if err != nil {
		switch err {
		default:
			ErrorWithJSON(w, "Database error", http.StatusInternalServerError)
			log.Println("Failed update book: ", err)
			return
		case mgo.ErrNotFound:
			ErrorWithJSON(w, "Book not found", http.StatusNotFound)
			return
		}
	}

	w.WriteHeader(http.StatusNoContent)
}

func DeleteClass(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
	session := mgoSession.Copy()
	defer session.Close()

	id := ps.ByName("classid")

	c := session.DB("yzschool").C("class")

	err := c.Remove(bson.M{"classID": id})
	if err != nil {
		switch err {
		default:
			ErrorWithJSON(w, "Database error", http.StatusInternalServerError)
			log.Println("Failed delete book: ", err)
			return
		case mgo.ErrNotFound:
			ErrorWithJSON(w, "Book not found", http.StatusNotFound)
			return
		}
	}

	w.WriteHeader(http.StatusNoContent)
}

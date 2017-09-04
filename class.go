package main

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"

	"gopkg.in/mgo.v2"
	"gopkg.in/mgo.v2/bson"

	"time"

	"github.com/julienschmidt/httprouter"
	"github.com/nu7hatch/gouuid"
)

/*
example JSON class
{
  "classid":"1202434354545",
  "classname":"yz english AA",
  "starttime":"2017-07-11",
  "usertype":"男",
  "endtime":"2017-09-11",
  "peroid":"每周三晚上",
  "classtime":"19:30-20:30",
  "level":"中等",
  "city":"深圳",
  "district":"福田",
  "building":"香蜜湖小区",
  "latitude":114.026694,
  "longitude":22.549416,
  "creater":"xiao lee",
  "openid":"adsf2324sdfa",
  "tel":"13900000000",
  "creattime":"2017-07-01",
  "grade":"三年纪",
  "teachertel":"13900000000",
  "price":"1600",
  "studentsid":["09882342", "09882342", "09882342"]
}
*/

/*Class is GO struct for Class JSON */
type Class struct {
	Classid    string   `json:"classid"`
	Classname  string   `json:"classname"`
	Starttime  string   `json:"starttime"`
	Usertype   string   `json:"usertype"`
	Endtime    string   `json:"endtime"`
	Peroid     string   `json:"peroid"`
	Classtime  string   `json:"classtime"`
	Level      string   `json:"level"`
	City       string   `json:"city"`
	District   string   `json:"district"`
	Building   string   `json:"building"`
	Latitude   float64  `json:"latitude"`
	Longitude  float64  `json:"longitude"`
	Creater    string   `json:"creater"`
	Openid     string   `json:"openid"`
	Tel        string   `json:"tel"`
	Creattime  string   `json:"creattime"`
	Grade      string   `json:"grade"`
	Teachertel string   `json:"teachertel"`
	Price      string   `json:"price"`
	Studentsid []string `json:"studentsid"`
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

	/* generate UUID for the class ID */
	u4, err := uuid.NewV4()
	if err != nil {
		fmt.Println("error:", err)
		return
	}
	class.Classid = u4.String()
	class.Creattime = time.Now().String()

	c := session.DB("yzschool").C("class")

	err = c.Insert(class)
	if err != nil {
		if mgo.IsDup(err) {
			ErrorWithJSON(w, "Class with this classid already exists", http.StatusBadRequest)
			return
		}

		ErrorWithJSON(w, "Database error", http.StatusInternalServerError)
		log.Println("Failed insert book: ", err)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)

}

func GetClassByID(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
	session := mgoSession.Copy()
	defer session.Close()

	id := ps.ByName("classid")

	c := session.DB("yzschool").C("class")
	fmt.Println("GetClassByID classid is ", id)

	var class Class
	err := c.Find(bson.M{"classid": id}).One(&class)
	if err != nil {
		ErrorWithJSON(w, "Database error", http.StatusInternalServerError)
		log.Println("Failed find book: ", err)
		return
	}

	if class.Classid == "" {
		ErrorWithJSON(w, "Class not found", http.StatusNotFound)
		return
	}

	respBody, err := json.MarshalIndent(class, "", "  ")
	if err != nil {
		log.Fatal(err)
	}

	ResponseWithJSON(w, respBody, http.StatusOK)
}

func GetClassByName(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
	session := mgoSession.Copy()
	defer session.Close()

	classname := ps.ByName("classname")

	c := session.DB("yzschool").C("class")
	fmt.Println("GetClassByClassName classname is ", classname)

	var class []Class
	err := c.Find(bson.M{"classname": bson.M{"$regex": bson.RegEx{classname, "i"}}}).All(&class)
	if err != nil  || class == nil{	
		ErrorWithJSON(w, "Database error", http.StatusInternalServerError)
		return
	}

	if class[0].Classname == "" {
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
	fmt.Println("classid is: " + id)

	var student Student
	decoder := json.NewDecoder(r.Body)
	err := decoder.Decode(&student)
	if err != nil {
		ErrorWithJSON(w, "Incorrect body", http.StatusBadRequest)
		return
	}

	var studentID string
	studentID = GetStudentByName(student.Name, student.Phone)
	if studentID == "" {
		studentID = add_student(student)
	}

	var class Class
	c := session.DB("yzschool").C("class")
	err = c.Find(bson.M{"classid": id}).One(&class)
	if err != nil {
		ErrorWithJSON(w, "Database error", http.StatusInternalServerError)
		fmt.Println("Failed find class: ", err)
		return
	}

	registered := false
	for _, value := range class.Studentsid {
		if value == studentID {
			registered = true
		}
	}

	if registered == false {
		fmt.Println("add class: " + id +" student: " + studentID )
		new_id := append(class.Studentsid, studentID)
		class.Studentsid = new_id

		err = c.Update(bson.M{"classid": id}, &class)
		if err != nil {
			switch err {
			default:
				ErrorWithJSON(w, "Database error", http.StatusInternalServerError)
				log.Println("Failed update book: ", err)
				return
			case mgo.ErrNotFound:
				ErrorWithJSON(w, "Class not found", http.StatusNotFound)
				return
			}
		}

		w.WriteHeader(http.StatusCreated)
	} else {
		w.WriteHeader(http.StatusAlreadyReported)
	}
}

func DeleteClass(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
	session := mgoSession.Copy()
	defer session.Close()

	id := ps.ByName("classid")

	c := session.DB("yzschool").C("class")

	err := c.Remove(bson.M{"classid": id})
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

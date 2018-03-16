
package main

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"time"
	"gopkg.in/mgo.v2/bson"

	"github.com/julienschmidt/httprouter"
)

/*
{
	"name": "张生",
	"gender": "男",
	"birthdate": "2017-07-11",
	"grade": "三年级",
	"tel": "13546984569",
	"school": "南山外国语小学",
	"address": "深圳-南山-科技园",
	"updatetime": "2018-01-01"
  }
  */

type CourseAppDataTable struct {
	Data []Courseapplication `json:"data"`
}

type Courseapplication struct {
	Name       string `json:"name"`
	Gender     string `json:"gender"`
	Birthdate  string `json:"birthdate"`
	Grade      string `json:"grade"`
	Tel        string `json:"tel"`
	School     string `json:"school"`
	Address    string `json:"address"`
	Updatetime string `json:"updatetime"`
}


/*SubmitApplication application submit application form */
func SubmitCourseApplication(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
	//fmt.Println("SubmitApplication get request")
	session := mgoSession.Copy()
	defer session.Close()

	c := session.DB("yzschool").C("courseapplication")
	
	var app Courseapplication
	decoder := json.NewDecoder(r.Body)
	err := decoder.Decode(&app)
	if err != nil {
		ErrorWithJSON(w, "Incorrect body", http.StatusBadRequest)
		return
	}
	app.Updatetime = time.Now().Format("2006-01-02 3:4:5 PM")

	err = c.Insert(app)
	if err != nil {
		ErrorWithJSON(w, "Database error", http.StatusInternalServerError)
		log.Println("Failed insert partner: ", err)
		return
	}

	fmt.Fprintf(w, "{message: success}")
	w.WriteHeader(http.StatusCreated)
}


func GetAllCourseApplicationDataTable(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
	session := mgoSession.Copy()
	defer session.Close()

	c := session.DB("yzschool").C("courseapplication")

	var dbApps []Courseapplication
	err := c.Find(bson.M{}).All(&dbApps)

	if err != nil || len(dbApps) == 0 {
		ErrorWithJSON(w, "Can not found in the database!", http.StatusNotFound)
		return
	}

	var data CourseAppDataTable
	data.Data = dbApps

	respBody, err := json.MarshalIndent(data, "", "  ")
	if err != nil {
		ErrorWithJSON(w, "Database error", http.StatusInternalServerError)
		log.Fatal(err)
		return
	}
	
	ResponseWithJSON(w, respBody, http.StatusOK)
}

func GetAllCourseApplication(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
	session := mgoSession.Copy()
	defer session.Close()

	c := session.DB("yzschool").C("courseapplication")

	var dbApps []Courseapplication
	err := c.Find(bson.M{}).All(&dbApps)

	if err != nil || len(dbApps) == 0 {
		ErrorWithJSON(w, "Can not found in the database!", http.StatusNotFound)
		return
	}

	respBody, err := json.MarshalIndent(dbApps, "", "  ")
	if err != nil {
		ErrorWithJSON(w, "Database error", http.StatusInternalServerError)
		log.Fatal(err)
		return
	}

	ResponseWithJSON(w, respBody, http.StatusOK)
}
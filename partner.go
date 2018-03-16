package main

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"time"

	"gopkg.in/mgo.v2"
	"gopkg.in/mgo.v2/bson"

	"github.com/julienschmidt/httprouter"
)

/*
{
  "partnername": "张生",
  "gender": "男",
  "birthdate": "2017-07-11",
  "marital": "未婚",
  "tel": "13546984569",
  "email": "566vjj@163.com",
  "address": "深圳-南山-科技园",
  "updatetime": "2018-01-01"
}
*/

type PartnerDataTable struct {
	Data []Partner `json:"data"`
}

type Partner struct {
	Partnername string `json:"partnername"`
	Gender      string `json:"gender"`
	Birthdate   string `json:"birthdate"`
	Marital     string `json:"marital"`
	Tel         string `json:"tel"`
	Email       string `json:"email"`
	Address     string `json:"address"`
	Updatetime  string `json:"updatetime"`
}

func ensureIndexPartner(c *mgo.Collection) {

	index := mgo.Index{
		Key:        []string{"tel"},
		Unique:     true,
		DropDups:   true,
		Background: true,
		Sparse:     true,
	}

	err := c.EnsureIndex(index)
	if err != nil {
		fmt.Println("EnsureIndex in mongodb failure", err)
		panic(err)
	}
}


/*SubmitApplication application submit application form */
func SubmitPartnerApplication(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
	//fmt.Println("SubmitApplication get request")
	session := mgoSession.Copy()
	defer session.Close()

	c := session.DB("yzschool").C("partner")
	ensureIndexPartner(c)

	var p Partner
	decoder := json.NewDecoder(r.Body)
	err := decoder.Decode(&p)
	if err != nil {
		ErrorWithJSON(w, "Incorrect body", http.StatusBadRequest)
		return
	}
	p.Updatetime = time.Now().Format("2006-01-02 3:4:5 PM")

	err = c.Insert(p)
	if err != nil {
		if mgo.IsDup(err) {
			log.Println("Failed insert partner: ", err)
			ErrorWithJSON(w, "The telphone number already exists", http.StatusBadRequest)
			return
		}

		ErrorWithJSON(w, "Database error", http.StatusInternalServerError)
		log.Println("Failed insert partner: ", err)
		return
	}

	fmt.Fprintf(w, "{message: success}")
	w.WriteHeader(http.StatusCreated)
}


func GetAllPartnerDataTable(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
	session := mgoSession.Copy()
	defer session.Close()

	c := session.DB("yzschool").C("partner")

	var dbPartners []Partner
	err := c.Find(bson.M{}).All(&dbPartners)

	if err != nil || len(dbPartners) == 0 {
		ErrorWithJSON(w, "Can not found in the database!", http.StatusNotFound)
		return
	}

	var data PartnerDataTable
	data.Data = dbPartners

	respBody, err := json.MarshalIndent(data, "", "  ")
	if err != nil {
		ErrorWithJSON(w, "Database error", http.StatusInternalServerError)
		log.Fatal(err)
		return
	}
	
	ResponseWithJSON(w, respBody, http.StatusOK)
}

func GetAllPartner(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
	session := mgoSession.Copy()
	defer session.Close()

	c := session.DB("yzschool").C("partner")

	var dbPartners []Partner
	err := c.Find(bson.M{}).All(&dbPartners)

	if err != nil || len(dbPartners) == 0 {
		ErrorWithJSON(w, "Can not found in the database!", http.StatusNotFound)
		return
	}

	respBody, err := json.MarshalIndent(dbPartners, "", "  ")
	if err != nil {
		ErrorWithJSON(w, "Database error", http.StatusInternalServerError)
		log.Fatal(err)
		return
	}

	ResponseWithJSON(w, respBody, http.StatusOK)
}
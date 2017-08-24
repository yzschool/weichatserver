package main

import (
	"fmt"
	"log"
	"net/http"

	"gopkg.in/mgo.v2"

	"github.com/julienschmidt/httprouter"
)

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

var mgoSession *mgo.Session

func ensureIndex(s *mgo.Session) {
	session := s.Copy()
	defer session.Close()

	c := session.DB("yzschool").C("class")

	index := mgo.Index{
		Key:        []string{"classid"},
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

func main() {

	session, err := mgo.Dial("localhost")
	if err != nil {
		fmt.Println("connect to mongodb failure", err)
		panic(err)
	}
	defer session.Close()

	session.SetMode(mgo.Monotonic, true)
	mgoSession = session
	ensureIndex(session)

	router := httprouter.New()
	router.GET("/", Index)

	router.GET("/weichat/class", GetClass)
	router.POST("/weichat/class", CreateClass)
	router.GET("/weichat/class/:classid", GetClassByID)
	router.GET("/weichat/classname/:classname", GetClassByName)
	router.PUT("/weichat/class/:classid", UpdateClass)
	router.DELETE("/weichat/class/:classid", DeleteClass)
	router.GET("/weichat/student", GetStudent)
	router.POST("/weichat/student", CreateStudent)
	router.GET("/weichat/studentname/:studentname", HTTPGetClassByName)
	router.GET("/weichat/student/:studentid", HTTPGetStudentById)

	log.Fatal(http.ListenAndServe("0.0.0.0:8080", router))
}

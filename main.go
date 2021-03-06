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

	conn_class := session.DB("yzschool").C("class")

	index := mgo.Index{
		Key:        []string{"classid"},
		Unique:     true,
		DropDups:   true,
		Background: true,
		Sparse:     true,
	}
	err := conn_class.EnsureIndex(index)
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

	/* code for library book management */
	/* List all the admin in the library */
	router.GET("/weichat/book_admin", GetAllAdmin)
	/* List all the book in the library */
	router.POST("/weichat/books", GetAllBook)
	/* add book to the library */
	router.POST("/weichat/book", CreateBook)
	/* Change book status, when some borrow it */
	router.POST("/weichat/book_update", UpdateBook)
	/* Remove book from the library */
	router.POST("/weichat/book_delete", DeleteBook)
	/* Query book by ISBN */
	router.GET("/weichat/book/isbn/:isbn", GetBookByISBN)
	/* Query book by Name */
	router.POST("/weichat/book/name", GetBookByName)
	/* Query book by ID */
	router.POST("/weichat/book/id", GetBookByID)

	/* recruit API */
	/* candidate API */
	router.POST("/recruit/submitapplication", SubmitApplication)
	router.POST("/recruit/getcandidate", GetCandidate)
	router.GET("/recruit/getallcandidate", GetAllCandidate)
	router.GET("/recruit/getallcandidatetable", GetAllCandidateDataTable)
	/* exam API */
	router.POST("/recruit/submitexam", SubmitExam)
	router.POST("/recruit/getexam", GetExam)
	router.GET("/recruit/getallexam", GetAllExam)
	router.GET("/recruit/getallexamtable", GetAllExamDataTable)
	/* partner API */
	router.POST("/recruit/submitpartner", SubmitPartnerApplication)		
	router.GET("/recruit/getallpartner", GetAllPartner)
	router.GET("/recruit/getallpartnertable", GetAllPartnerDataTable)
	/* Course Application API */
	router.POST("/recruit/submitcourseapp", SubmitCourseApplication)		
	router.GET("/recruit/getallcourseapp", GetAllCourseApplication)
	router.GET("/recruit/getallcourseapptable", GetAllCourseApplicationDataTable)

	log.Fatal(http.ListenAndServe("0.0.0.0:8080", router))
}

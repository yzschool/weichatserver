package main

import (
	"fmt"
	"net/http"

	mgo "gopkg.in/mgo.v2"
	"gopkg.in/mgo.v2/bson"

	"github.com/julienschmidt/httprouter"
	uuid "github.com/nu7hatch/gouuid"
)

/*
{
  "studentid":"09882342",
  "name":"张三",
  "phone":"13900000000",
  "studentopenid":"adsf2324sdfa",
  "email":"aaa@163.com",
  "school":"实验小学",
  "grade":"三年级",
  "capibility":"优秀",
  "address":"福田小区",
  "mothername":"张二",
  "motherphone":"13800000000",
  "motheropenid":"adsf2324sdfa",
  "fathername":"张二",
  "fatherphone":"13800000000",
  "fatheropenid":"adsf2324sdfa",
  "updatetime":"2017-07-11"
}
*/

/*Student is the GO struct for student JSON */
type Student struct {
	Studentid     string `json:"studentid"`
	Name          string `json:"name"`
	Phone         string `json:"phone"`
	Studentopenid string `json:"studentopenid"`
	Email         string `json:"email"`
	School        string `json:"school"`
	Grade         string `json:"grade"`
	Capibility    string `json:"capibility"`
	Address       string `json:"address"`
	Mothername    string `json:"mothername"`
	Motherphone   string `json:"motherphone"`
	Motheropenid  string `json:"motheropenid"`
	Fathername    string `json:"fathername"`
	Fatherphone   string `json:"fatherphone"`
	Fatheropenid  string `json:"fatheropenid"`
	Updatetime    string `json:"updatetime"`
}

func GetStudent(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
	fmt.Fprintf(w, "hello, %s!\n", ps.ByName("name"))
}

func GetStudentByName(name string, phone string) string {

	fmt.Println("Get student by Name:" + name + " phone " + phone)
	session := mgoSession.Copy()
	defer session.Close()
	var student Student
	c := session.DB("yzschool").C("student")

	iter := c.Find(bson.M{"name": name}).Iter()
	for iter.Next(&student) {
		fmt.Println("Found student name " + student.Name + " phone " + student.Phone)

		if student.Phone == phone {
			fmt.Println("found student id:", student.Studentid)
			return student.Studentid
		}
	}

	fmt.Println("not found student")
	return ""

}

func CreateStudent(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
	fmt.Fprintf(w, "hello, %s!\n", ps.ByName("name"))
}

func add_student(student Student) string {
	fmt.Println("add student to DB")

	session := mgoSession.Copy()
	defer session.Close()

	c := session.DB("yzschool").C("student")

	/* generate UUID for the student ID */
	u4, err := uuid.NewV4()
	if err != nil {
		fmt.Println("error:", err)
		return ""
	}
	student.Studentid = u4.String()

	err = c.Insert(student)
	if err != nil {
		if mgo.IsDup(err) {
			fmt.Println("duplicated student id")
			return ""
		}

		fmt.Println("Data base err")
		return ""
	}

	return student.Studentid

}

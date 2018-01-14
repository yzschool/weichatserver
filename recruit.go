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

/* example JSON for candidate
{
  "name": "张三",
  "gender": "男",
  "dateofbirth": "9787500128199",
  "schoolname": "深圳大学",
  "department": "外语学院",
  "major": "英语",
  "age": "20",
  "workexperience": "英语教师一年",
  "positionapplied": "英语教师",
  "telephone": "13900000000",
  "mail": "test@123.com",
  "residence": "深圳",
  "selfintroduce": "我是...",
  "jobdecription": "英语教师...",
  "updatetime": "2018"
}
*/

type Candidate struct {
	Name            string `json:"name"`
	Gender          string `json:"gender"`
	Dateofbirth     string `json:"dateofbirth"`
	Schoolname      string `json:"schoolname"`
	Department      string `json:"department"`
	Major           string `json:"major"`
	Age             string `json:"age"`
	Workexperience  string `json:"workexperience"`
	Positionapplied string `json:"positionapplied"`
	Telephone       string `json:"telephone"`
	Mail            string `json:"mail"`
	Residence       string `json:"residence"`
	Selfintroduce   string `json:"selfintroduce"`
	Jobdecription   string `json:"jobdecription"`
	Updatetime      string `json:"updatetime"`
}

/* example JSON for exam
{
  "exam": "入学考试1",
  "tester": "张三",
  "telephone": "13900000000",
  "mail": "test@123.com",
  "result": [
    {
      "question": "你的目标1",
      "answer": "好1"
    },
    {
      "question": "你的目标2",
      "answer": "好2"
    }
  ],
  "score": 90,
  "updatetime": "2018"
}
*/
type Exam struct {
	Exam      string `json:"exam"`
	Tester    string `json:"tester"`
	Telephone string `json:"telephone"`
	Mail      string `json:"mail"`
	Result    []struct {
		Question string `json:"question"`
		Answer   string `json:"answer"`
	} `json:"result"`
	Score      int    `json:"score"`
	Updatetime string `json:"updatetime"`
}

/*SubmitApplication application submit application form */
func SubmitApplication(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
	fmt.Println("SubmitApplication get request")
	session := mgoSession.Copy()
	defer session.Close()

	c := session.DB("yzschool").C("candidate")

	var cdd Candidate
	decoder := json.NewDecoder(r.Body)
	err := decoder.Decode(&cdd)
	if err != nil {
		ErrorWithJSON(w, "Incorrect body", http.StatusBadRequest)
		return
	}
	cdd.Updatetime = time.Now().Format("2006-01-02 3:4:5 PM")

	err = c.Insert(cdd)
	if err != nil {
		ErrorWithJSON(w, "Database error", http.StatusInternalServerError)
		log.Println("Failed insert candidate: ", err)
		return
	}

	fmt.Fprintf(w, "{message: success}")
	w.WriteHeader(http.StatusCreated)
}

func GetCandidate(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
	session := mgoSession.Copy()
	defer session.Close()

	c := session.DB("yzschool").C("candidate")

	var cdd Candidate
	decoder := json.NewDecoder(r.Body)
	err := decoder.Decode(&cdd)
	if err != nil {
		ErrorWithJSON(w, "Incorrect body", http.StatusBadRequest)
		return
	}

	var dbCdds []Candidate
	if cdd.Name != "" {
		err = c.Find(bson.M{"name": bson.M{"$regex": bson.RegEx{cdd.Name, "i"}}}).All(&dbCdds)
	} else if cdd.Mail != "" {
		err = c.Find(bson.M{"mail": bson.M{"$regex": bson.RegEx{cdd.Mail, "i"}}}).All(&dbCdds)
	} else if cdd.Mail != "" {
		err = c.Find(bson.M{"mail": bson.M{"$regex": bson.RegEx{cdd.Mail, "i"}}}).All(&dbCdds)
	}

	if err != nil || len(dbCdds) == 0 {
		ErrorWithJSON(w, "Can not found in the database!", http.StatusNotFound)
		return
	}

	respBody, err := json.MarshalIndent(dbCdds, "", "  ")
	if err != nil {
		ErrorWithJSON(w, "Database error", http.StatusInternalServerError)
		log.Fatal(err)
		return
	}

	ResponseWithJSON(w, respBody, http.StatusOK)
}

func GetAllCandidate(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
	session := mgoSession.Copy()
	defer session.Close()

	c := session.DB("yzschool").C("candidate")

	var dbCdds []Candidate
	err := c.Find(bson.M{}).All(&dbCdds)

	if err != nil || len(dbCdds) == 0 {
		ErrorWithJSON(w, "Can not found in the database!", http.StatusNotFound)
		return
	}

	respBody, err := json.MarshalIndent(dbCdds, "", "  ")
	if err != nil {
		ErrorWithJSON(w, "Database error", http.StatusInternalServerError)
		log.Fatal(err)
		return
	}

	ResponseWithJSON(w, respBody, http.StatusOK)
}

func SubmitExam(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
	fmt.Println("SubmitApplication get request")
	session := mgoSession.Copy()
	defer session.Close()

	c := session.DB("yzschool").C("exam")

	var exam Exam
	decoder := json.NewDecoder(r.Body)
	err := decoder.Decode(&exam)
	if err != nil {
		ErrorWithJSON(w, "Incorrect body", http.StatusBadRequest)
		return
	}
	exam.Updatetime = time.Now().Format("2006-01-02 3:4:5 PM")

	err = c.Insert(exam)
	if err != nil {
		ErrorWithJSON(w, "Database error", http.StatusInternalServerError)
		log.Println("Failed insert candidate: ", err)
		return
	}

	fmt.Fprintf(w, "{message: success}")
	w.WriteHeader(http.StatusCreated)
}

func GetExam(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
	session := mgoSession.Copy()
	defer session.Close()

	c := session.DB("yzschool").C("exam")

	var exam Exam
	decoder := json.NewDecoder(r.Body)
	err := decoder.Decode(&exam)
	if err != nil {
		ErrorWithJSON(w, "Incorrect body", http.StatusBadRequest)
		return
	}
	//fmt.Println(exam)

	var dbExam []Exam
	if exam.Exam != "" {
		err = c.Find(bson.M{"exam": bson.M{"$regex": bson.RegEx{exam.Exam, "i"}}}).All(&dbExam)
	} else if exam.Tester != "" {
		err = c.Find(bson.M{"tester": bson.M{"$regex": bson.RegEx{exam.Tester, "i"}}}).All(&dbExam)
	} else if exam.Telephone != "" {
		err = c.Find(bson.M{"telephone": bson.M{"$regex": bson.RegEx{exam.Telephone, "i"}}}).All(&dbExam)
	} else if exam.Mail != "" {
		err = c.Find(bson.M{"mail": bson.M{"$regex": bson.RegEx{exam.Mail, "i"}}}).All(&dbExam)
	}

	if err != nil || len(dbExam) == 0 {
		ErrorWithJSON(w, "Can not found in the database!", http.StatusNotFound)
		return
	}

	respBody, err := json.MarshalIndent(dbExam, "", "  ")
	if err != nil {
		ErrorWithJSON(w, "Database error", http.StatusInternalServerError)
		log.Fatal(err)
		return
	}

	ResponseWithJSON(w, respBody, http.StatusOK)
}

func GetAllExam(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
	session := mgoSession.Copy()
	defer session.Close()

	c := session.DB("yzschool").C("exam")

	var dbExams []Exam
	err := c.Find(bson.M{}).All(&dbExams)

	if err != nil || len(dbExams) == 0 {
		ErrorWithJSON(w, "Can not found in the database!", http.StatusNotFound)
		return
	}

	respBody, err := json.MarshalIndent(dbExams, "", "  ")
	if err != nil {
		ErrorWithJSON(w, "Database error", http.StatusInternalServerError)
		log.Fatal(err)
		return
	}

	ResponseWithJSON(w, respBody, http.StatusOK)
}

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
  "id":"1202434354545",
  "bookname":"yz english AA",
  "isbn":"13900000000",
  "username":"2017-07-11",
  "updatetime":"2017-07-11"
}
*/

/*Class is GO struct for Class JSON */
type Book struct {
	ID         string `json:"id"`
	Bookname   string `json:"bookname"`
	Isbn       string `json:"isbn"`
	Username   string `json:"username"`
	Updatetime string `json:"updatetime"`
}

/* List all the book in the library */

func GetAllBook(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
	session := mgoSession.Copy()
	defer session.Close()

	c := session.DB("yzschool").C("library")

	var books []Book
	err := c.Find(bson.M{}).All(&books)
	if err != nil {
		ErrorWithJSON(w, "Database error", http.StatusInternalServerError)
		log.Println("Failed get all classes: ", err)
		return
	}

	respBody, err := json.MarshalIndent(books, "", "  ")
	if err != nil {
		log.Fatal(err)
	}

	ResponseWithJSON(w, respBody, http.StatusOK)
}

/* add book to the library */
func CreateBook(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
	session := mgoSession.Copy()
	defer session.Close()

	var book Book
	decoder := json.NewDecoder(r.Body)
	err := decoder.Decode(&book)
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
	book.ID = u4.String()
	book.Updatetime = time.Now().String()

	c := session.DB("yzschool").C("library")

	err = c.Insert(book)
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

func UpdateBook(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
}

/* Remove book from the library */
func DeleteBook(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
}

/* Query book by ISBN */
func GetBookByISBN(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
}

/* Query book by Name */
func GetBookByName(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
}

/* Query book by ID */
func GetBookByID(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
}

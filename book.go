package main

import (
	"encoding/json"
	"fmt"
	"time"
	//"fmt"
	"log"
	"net/http"

	"gopkg.in/mgo.v2"
	"gopkg.in/mgo.v2/bson"

	"github.com/julienschmidt/httprouter"
)

/*
example JSON class
{
  "id":"1202434354545",
  "bookname":"yz english AA",
  "isbn":"9787500128199",
  "borrowname":"张",
  "location":"乔希家",
  "updatetime":"2017-07-11"
}
*/

/*Class is GO struct for Class JSON */
type Book struct {
	ID         string `json:"id"`
	Bookname   string `json:"bookname"`
	Isbn       string `json:"isbn"`
	Borrowname string `json:"borrowname"`
	Location   string `json:"location"`
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

func ensureIndex_lib(c *mgo.Collection) {

	index := mgo.Index{
		Key:        []string{"id"},
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
	book.Updatetime = time.Now().Format("2006-01-02 3:4:5 PM")

	c := session.DB("yzschool").C("library")
	ensureIndex_lib(c)

	err = c.Insert(book)
	if err != nil {
		if mgo.IsDup(err) {
			ErrorWithJSON(w, "book with this id already exists", http.StatusBadRequest)
			return
		}

		ErrorWithJSON(w, "Database error", http.StatusInternalServerError)
		log.Println("Failed insert book: ", err)
		return
	}

	//w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)

}

func UpdateBook(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
	session := mgoSession.Copy()
	defer session.Close()
	c := session.DB("yzschool").C("library")
	ensureIndex_lib(c)

	var book Book
	decoder := json.NewDecoder(r.Body)
	err := decoder.Decode(&book)
	if err != nil {
		ErrorWithJSON(w, "Incorrect body", http.StatusBadRequest)
		return
	}
	//log.Println("Book info from user : " + book.ID + book.Bookname + book.Borrowname)

	query := bson.M{"id": book.ID}

	update := bson.M{"$set": bson.M{"borrowname": book.Borrowname, "updatetime": time.Now().Format("2006-01-02 3:4:5 PM")}}

	err = c.Update(query, update)

	if err != nil {
		ErrorWithJSON(w, "Database error", http.StatusInternalServerError)
		log.Println("Failed insert book: ", err)
		return
	}
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
}

/* Remove book from the library */
func DeleteBook(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
	session := mgoSession.Copy()
	defer session.Close()
	c := session.DB("yzschool").C("library")
	ensureIndex_lib(c)

	var book Book
	decoder := json.NewDecoder(r.Body)
	err := decoder.Decode(&book)
	if err != nil {
		ErrorWithJSON(w, "Incorrect body", http.StatusBadRequest)
		return
	}
	//log.Println("Book info from user : " + book.ID + book.Bookname + book.Borrowname)

	err = c.Remove(bson.M{"id": book.ID})
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

/* Query book by ISBN */
func GetBookByISBN(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {

}

/* Query book by Name */
func GetBookByName(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
	session := mgoSession.Copy()
	defer session.Close()

	c := session.DB("yzschool").C("library")

	bookname := ps.ByName("name")

	var book Book
	err := c.Find(bson.M{"bookname": bson.M{"$regex": bson.RegEx{bookname, "i"}}}).One(&book)
	if err != nil {
		ErrorWithJSON(w, "Can not found the id in the database!", http.StatusNotFound)
		log.Println("Failed get all classes: ", err)
		return
	}

	respBody, err := json.MarshalIndent(book, "", "  ")
	if err != nil {
		ErrorWithJSON(w, "Database error", http.StatusInternalServerError)
		log.Fatal(err)
	}

	ResponseWithJSON(w, respBody, http.StatusOK)
}

/* Query book by ID */
func GetBookByID(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
	session := mgoSession.Copy()
	defer session.Close()

	c := session.DB("yzschool").C("library")

	id := ps.ByName("id")

	var book Book
	err := c.Find(bson.M{"id": id}).One(&book)
	if err != nil {
		ErrorWithJSON(w, "Can not found the id in the database!", http.StatusNotFound)
		//log.Println("Failed get all classes: ", err)
		return
	}

	respBody, err := json.MarshalIndent(book, "", "  ")
	if err != nil {
		ErrorWithJSON(w, "Database error", http.StatusInternalServerError)
		log.Fatal(err)
	}

	ResponseWithJSON(w, respBody, http.StatusOK)
}

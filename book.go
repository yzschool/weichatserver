package main

import (
	"encoding/json"
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
	session := mgoSession.Copy()
	defer session.Close()
	c := session.DB("yzschool").C("library")

	/*
		id := ps.ByName("id")
		fmt.Println("id is: " + id)
	*/
	var book Book
	decoder := json.NewDecoder(r.Body)
	err := decoder.Decode(&book)
	if err != nil {
		ErrorWithJSON(w, "Incorrect body", http.StatusBadRequest)
		return
	}

	var book_database Book
	err = c.Find(bson.M{"id": book.ID}).All(&book_database)

	book_database.Borrowname = book.Borrowname
	c.Update(book_database, &book)

	/* TODO
	   update 	Borrowname	Updatetime
	*/

}

/* Remove book from the library */
func DeleteBook(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
	/*
		    session := mgoSession.Copy()
			defer session.Close()

			id := ps.ByName("isbn")

			c := session.DB("yzschool").C("library")

			err := c.Remove(bson.M{"isbn": id})
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
	*/
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

	var books []Book
	err := c.Find(bson.M{"bookname": bson.M{"$regex": bson.RegEx{bookname, "i"}}}).All(&books)
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

/* Query book by ID */
func GetBookByID(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
	session := mgoSession.Copy()
	defer session.Close()

	c := session.DB("yzschool").C("library")

	id := ps.ByName("id")

	var books []Book
	err := c.Find(bson.M{"id": id}).All(&books)
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

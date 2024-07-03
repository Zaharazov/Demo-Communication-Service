package pages

import (
	"context"
	"html/template"
	"log"
	"main/internal/domain"
	"main/pkg/database/mongodb"
	"net/http"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo/options"
)

func Users_page(w http.ResponseWriter, r *http.Request) { // w - ответ сайту, r - запрос к сайту

	name_for_search := r.FormValue("search_by_name")
	db_users := findUsers(name_for_search)

	tmpl, _ := template.ParseFiles("templates/users_page.html")
	tmpl.Execute(w, db_users) // если ничего не передаем, то пишем nil
}

func findUsers(s_name string) []domain.User {
	result := getAllUsers(s_name)
	return result
}

func getAllUsers(param string) []domain.User {

	// Pass these options to the Find method
	options := options.Find()
	options.SetLimit(100)

	filter := bson.D{}

	if len(param) > 0 {
		filter = bson.D{{"login", param}}
	}

	// Here's an array in which you can store the decoded documents
	var db_users []domain.User // можно добавить *

	// Passing nil as the filter matches all documents in the collection
	cur, err := mongodb.Collection.Find(context.TODO(), filter, options)
	if err != nil {
		log.Fatal(err)
	}

	for cur.Next(context.TODO()) {

		// create a value into which the single document can be decoded
		var elem domain.User
		err := cur.Decode(&elem)
		if err != nil {
			log.Fatal(err)
		}

		db_users = append(db_users, elem) // можно добавить &
	}

	if err := cur.Err(); err != nil {
		log.Fatal(err)
	}

	// Close the cursor once finished
	cur.Close(context.TODO())

	//log.Printf("Found multiple documents (array of pointers): %+v\n", db_users)

	return db_users
}

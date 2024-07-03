package pages

import (
	"context"
	"fmt"
	"html/template"
	"log"
	"main/internal/domain"
	"main/pkg/database/mongodb"
	"net/http"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
)

func Home_page(w http.ResponseWriter, r *http.Request) { // w - ответ сайту, r - запрос к сайту
	tmpl, _ := template.ParseFiles("templates/home_page.html")

	data := domain.User{
		Id:       1,
		Login:    r.FormValue("login"),
		Password: r.FormValue("password"),
		ULAccess: false,
	}

	if len(data.Login) > 0 && len(data.Password) > 0 {

		if err, _ := getUser(data.Login, data.Password); err != nil {
			createUser(data)
		} else {
			log.Println("User is already exist!")
		}

	}

	db_users := getAllUsers("")
	user_login := r.FormValue("user_login")
	user_pass := r.FormValue("user_pass")
	db_user := domain.User{}

	for _, user := range db_users {
		if user_login == user.Login && user_pass == user.Password {

			updateResult := changeUserULAccess(user.Login, user.Password, true)

			log.Printf("Matched %v documents and updated %v documents.\n", updateResult.MatchedCount, updateResult.ModifiedCount)
			_, db_user = getUser(user_login, user_pass)
		}
	}

	tmpl.Execute(w, db_user) // если ничего не передаем, то пишем nil
}

func changeUserULAccess(l string, p string, value bool) *mongo.UpdateResult {
	filter := bson.D{{"login", l}, {"password", p}}
	update := bson.D{
		{"$set", bson.D{
			{"ulaccess", value},
		}},
	}

	updateResult, err := mongodb.Collection.UpdateOne(context.TODO(), filter, update)
	if err != nil {
		log.Fatal(err)
	}

	return updateResult
}

func createUser(data domain.User) {
	insertResult, err := mongodb.Collection.InsertOne(context.TODO(), data)
	if err != nil {
		log.Fatal(err)
	}
	log.Println("Inserted a single document: ", insertResult.InsertedID)
}

func getUser(l string, p string) (error, domain.User) {
	// create a value into which the result can be decoded
	var result domain.User

	filter := bson.D{{"login", l}, {"password", p}}
	err := mongodb.Collection.FindOne(context.TODO(), filter).Decode(&result)
	if err != nil {
		return err, result
	}

	fmt.Printf("Found a single document: %+v\n", result)

	return nil, result
}

package main

import (
	"context"
	"fmt"
	"html/template"
	"log"
	"net/http"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	// go get go.mongodb.org/mongo-driver/mongo
)

type User struct {
	Id       int
	Login    string
	Password string
	ULAccess bool
}

var users = []User{}

var client, err = mongo.NewClient(options.Client().ApplyURI("mongodb://127.0.0.1:27017"))
var collection = client.Database("test2").Collection("users")

func home_page(w http.ResponseWriter, r *http.Request) { // w - ответ сайту, r - запрос к сайту
	tmpl, _ := template.ParseFiles("templates/home_page.html")

	data := User{
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
	db_user := User{}

	for _, user := range db_users {
		if user_login == user.Login && user_pass == user.Password {

			updateResult := changeUserULAccess(user.Login, user.Password, true)

			log.Printf("Matched %v documents and updated %v documents.\n", updateResult.MatchedCount, updateResult.ModifiedCount)
			_, db_user = getUser(user_login, user_pass)
		}
	}

	tmpl.Execute(w, db_user) // если ничего не передаем, то пишем nil
}

func users_page(w http.ResponseWriter, r *http.Request) { // w - ответ сайту, r - запрос к сайту

	name_for_search := r.FormValue("search_by_name")
	db_users := findUsers(name_for_search)

	tmpl, _ := template.ParseFiles("templates/users_page.html")
	tmpl.Execute(w, db_users) // если ничего не передаем, то пишем nil
}

func handleRequest() {
	http.HandleFunc("/", home_page)       // отслеживаем переход по URL (/ - переход на главную страницу)
	http.HandleFunc("/users", users_page) // ВАЖНО - в конце URL дописываем /, чтобы он корректно обрабатывался

	http.ListenAndServe(":8080", nil) // запускаем локальный сервер на порту 8080 (параметры: порт и настройки запуска)
}

func getAllUsers(param string) []User {

	// Pass these options to the Find method
	options := options.Find()
	options.SetLimit(100)

	filter := bson.D{}

	if len(param) > 0 {
		filter = bson.D{{"login", param}}
	}

	// Here's an array in which you can store the decoded documents
	var db_users []User // можно добавить *

	// Passing nil as the filter matches all documents in the collection
	cur, err := collection.Find(context.TODO(), filter, options)
	if err != nil {
		log.Fatal(err)
	}

	for cur.Next(context.TODO()) {

		// create a value into which the single document can be decoded
		var elem User
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

func findUsers(s_name string) []User {
	result := getAllUsers(s_name)
	return result
}

func getUser(l string, p string) (error, User) {
	// create a value into which the result can be decoded
	var result User

	filter := bson.D{{"login", l}, {"password", p}}
	err = collection.FindOne(context.TODO(), filter).Decode(&result)
	if err != nil {
		return err, result
	}

	fmt.Printf("Found a single document: %+v\n", result)

	return nil, result
}

func createUser(data User) {
	insertResult, err := collection.InsertOne(context.TODO(), data)
	if err != nil {
		log.Fatal(err)
	}
	log.Println("Inserted a single document: ", insertResult.InsertedID)
}

func changeUserULAccess(l string, p string, value bool) *mongo.UpdateResult {
	filter := bson.D{{"login", l}, {"password", p}}
	update := bson.D{
		{"$set", bson.D{
			{"ulaccess", value},
		}},
	}

	updateResult, err := collection.UpdateOne(context.TODO(), filter, update)
	if err != nil {
		log.Fatal(err)
	}

	return updateResult
}

func connectToMongo() {
	// Create connect
	err = client.Connect(context.TODO())
	if err != nil {
		log.Fatal(err)
	}

	// Check the connection
	err = client.Ping(context.TODO(), nil)
	if err != nil {
		log.Fatal(err)
	}

	fmt.Println("Connected to MongoDB!")
}

func main() {
	connectToMongo()

	handleRequest()
}

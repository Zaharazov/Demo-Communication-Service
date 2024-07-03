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

func Jobs_page(w http.ResponseWriter, r *http.Request) {
	tmpl, _ := template.ParseFiles("templates/jobs_page.html")

	job := domain.Job{
		Title:       r.FormValue("title"),
		Description: r.FormValue("description"),
		Author:      r.FormValue("author"),
		Money:       r.FormValue("money"),
	}

	if len(job.Title) > 0 && len(job.Author) > 0 {
		createJob(job)
	}

	db_jobs := getAllJobs()
	tmpl.Execute(w, db_jobs)
}

func createJob(job domain.Job) {
	insertResult, err := mongodb.Collection2.InsertOne(context.TODO(), job)
	if err != nil {
		log.Fatal(err)
	}
	log.Println("Inserted a single document: ", insertResult.InsertedID)
}

func getAllJobs() []domain.Job {

	// Pass these options to the Find method
	options := options.Find()
	options.SetLimit(100)

	filter := bson.D{}

	// Here's an array in which you can store the decoded documents
	var db_jobs []domain.Job // можно добавить *

	// Passing nil as the filter matches all documents in the collection
	cur, err := mongodb.Collection2.Find(context.TODO(), filter, options)
	if err != nil {
		log.Fatal(err)
	}

	for cur.Next(context.TODO()) {

		// create a value into which the single document can be decoded
		var elem domain.Job
		err := cur.Decode(&elem)
		if err != nil {
			log.Fatal(err)
		}

		db_jobs = append(db_jobs, elem) // можно добавить &
	}

	if err := cur.Err(); err != nil {
		log.Fatal(err)
	}

	// Close the cursor once finished
	cur.Close(context.TODO())

	//log.Printf("Found multiple documents (array of pointers): %+v\n", db_users)

	return db_jobs
}

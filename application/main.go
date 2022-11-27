package main

import (
	"context"
	"fmt"
	"io"
	"net/http"
	"os"
	"strconv"
	"time"

	"github.com/google/uuid"

	"github.com/gin-gonic/gin"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

// Path to certificate
var mongoCertificate = os.Getenv("MONGODB_CERTIFICATE")

// MongoDB connection URI
var uri = "mongodb+srv://hosteldb.e3ayhyn.mongodb.net/?" +
	"retryWrites=true&w=majority" +
	"&authSource=%24external&authMechanism=MONGODB-X509&tlsCertificateKeyFile=" + mongoCertificate

func main() {
	r := gin.Default() // initialize a gin router with default middlewares

	// Create a new client and connect to the server
	mongoClient, err := mongo.Connect(context.TODO(), options.Client().ApplyURI(uri))
	if err != nil {
		panic(err)
	}
	defer func() {
		if err = mongoClient.Disconnect(context.TODO()); err != nil {
			panic(err)
		}
	}()

	// Group using gin.BasicAuth() middleware
	authorized := r.Group("/user", gin.BasicAuth(gin.Accounts{
		"foo":    "bar",
		"austin": "1234",
		"lena":   "hello2",
		"manu":   "4321",
	}))

	// /user/profile endpoint
	authorized.GET("/profile", func(ctx *gin.Context) {
		// get user, it was set by the BasicAuth middleware
		user := ctx.MustGet(gin.AuthUserKey).(string)
		var result bson.M

		collection := mongoClient.Database("University").Collection("Students")
		collection.FindOne(ctx, bson.D{{Key: "name", Value: user}}).Decode(&result)
		ctx.String(http.StatusOK, fmt.Sprint(result["password"]))
	})

	// /user/apply endpoint
	authorized.POST("/apply", func(ctx *gin.Context) {
		// get user, it was set by the BasicAuth middleware
		user := ctx.MustGet(gin.AuthUserKey).(string)

		// Multipart form
		form, err := ctx.MultipartForm()
		if err != nil {
			ctx.String(http.StatusBadRequest, "get form err: %s", err.Error())
			return
		}
		files := form.File["upload[]"]
		var documents [][]byte
		for _, file := range files {
			reader, err := file.Open()
			if err != nil {
				panic(err)
			}
			document, err := io.ReadAll(reader)
			if err != nil {
				panic(err)
			}
			documents = append(documents, document)
		}
		collection := mongoClient.Database("University").Collection("Students")

		formId := ctx.PostForm("formId")
		if formId == "" {
			formId = uuid.New().String()
		}

		filter := bson.M{
			"forms": bson.M{
				"$exists": formId,
			},
			"name": user,
		}
		update := bson.M{
			"$set": bson.M{
				fmt.Sprintf("forms.%s", formId): bson.M{
					"quotas":       form.Value["quotas[]"],
					"files":        documents,
					"status":       "SUBMITTED",
					"submission":   time.Now().UTC().Format(time.RFC3339),
					"verification": nil,
					"payment":      nil,
				},
			},
		}

		_, err = collection.UpdateOne(context.TODO(), filter, update)
		if err != nil {
			panic(err)
		}

		ctx.String(http.StatusOK, "Application ID=%s.", formId)
	})

	r.GET("/download", func(ctx *gin.Context) {
		collection := mongoClient.Database("University").Collection("Students")
		formId := ctx.Query("formId")

		filter := bson.D{{Key: "formId", Value: formId}}

		var form bson.M
		err := collection.FindOne(ctx, filter).Decode(&form)
		if err != nil {
			panic(err)
		}

		for i, doc := range form["files"].(bson.A) {
			os.WriteFile("downloaded"+strconv.Itoa(i), doc.(primitive.Binary).Data, 0777)
		}
	})

	r.Run(":8080") // listen and serve on 0.0.0.0:8080 (for windows "localhost:8080")
}

package main

import (
	"context"
	"fmt"
	"net/http"
	"os"
	"path/filepath"

	"github.com/gin-gonic/gin"
	"go.mongodb.org/mongo-driver/bson"
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

		for _, file := range files {
			filename := filepath.Base(file.Filename)
			if err := ctx.SaveUploadedFile(file, "recieved_"+filename); err != nil {
				ctx.String(http.StatusBadRequest, "upload file err: %s", err.Error())
				return
			}
		}

		ctx.String(http.StatusOK, "Uploaded successfully %d files with fields name=%s and email=%s.", len(files), user, form.Value["quotas[]"])
	})

	r.Run(":8080") // listen and serve on 0.0.0.0:8080 (for windows "localhost:8080")
}

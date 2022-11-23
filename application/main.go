package main

import (
	"net/http"
	"path/filepath"

	"github.com/gin-gonic/gin"
)

func main() {
	r := gin.Default() // initialize a gin router with default middlewares

	// Group using gin.BasicAuth() middleware
	authorized := r.Group("/user", gin.BasicAuth(gin.Accounts{
		"foo":    "bar",
		"austin": "1234",
		"lena":   "hello2",
		"manu":   "4321",
	}))

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

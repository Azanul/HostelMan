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
	// hit "localhost:8080/user/user
	authorized.POST("/apply", func(c *gin.Context) {
		// get user, it was set by the BasicAuth middleware
		// user := c.MustGet(gin.AuthUserKey).(string)

		name := c.PostForm("name")

		// Multipart form
		form, err := c.MultipartForm()
		if err != nil {
			c.String(http.StatusBadRequest, "get form err: %s", err.Error())
			return
		}
		files := form.File["upload[]"]

		for _, file := range files {
			filename := filepath.Base(file.Filename)
			if err := c.SaveUploadedFile(file, filename); err != nil {
				c.String(http.StatusBadRequest, "upload file err: %s", err.Error())
				return
			}
		}

		c.String(http.StatusOK, "Uploaded successfully %d files with fields name=%s and email=%s.", len(files), name, form.Value["quotas[]"])
	})

	r.Run(":8080") // listen and serve on 0.0.0.0:8080 (for windows "localhost:8080")
}

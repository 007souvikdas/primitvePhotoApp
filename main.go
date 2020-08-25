package main

import (
	"fmt"
	"io"
	"log"
	"os"
	"os/exec"
	"strconv"
	"strings"

	"github.com/gin-gonic/gin"
)

func main() {
	r := gin.Default()
	r.LoadHTMLFiles("index.html", "error.html")
	r.GET("/page", func(c *gin.Context) {
		c.HTML(200, "index.html", nil)
	})
	r.POST("/process", func(c *gin.Context) {
		fmt.Println(c.Request)
		file, header, err := c.Request.FormFile("image")
		n := c.PostForm("number")
		m := c.PostForm("shape")
		n = strings.Trim(n, " ")
		if n == "" {
			n = "50"
		}
		int2, _ := strconv.ParseInt(n, 2, 32)
		if int2 > 500 {
			c.HTML(200, "error.html", gin.H{
				"message": "n cannot be greater than 500",
			})
			return
		}

		if err != nil {
			c.HTML(200, "error.html", gin.H{
				"message": "Unable to read file",
			})
			return
		}
		filename := header.Filename
		fmt.Println(header.Filename)
		out, err := os.Create("./images/" + filename)
		if err != nil {
			log.Fatal(err)
			c.HTML(200, "error.html", gin.H{
				"message": "Unable to save the file in server",
			})
			return
			//return error in this case
		}
		defer out.Close()
		_, err = io.Copy(out, file)
		if err != nil {
			log.Fatal(err)
			c.HTML(200, "error.html", gin.H{
				"message": "Unable to copy the contents of the file",
			})
			return
			//return error in this case
		}

		cmd := exec.Command("primitive", "-i", "./images/"+filename, "-o", "./output/"+filename, "-n", n, "-m", m)

		_, err2 := cmd.CombinedOutput()
		if err2 != nil {
			fmt.Println(err2)
			c.HTML(200, "error.html", gin.H{
				"message": "Failed to perform primitive ops on uploaded image ",
			})
			return
		}

		c.Header("Content-Description", "File Transfer")
		c.Header("Content-Transfer-Encoding", "binary")
		c.Header("Content-Disposition", "attachment; filename="+filename)
		c.Header("Content-Type", "application/octet-stream")
		c.File("./output/" + filename)

		//remove the files too
		os.Remove("./output/" + filename)
		os.Remove("./images/" + filename)
	})
	r.Run(":8112")
}

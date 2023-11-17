package main

import (
	cmd "gin/cmd"

	"github.com/gin-gonic/gin"
)


func main() {
	r := gin.Default()
	r.POST("/readFile", cmd.FileData)
	r.Run(":3000")
}


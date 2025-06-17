package main

import (
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
)

type Sheet struct {
	Id          int    `json:"id"`
	Title       string `json:"title"`
	Description string `json:"description"`
}

var sheets []Sheet
var nextId = 1

func main() {
	router := gin.Default()

	router.GET("/", func(c *gin.Context) {
		c.JSON(http.StatusOK, "Hello from route")

	})

	router.GET("/sheets", func(c *gin.Context) {
		c.JSON(http.StatusOK, sheets)

	})
	router.POST("/sheets", func(c *gin.Context) {
		var sheet Sheet
		if err := c.ShouldBindJSON(&sheet); err != nil {
			c.JSON(http.StatusBadRequest, "Invalid Json")
			return
		}
		sheet.Id = nextId
		nextId++

		sheets = append(sheets, sheet)
		c.JSON(http.StatusCreated, "A new sheet created with title"+sheet.Title)

	})
	router.GET("/sheets/:id", func(c *gin.Context) {
		sheetId := c.Param("id")
		id, err := strconv.Atoi(sheetId)
		if err != nil {
			c.JSON(http.StatusBadRequest, "Invalid Id")
			return
		}
		for _, sheet := range sheets {
			if sheet.Id == id {
				c.JSON(http.StatusOK, sheet)
				return

			}
		}
		c.JSON(http.StatusNotFound, "Sheet not found")

	})

	router.Run(":3000")
}

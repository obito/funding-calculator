package main

import (
	"fmt"
	"io/ioutil"
	"net/http"
	"time"

	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
	"github.com/gin-gonic/gin/binding"
)

type FundingRates struct {
	Future    string    `form:"future" binding:"required"`
	StartTime time.Time `form:"start_time" binding:"required" time_format:"2006-01-02"`
	EndTime   time.Time `form:"end_time" binding:"required" time_format:"2006-01-02"`
}

type FTXResponse struct {
	Success bool `json:"success"`
	Result  []struct {
		Future string    `json:"future"`
		Rate   float64   `json:"rate"`
		Time   time.Time `json:"time"`
	} `json:"result"`
}

func setupRouter() *gin.Engine {
	// Disable Console Color
	// gin.DisableConsoleColor()
	r := gin.Default()

	config := cors.DefaultConfig()
	config.AllowAllOrigins = true
	config.AddAllowHeaders("OPTIONS")

	r.Use(cors.New(config))

	// Ping test
	r.GET("/ping", func(c *gin.Context) {
		c.String(http.StatusOK, "pong")
	})

	r.GET("/fundings", func(c *gin.Context) {

		var fr FundingRates
		if err := c.ShouldBindWith(&fr, binding.Query); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}

		url := fmt.Sprintf("https://ftx.com/api/funding_rates?future=%s&start_time=%d&end_time=%d", fr.Future, fr.StartTime.Unix(), fr.EndTime.Unix())

		client := &http.Client{}
		req, err := http.NewRequest("GET", url, nil)

		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{
				"message": "Something bad happened!",
				"error":   err.Error(),
			})
			return
		}
		res, err := client.Do(req)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{
				"message": "Something bad happened!",
				"error":   err.Error(),
			})
			return
		}
		defer res.Body.Close()

		if res.StatusCode != 200 {
			c.JSON(http.StatusInternalServerError, gin.H{
				"message": "Something bad happened!",
				"error":   "Status code from FTX is not 200",
			})
			return
		}

		body, err := ioutil.ReadAll(res.Body)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{
				"message": "Something bad happened!",
				"error":   err.Error(),
			})
			return
		}

		c.Data(http.StatusOK, "application/json", body)
	})

	return r
}

func main() {
	r := setupRouter()
	// Listen and Server in 0.0.0.0:8080
	r.Run(":10000")
}

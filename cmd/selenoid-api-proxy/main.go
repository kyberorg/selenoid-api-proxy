package main

import (
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/kyberorg/selenoid-api-proxy/cmd/selenoid-api-proxy/config"
	"io/ioutil"
	"log"
	"net/http"
	"strconv"
	"strings"
)

var httpClient = &http.Client{}

func main() {
	r := gin.Default()
	r.GET("/ping", ping)
	r.DELETE("/videos/:filename", deleteVideo)

	selenoidApiAvailable, selenoidStatusMessage := isSelenoidApiAvailable()
	if !selenoidApiAvailable {
		log.Fatal(fmt.Printf("Selenoid API is not available. Error is %s", selenoidStatusMessage))
	}

	portString := strconv.Itoa(int(config.GetAppConfig().Port))
	err := r.Run(":" + portString)
	if err != nil {
		log.Fatal(fmt.Printf("failed to start server. %s\n", err.Error()))
		return
	}
}

func isSelenoidApiAvailable() (bool, string) {
	resp, err := http.Get((*config.GetAppConfig().SelenoidApiUrl).String())
	if err != nil {
		return false, err.Error()
	}
	if resp.StatusCode != 200 {
		return false, "Selenoid API replied with status " + resp.Status
	}

	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return false, err.Error()
	}
	bodyString := string(body)
	if strings.Contains(bodyString, "Selenoid") {
		return true, "OK"
	} else {
		return false, "Not a Selenoid instance"
	}
}

func ping(c *gin.Context) {
	c.JSON(200, gin.H{
		"message": "pong",
	})
}

func deleteVideo(c *gin.Context) {
	name := c.Param("filename")
	if len(name) == 0 || len(strings.TrimSpace(name)) == 0 {
		c.JSON(400, gin.H{
			"message": "filename must not be empty",
		})
	}
	if strings.Contains(name, ".mp4") {
		name = strings.ReplaceAll(name, ".mp4", "")
	}

	if strings.Contains(name, "mp4") {
		name = strings.ReplaceAll(name, "mp4", "")
	}

	token := c.GetHeader("X-Token")
	tokenValid := isTokenValid(token)
	if tokenValid {
		resp, err := doDeleteVideo(name)
		if err != nil {
			log.Printf("failed to make Delete request, error: %s", err.Error())
			c.JSON(500, gin.H{
				"message": "Error making API request",
			})
		}
		if resp.StatusCode == 200 {
			c.JSON(200, gin.H{
				"message": "Video deleted",
			})
		} else if resp.StatusCode == 404 {
			c.JSON(404, gin.H{
				"message": "No such video " + name,
			})
		} else {
			log.Printf("Unknown status in response %s", resp.Status)
			c.JSON(500, gin.H{
				"message": "Unexpected reply from Selenoid API",
			})
		}
	} else {
		c.JSON(401, gin.H{
			"message": "Token is invalid",
		})
	}
}

func isTokenValid(token string) bool {
	return token == config.GetAppConfig().Token
}

func doDeleteVideo(filename string) (*http.Response, error) {
	deleteEndpoint :=
		(*config.GetAppConfig().SelenoidApiUrl).String() + "/video/" + filename + ".mp4"
	//Create Request
	req, err := http.NewRequest(http.MethodDelete, deleteEndpoint, nil)
	if err != nil {
		return nil, err
	}
	//Send Request
	resp, err := httpClient.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()
	return resp, nil
}

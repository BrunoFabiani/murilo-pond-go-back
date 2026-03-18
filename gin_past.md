package main

import (
    "net/http"

    "github.com/gin-gonic/gin"
)

// album represents data about a record album.
type People struct {
    ID     string  `json:"id"`
    Name  string  `json:"name"`
}

// albums slice to seed record album data.
var group = []People{

}

// postAlbums adds an album from JSON received in the request body.
func postGroup(c *gin.Context) {
    var newGroup People

    // Call BindJSON to bind the received JSON
    if err := c.BindJSON(&newGroup); err != nil {
        return
    }

    // Add the new album to the slice.
    group = append(group, newGroup)
    c.IndentedJSON(http.StatusCreated, newGroup)
}

func main() {
    router := gin.Default()
    router.POST("/group", postGroup)
    
    router.Run(":8080")
}
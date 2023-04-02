package main

import (
	"log"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/gocql/gocql"
	"github.com/zacker/cassandra/taskapp/db"
)

func main() {
	r := gin.Default()
	r.GET("/", func(c *gin.Context) {
		c.JSON(200, gin.H{
			"message": "HELLO WORLD",
		})
	})

	cluster := gocql.NewCluster("localhost")
	cluster.Keyspace = "app"
	cluster.Consistency = gocql.Quorum

	session, err := cluster.CreateSession()
	if err != nil {
		log.Fatalf("error while connecting db: %v", err)
	}

	userRepository := db.NewUserRepository(session)
	r.GET("/users", func(c *gin.Context) {
		users, err := userRepository.FetchUsers(100)
		if err != nil {
			log.Printf("error: %v\n", err)
			c.Status(http.StatusInternalServerError)
			return
		}
		c.JSON(http.StatusOK, users)
	})

	r.Run()
}

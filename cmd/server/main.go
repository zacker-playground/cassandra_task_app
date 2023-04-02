package main

import (
	"log"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/gocql/gocql"
	"github.com/zacker/cassandra/taskapp/db"
)

type Server struct {
	userRepository db.UserRepository
	taskRepository db.TaskRepository
}

func main() {
	r := gin.Default()

	cluster := gocql.NewCluster("localhost")
	cluster.Keyspace = "app"
	cluster.Consistency = gocql.Quorum
	cluster.Timeout = time.Second * 10

	session, err := cluster.CreateSession()
	if err != nil {
		log.Fatalf("error while connecting db: %v", err)
	}

	server := Server{
		userRepository: db.NewUserRepository(session),
		taskRepository: db.NewTaskRepository(session),
	}

	r.GET("/", server.Index)
	r.GET("/users", server.UserIndex)
	r.GET("/users/:id", server.UserTasks)

	r.Run()
}

func (r Server) Index(c *gin.Context) {
	c.JSON(200, gin.H{
		"message": "HELLO WORLD",
	})
}

func (r Server) UserIndex(c *gin.Context) {
	users, err := r.userRepository.FetchUsers(100)
	if err != nil {
		log.Printf("error: %v\n", err)
		c.Status(http.StatusInternalServerError)
		return
	}
	c.JSON(http.StatusOK, users)
}

func (r Server) UserTasks(c *gin.Context) {
	tasks, err := r.taskRepository.FindTasks(c.Param("id"))
	if err != nil {
		log.Printf("error: %v\n", err)
		c.Status(http.StatusInternalServerError)
		return
	}
	c.JSON(http.StatusOK, tasks)
}

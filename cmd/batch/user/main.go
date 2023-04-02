package main

import (
	"log"
	"sync"
	"time"

	"github.com/gocql/gocql"
	"github.com/goombaio/namegenerator"
)

func main() {
	cluster := gocql.NewCluster("localhost")
	cluster.Keyspace = "app"
	cluster.Consistency = gocql.Quorum

	session, err := cluster.CreateSession()
	if err != nil {
		log.Fatalf("error while connecting db: %v", err)
	}

	// 10並列で書き込みを行う
	var wg sync.WaitGroup
	sem := make(chan struct{}, 20)

	seed := time.Now().UTC().UnixNano()
	nameGenerator := namegenerator.NewNameGenerator(seed)

	start := time.Now()
	total := 1000000
	for i := 0; i < total; i++ {
		sem <- struct{}{}
		wg.Add(1)

		go func(n int) {
			defer wg.Done()
			defer func() { <-sem }()

			name := nameGenerator.Generate()
			age := 10 + n%60
			CreateUserData(name, age, session)
		}(i)
	}

	wg.Wait()

	elapsed := time.Since(start)
	log.Printf("Batch Create User took %s", elapsed)
}

func CreateUserData(name string, age int, session *gocql.Session) {
	query := session.Query(
		"INSERT INTO app.users (id, name, age, created_at, updated_at) VALUES (uuid(), ?, ?, toTimeStamp(now()), toTimeStamp(now()))",
		name,
		age,
	)
	if err := query.Exec(); err != nil {
		log.Fatalf("error while inserting data to db: %v", err)
	}
}

package main

import (
	"fmt"
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
	sem := make(chan struct{}, 3)

	seed := time.Now().UTC().UnixNano()
	nameGenerator := namegenerator.NewNameGenerator(seed)

	start := time.Now()

	// 100万ユーザーを生成
	total := 100000
	for i := 0; i < total; i++ {
		sem <- struct{}{}
		wg.Add(1)

		go func(n int, name string) {
			defer wg.Done()
			defer func() { <-sem }()

			batch := session.NewBatch(gocql.LoggedBatch)
			userId := gocql.MustRandomUUID()
			age := 10 + n%60

			CreateUserData(userId, name, age, batch)

			// 1ユーザーごと100個のタスクを持つ
			for j := 0; j < 100; j++ {
				CreateTask(userId, fmt.Sprintf("The task of %s [%d / 100]", name, j), batch)
			}

			if err := session.ExecuteBatch(batch); err != nil {
				log.Printf("error while inserting data: %v\n", err)
			}
		}(i, nameGenerator.Generate())
	}

	wg.Wait()

	elapsed := time.Since(start)
	log.Printf("Batch Create User took %s", elapsed)
}

func CreateUserData(userId gocql.UUID, name string, age int, session *gocql.Batch) {
	session.Query(
		"INSERT INTO app.users (id, name, age, created_at, updated_at) VALUES (?, ?, ?, toTimeStamp(now()), toTimeStamp(now()))",
		userId,
		name,
		age,
	)
}

func CreateTask(userId gocql.UUID, title string, session *gocql.Batch) {
	session.Query(
		"INSERT INTO app.tasks (id, user_id, title, checked, created_at, updated_at) VALUES (uuid(), ?, ?, FALSE, toTimestamp(now()), toTimestamp(now()))",
		userId,
		title,
	)
}

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
	sem := make(chan struct{}, 20)

	seed := time.Now().UTC().UnixNano()
	nameGenerator := namegenerator.NewNameGenerator(seed)

	start := time.Now()

	// 100万ユーザーを生成
	total := 100000
	for i := 0; i < total; i++ {
		sem <- struct{}{}
		wg.Add(1)

		go func(n int) {
			defer wg.Done()
			defer func() { <-sem }()

			name := nameGenerator.Generate()
			age := 10 + n%60
			userId := CreateUserData(name, age, session)

			// 1ユーザーごと100個のタスクを持つ
			for j := 0; j < 100; j++ {
				CreateTask(fmt.Sprintf("The task of %s [%d / 100]", name, j), userId, session)
			}
		}(i)
	}

	wg.Wait()

	elapsed := time.Since(start)
	log.Printf("Batch Create User took %s", elapsed)
}

func CreateUserData(name string, age int, session *gocql.Session) gocql.UUID {
	userId := gocql.MustRandomUUID()
	query := session.Query(
		"INSERT INTO app.users (id, name, age, created_at, updated_at) VALUES (?, ?, ?, toTimeStamp(now()), toTimeStamp(now()))",
		userId,
		name,
		age,
	)
	if err := query.Exec(); err != nil {
		log.Fatalf("error while inserting user data: %v", err)
	}

	return userId
}

func CreateTask(title string, userId gocql.UUID, session *gocql.Session) {
	query := session.Query(
		"INSERT INTO app.tasks (id, user_id, title, checked, created_at, updated_at) VALUES (uuid(), ?, ?, FALSE, toTimestamp(now()), toTimestamp(now()))",
		userId,
		title,
	)
	if err := query.Exec(); err != nil {
		log.Fatalf("error while inserting task data: %v", err)
	}
}

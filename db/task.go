package db

import (
	"time"

	"github.com/gocql/gocql"
)

type TaskEntity struct {
	Id        string
	UserId    string
	Title     string
	Checked   bool
	CreatedAt time.Time
	UpdatedAt time.Time
}

type TaskRepository struct {
	session *gocql.Session
}

func NewTaskRepository(session *gocql.Session) TaskRepository {
	return TaskRepository{
		session: session,
	}
}

func (r TaskRepository) FindTasks(userId string) ([]TaskEntity, error) {
	query := r.session.Query(
		"SELECT id, user_id, title, checked, created_at, updated_at FROM app.tasks WHERE user_id = ?",
		userId,
	)
	iter := query.Iter()

	var tasks []TaskEntity
	var task TaskEntity
	for iter.Scan(
		&task.Id,
		&task.UserId,
		&task.Title,
		&task.Checked,
		&task.CreatedAt,
		&task.UpdatedAt,
	) {
		tasks = append(tasks, task)
	}

	if err := iter.Close(); err != nil {
		return nil, err
	}

	return tasks, nil
}

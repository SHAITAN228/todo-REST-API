package database

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"os"
	"todo-app/models"

	_ "github.com/lib/pq" // Драйвер для PostgreSQL
)

var DB *sql.DB

func InitDB() error {
	var err error

	configData, err := os.ReadFile("config.json")
	if err != nil {
		fmt.Printf("Не удалось прочитать файл конфигурации: %v\n", err)
		return err
	}

	var configMap map[string]interface{}
	err = json.Unmarshal(configData, &configMap)
	if err != nil {
		return fmt.Errorf("не удалось распарсить JSON: %v", err)
	}

	dataSourceName := fmt.Sprintf(
		"user=%s password=%s host=%s port=%d dbname=%s sslmode=%s",
		configMap["user"], configMap["password"], configMap["host"],
		int(configMap["port"].(float64)), configMap["dbname"], configMap["sslmode"],
	)

	DB, err = sql.Open("postgres", dataSourceName)

	if err != nil {
		return err
	}

	_, err = DB.Exec(`
	CREATE TABLE IF NOT EXISTS tasks (
	 id SERIAL PRIMARY KEY,
	 title TEXT NOT NULL,
	 content TEXT,
	 completed BOOLEAN DEFAULT FALSE
	);
   `)
	if err != nil {
		return err
	}

	fmt.Print("Успешно подключено")

	return nil
}

func InsertTask(title, content string, completed bool) (int64, error) {
	var id int64
	query := `
	 INSERT INTO tasks (title, content, completed)
	 VALUES ($1, $2, $3)
	 RETURNING id
	`
	err := DB.QueryRow(query, title, content, completed).Scan(&id)
	if err != nil {
		return 0, err
	}
	return id, nil
}

func UpdateTask(id int64, title, content string, completed bool) (int64, error) {
	query := `
        UPDATE tasks
        SET title = $1, content = $2, completed = $3
        WHERE id = $4
        RETURNING id
    `

	var affectedID int64
	err := DB.QueryRow(query, title, content, completed, id).Scan(&affectedID)
	if err != nil {
		if err == sql.ErrNoRows {
			return 0, fmt.Errorf("task with id %d not found", id)
		}
		return 0, err
	}

	return affectedID, nil
}

func GetAllTasks() ([]models.Task, error) {

	query := `
        SELECT id, title, content, completed
        FROM tasks
    `

	rows, err := DB.Query(query)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var tasks []models.Task
	for rows.Next() {
		var task models.Task
		err := rows.Scan(&task.ID, &task.Title, &task.Content, &task.Completed)
		if err != nil {
			return nil, err
		}
		tasks = append(tasks, task)
	}

	return tasks, nil
}

func GetTaskByID(id int64) (*models.Task, error) {

	query := `
        SELECT id, title, content, completed
        FROM tasks
        WHERE id = $1
    `

	var task models.Task

	err := DB.QueryRow(query, id).Scan(&task.ID, &task.Title, &task.Content, &task.Completed)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, fmt.Errorf("task with id %d not found", id)
		}
		return nil, err
	}

	return &task, nil
}

func DeleteTask(id int64) error {
	query := `
	 DELETE FROM tasks
	 WHERE id = $1
	`
	_, err := DB.Exec(query, id)
	return err
}

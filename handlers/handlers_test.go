package handlers_test

import (
	"bytes"
	"encoding/json"
	"github.com/gorilla/mux"
	"github.com/spf13/viper"
	"log"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/mussyaroslav/otus_final_project/database"
	"github.com/mussyaroslav/otus_final_project/handlers"
	"github.com/mussyaroslav/otus_final_project/models"
	"github.com/stretchr/testify/assert"
)

func setup() {
	viper.SetConfigFile("../config.yaml") // полный путь к файлу конфигурации
	if err := viper.ReadInConfig(); err != nil {
		log.Printf("Config file not found, using default settings: %v", err)
	}

	// Инициализируем in-memory базу данных
	database.InitInMemoryDB()
	database.ClearTasks() // Очищаем задачи перед каждым тестом
}

func TestCreateTask(t *testing.T) {
	setup() // Настраиваем in-memory базу данных перед тестом

	task := models.Task{
		Title:       "Test Task",
		Description: "Test Description",
		Status:      "in-progress",
	}

	taskJSON, _ := json.Marshal(task)
	req, err := http.NewRequest("POST", "/tasks", bytes.NewBuffer(taskJSON))
	assert.NoError(t, err)

	rr := httptest.NewRecorder()
	handler := http.HandlerFunc(handlers.CreateTask)
	handler.ServeHTTP(rr, req)

	assert.Equal(t, http.StatusOK, rr.Code)

	var createdTask models.Task
	err = json.NewDecoder(rr.Body).Decode(&createdTask)
	assert.NoError(t, err)
	assert.Equal(t, "Test Task", createdTask.Title)
	assert.Equal(t, "in-progress", createdTask.Status)
}

func TestGetTasks(t *testing.T) {
	setup()

	// Добавляем тестовую задачу
	task := models.Task{
		Title:       "Task 1",
		Description: "Test Task 1",
		Status:      "in-progress",
		CreatedAt:   time.Now(),
		UpdatedAt:   time.Now(),
	}
	database.AddTask(task)

	req, err := http.NewRequest("GET", "/tasks", nil)
	assert.NoError(t, err)

	rr := httptest.NewRecorder()
	handler := http.HandlerFunc(handlers.GetTasks)
	handler.ServeHTTP(rr, req)

	assert.Equal(t, http.StatusOK, rr.Code)

	var tasks []models.Task
	err = json.NewDecoder(rr.Body).Decode(&tasks)
	assert.NoError(t, err)
	assert.Len(t, tasks, 1)
	assert.Equal(t, "Task 1", tasks[0].Title)
}

func TestUpdateTask(t *testing.T) {
	setup()

	// Добавляем тестовую задачу
	task := models.Task{
		ID:          1,
		Title:       "Original Task",
		Description: "Original Description",
		Status:      "in-progress",
		CreatedAt:   time.Now(),
		UpdatedAt:   time.Now(),
	}
	database.AddTask(task)

	updatedTask := models.Task{
		Title:       "Updated Task",
		Description: "Updated Description",
		Status:      "completed",
	}

	taskJSON, _ := json.Marshal(updatedTask)
	req, err := http.NewRequest("PUT", "/tasks/1", bytes.NewBuffer(taskJSON))
	assert.NoError(t, err)

	rr := httptest.NewRecorder()
	router := mux.NewRouter()
	router.HandleFunc("/tasks/{id}", handlers.UpdateTask).Methods("PUT")
	router.ServeHTTP(rr, req)

	assert.Equal(t, http.StatusOK, rr.Code)

	var returnedTask models.Task
	err = json.NewDecoder(rr.Body).Decode(&returnedTask)
	assert.NoError(t, err)
	assert.Equal(t, "Updated Task", returnedTask.Title)
	assert.Equal(t, "completed", returnedTask.Status)
}

func TestDeleteTask(t *testing.T) {
	setup()

	// Добавляем тестовую задачу
	task := models.Task{
		ID:          1,
		Title:       "Task to Delete",
		Description: "Description to Delete",
		Status:      "in-progress",
		CreatedAt:   time.Now(),
		UpdatedAt:   time.Now(),
	}
	database.AddTask(task)

	req, err := http.NewRequest("DELETE", "/tasks/1", nil)
	assert.NoError(t, err)

	rr := httptest.NewRecorder()
	router := mux.NewRouter()
	router.HandleFunc("/tasks/{id}", handlers.DeleteTask).Methods("DELETE")
	router.ServeHTTP(rr, req)

	assert.Equal(t, http.StatusNoContent, rr.Code)

	// Проверяем, что задача была удалена
	tasks := database.GetTasks()
	assert.Len(t, tasks, 0)
}

func TestTaskHTML(t *testing.T) {
	setup()

	// Добавляем тестовую задачу
	task := models.Task{
		ID:          1,
		Title:       "Task 1",
		Description: "Test Task 1",
		Status:      "in-progress",
		CreatedAt:   time.Now(),
		UpdatedAt:   time.Now(),
	}
	database.AddTask(task)

	req, err := http.NewRequest("GET", "/tasks/1/view", nil)
	assert.NoError(t, err)

	rr := httptest.NewRecorder()
	router := mux.NewRouter()
	router.HandleFunc("/tasks/{id}/view", handlers.TaskHTML).Methods("GET")
	router.ServeHTTP(rr, req)

	assert.Equal(t, http.StatusOK, rr.Code)
	assert.Contains(t, rr.Body.String(), "Task 1")
}

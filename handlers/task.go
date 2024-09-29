package handlers

import (
	"encoding/json"
	"github.com/gorilla/mux"
	"github.com/mussyaroslav/otus_final_project/database"
	"github.com/mussyaroslav/otus_final_project/models"
	"log"
	"net/http"
	"strconv"
	"text/template"
	"time"
)

// Обработчик для главной страницы
func HomePage(w http.ResponseWriter, r *http.Request) {
	http.ServeFile(w, r, "./web/index.html")
}

// Получение всех задач
func GetTasks(w http.ResponseWriter, r *http.Request) {
	status := r.URL.Query().Get("status")
	createdAtStr := r.URL.Query().Get("created_at")

	tasks := database.GetTasks()
	var filteredTasks []models.Task

	// Если передан фильтр по статусу
	if status != "" {
		for _, task := range tasks {
			if task.Status == status {
				filteredTasks = append(filteredTasks, task)
			}
		}
	} else {
		filteredTasks = tasks
	}

	// Если передан фильтр по дате создания
	if createdAtStr != "" {
		createdAt, err := time.Parse("2006-01-02", createdAtStr)
		if err == nil {
			var dateFilteredTasks []models.Task
			for _, task := range filteredTasks {
				if task.CreatedAt.Format("2006-01-02") == createdAt.Format("2006-01-02") {
					dateFilteredTasks = append(dateFilteredTasks, task)
				}
			}
			filteredTasks = dateFilteredTasks
		}
	}

	json.NewEncoder(w).Encode(filteredTasks)
}

// Добавление новой задачи
func CreateTask(w http.ResponseWriter, r *http.Request) {
	var task models.Task
	if err := json.NewDecoder(r.Body).Decode(&task); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	// Если статус не передан, присваиваем статус "in-progress" по умолчанию
	if task.Status == "" {
		task.Status = "in-progress"
	}

	task.ID = uint(len(database.GetTasks()) + 1)
	task.CreatedAt = time.Now()
	task.UpdatedAt = time.Now()

	database.AddTask(task)
	json.NewEncoder(w).Encode(task)
}

// Обновление задачи по ID
func UpdateTask(w http.ResponseWriter, r *http.Request) {
	params := mux.Vars(r)
	id, _ := strconv.Atoi(params["id"])

	var updatedTask models.Task
	if err := json.NewDecoder(r.Body).Decode(&updatedTask); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	for _, task := range database.GetTasks() {
		if task.ID == uint(id) {
			updatedTask.ID = task.ID
			updatedTask.CreatedAt = task.CreatedAt
			updatedTask.UpdatedAt = time.Now()

			if database.UpdateTask(uint(id), updatedTask) {
				json.NewEncoder(w).Encode(updatedTask)
				return
			}
		}
	}
	http.Error(w, "Task not found", http.StatusNotFound)
}

// Удаление задачи по ID
func DeleteTask(w http.ResponseWriter, r *http.Request) {
	params := mux.Vars(r)
	id, _ := strconv.Atoi(params["id"])

	if database.DeleteTask(uint(id)) {
		w.WriteHeader(http.StatusNoContent)
	} else {
		http.Error(w, "Task not found", http.StatusNotFound)
	}
}

// Обработчик для рендеринга страницы задачи
func TaskHTML(w http.ResponseWriter, r *http.Request) {
	params := mux.Vars(r)
	id, err := strconv.Atoi(params["id"])
	if err != nil {
		http.Error(w, "Invalid task ID", http.StatusBadRequest)
		log.Printf("Invalid task ID: %v", err)
		return
	}

	// Поиск задачи по ID
	var foundTask *models.Task
	for _, task := range database.GetTasks() {
		if task.ID == uint(id) {
			foundTask = &task
			break
		}
	}

	if foundTask == nil {
		http.Error(w, "Task not found", http.StatusNotFound)
		log.Printf("Task with ID %d not found", id)
		return
	}

	// Парсинг шаблона
	tmpl, err := template.ParseFiles("../web/task.html")
	if err != nil {
		http.Error(w, "Error loading template", http.StatusInternalServerError)
		log.Printf("Error loading template: %v", err)
		return
	}

	// Рендеринг шаблона с данными задачи
	err = tmpl.Execute(w, foundTask)
	if err != nil {
		http.Error(w, "Error rendering template", http.StatusInternalServerError)
		log.Printf("Error rendering template: %v", err)
	}
}

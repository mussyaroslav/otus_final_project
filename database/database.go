package database

import (
	"encoding/json"
	"github.com/mussyaroslav/otus_final_project/models"
	"github.com/spf13/viper"
	"log"
	"os"
)

var tasks []models.Task

// InitDB инициализирует JSON файл для хранения задач
func InitDB() {
	filePath := viper.GetString("database.file")

	// Проверяем, существует ли файл
	if _, err := os.Stat(filePath); os.IsNotExist(err) {
		log.Printf("File %s not found, creating new one...", filePath)
		// Создаем пустой файл
		file, err := os.Create(filePath)
		if err != nil {
			log.Fatalf("Failed to create file: %v", err)
		}
		defer file.Close()

		// Записываем пустой массив задач
		file.Write([]byte("[]"))
	}

	// Читаем существующие задачи из файла
	file, err := os.Open(filePath)
	if err != nil {
		log.Fatalf("Failed to open file: %v", err)
	}
	defer file.Close()

	// Декодируем задачи из JSON в переменную tasks
	decoder := json.NewDecoder(file)
	if err := decoder.Decode(&tasks); err != nil {
		log.Fatalf("Failed to decode tasks from file: %v", err)
	}
}

// SaveTasks сохраняет текущий список задач в JSON файл
func SaveTasks() {
	filePath := viper.GetString("database.file")

	file, err := os.Create(filePath)
	if err != nil {
		log.Fatalf("Failed to open file for writing: %v", err)
	}
	defer file.Close()

	encoder := json.NewEncoder(file)
	if err := encoder.Encode(tasks); err != nil {
		log.Fatalf("Failed to encode tasks to file: %v", err)
	}
}

// GetTasks возвращает все задачи
func GetTasks() []models.Task {
	return tasks
}

// AddTask добавляет новую задачу
func AddTask(task models.Task) {
	tasks = append(tasks, task)
	SaveTasks()
}

// UpdateTask обновляет существующую задачу
func UpdateTask(id uint, updatedTask models.Task) bool {
	for i, task := range tasks {
		if task.ID == id {
			tasks[i] = updatedTask
			SaveTasks()
			return true
		}
	}
	return false
}

// DeleteTask удаляет задачу по ID
func DeleteTask(id uint) bool {
	for i, task := range tasks {
		if task.ID == id {
			tasks = append(tasks[:i], tasks[i+1:]...)
			SaveTasks()
			return true
		}
	}
	return false
}

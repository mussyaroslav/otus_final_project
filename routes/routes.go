package routes

import (
	"github.com/gorilla/mux"
	"github.com/mussyaroslav/otus_final_project/handlers"
	"net/http"
)

func InitRoutes() *mux.Router {
	r := mux.NewRouter()

	// Обработчик для главной страницы
	r.HandleFunc("/", handlers.HomePage).Methods("GET")

	// Маршрут для получения задачи в HTML формате
	r.HandleFunc("/tasks/{id}/view", handlers.TaskHTML).Methods("GET") // Добавляем /view для HTML

	// Получение, создание, обновление и удаление задач (JSON формат)
	r.HandleFunc("/tasks", handlers.GetTasks).Methods("GET")           // Для списка задач
	r.HandleFunc("/tasks", handlers.CreateTask).Methods("POST")        // Для создания задачи
	r.HandleFunc("/tasks/{id}", handlers.UpdateTask).Methods("PUT")    // Для обновления задачи
	r.HandleFunc("/tasks/{id}", handlers.DeleteTask).Methods("DELETE") // Для удаления задачи

	// Обслуживание статических файлов (например, HTML, CSS)
	r.PathPrefix("/web/").Handler(http.StripPrefix("/web/", http.FileServer(http.Dir("./web/"))))

	return r
}

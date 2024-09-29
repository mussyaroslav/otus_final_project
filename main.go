package main

import (
	"github.com/mussyaroslav/otus_final_project/config"
	"github.com/mussyaroslav/otus_final_project/database"
	"github.com/mussyaroslav/otus_final_project/routes"
	"log"
	"net/http"
)

func main() {
	// Инициализация конфигурации
	config.LoadConfig()

	// Инициализация базы данных
	database.InitDB()

	// Создание маршрутов
	r := routes.InitRoutes()

	// Обслуживание статических файлов
	fs := http.FileServer(http.Dir("./web"))
	http.Handle("/web/", http.StripPrefix("/web/", fs))

	// Запуск HTTP-сервера
	log.Println("Server starting on :8080...")
	log.Fatal(http.ListenAndServe(":8080", r))
}

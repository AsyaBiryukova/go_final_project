package main

import (
	"fmt"
	"log"
	"net/http"
	"os"

	"github.com/AsyaBiryukova/go_final_project/internal/App/models"
	"github.com/AsyaBiryukova/go_final_project/internal/api"
	"github.com/AsyaBiryukova/go_final_project/internal/server"
	db "github.com/AsyaBiryukova/go_final_project/internal/sqlite3"

	"github.com/joho/godotenv"
	_ "github.com/mattn/go-sqlite3"
)

func main() {
	// Загружаем переменные среды
	err := godotenv.Load(".env")
	if err != nil {
		fmt.Println(err)
	}
	// .env сам подгружается если мы используем docker compose для запуска, но для тестов удобнее запускать код напрямую, поэтому оставил godotenv

	dbFile := os.Getenv("TODO_DBFILE")

	// Если бд не существует, создаём
	if !db.DbExists(dbFile) {
		err = db.InstallDB()
		if err != nil {
			log.Println(err)
		}
	}

	// Запуск бд
	dbStorage, err := db.StartDB()
	defer func() {
		err := dbStorage.CloseDB()
		if err != nil {
			log.Println(err)
		}
	}()

	if err != nil {
		log.Fatal(err)
	}

	DateFormat := os.Getenv("TODO_DATEFORMAT")

	App := models.NewApp(dbStorage, DateFormat)

	Api := api.NewApi(*App)

	// Адрес для запуска сервера
	ip := ""
	port := os.Getenv("TODO_PORT")
	addr := fmt.Sprintf("%s:%s", ip, port)

	//Router
	r := server.NewRouter(*Api)
	// Запуск сервера
	err = http.ListenAndServe(addr, r)
	if err != nil {
		log.Println(err)
	}
	log.Printf("Server running on %s\n", port)

}

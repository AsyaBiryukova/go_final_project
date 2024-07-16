package models

import (
	"fmt"
	"log"
	"regexp"
	"strings"
	"time"

	nd "github.com/AsyaBiryukova/go_final_project/internal/nextdate"
)

type Repository interface {
	GetTasksList() ([]Task, error)
	GetTasksByTitle(search string) ([]Task, error)
	GetTasksByDate(search ...string) ([]Task, error)
	AddTask(task Task) (int64, error)
	PutTask(updateTask Task) error
	GetTaskByID(id string) (Task, error)
	DeleteTask(id string) error
}

func NewApp(repo2 Repository, df string) App {
	return App{repo: repo2, DateFormat: df}
}

type App struct {
	repo       Repository
	DateFormat string
}

type Task struct {
	ID      string `json:"id"`
	Date    string `json:"date"`
	Title   string `json:"title"`
	Comment string `json:"comment"`
	Repeat  string `json:"repeat"`
}

// formatTask проверяет переданную задачу Task на корректность полей, а так же корректирует дату задачи.
// Возвращает отформатированную задачу или ошибку.
func (a App) FormatTask(task Task) (Task, error) {
	var date time.Time
	var err error

	if len(task.Date) == 0 || strings.ToLower(task.Date) == "today" {
		date = time.Now()
		task.Date = date.Format(a.DateFormat)

	} else {
		date, err = time.Parse(a.DateFormat, task.Date)
		if err != nil {
			log.Println(err)
			return Task{}, err
		}
	}
	if isID, _ := regexp.Match("[0-9]+", []byte(task.ID)); !isID && task.ID != "" {
		err = fmt.Errorf("некорректный формат ID")
		return Task{}, err
	}

	// Даты с временем приведённым к 00:00:00
	dateTrunc := date.Truncate(time.Hour * 24)
	nowTrunc := time.Now().Truncate(time.Hour * 24)

	if dateTrunc.Before(nowTrunc) {
		switch {
		case len(task.Repeat) > 0:
			task.Date, err = nd.NextDate(time.Now(), task.Date, task.Repeat)
			if err != nil {
				log.Println(err)
				return Task{}, err
			}
		case len(task.Repeat) == 0:
			task.Date = time.Now().Format(a.DateFormat)
		}

	}
	return task, nil
}

func (a App) SearchTasks(search string) ([]Task, error) {
	// Проверяем может ли поисковой запрос содержать поиск по дате
	isDate, _ := regexp.Match("[0-9]{2}.[0-9]{2}.[0-9]{4}", []byte(search))

	if len(search) == 0 {
		return a.repo.GetTasksList()
	}
	if isDate {
		date, err := time.Parse("02.01.2006", search)

		if err == nil {
			search = date.Format(a.DateFormat)
			return a.repo.GetTasksByDate(search)
		}

	}
	search = fmt.Sprint("%" + search + "%")
	return a.repo.GetTasksByTitle(search)

}

func (a App) AddTask(task Task) (int64, error) {
	return a.repo.AddTask(task)
}

func (a App) PutTask(updateTask Task) error {
	return a.repo.PutTask(updateTask)
}

func (a App) GetTaskByID(id string) (Task, error) {
	return a.repo.GetTaskByID(id)
}

func (a App) DeleteTask(id string) error {
	return a.repo.DeleteTask(id)
}

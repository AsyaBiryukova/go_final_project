package api

import (
	"bytes"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"strconv"

	"time"

	models "github.com/AsyaBiryukova/go_final_project/internal/App/models"
	nd "github.com/AsyaBiryukova/go_final_project/internal/nextdate"
)

type Api struct {
	app models.App
}

func NewApi(app2 models.App) Api {
	return Api{app: app2}
}

// task.go содержит обработчики запросов к api/task

// PostTaskHandler обрабатывает запрос с методом POST.
// Если пользователь авторизован и задача отправлена в корректном формате, добавляет новую задачу в базу данных.
// Возвращает JSON {"id": string} или JSON {"error": error} в случае ошибки.
func (a Api) TaskHandler(w http.ResponseWriter, r *http.Request) {
	method := r.Method
	switch method {
	case http.MethodGet:
		a.getTask(w, r)
	case http.MethodPost:
		a.postTask(w, r)
	case http.MethodPut:
		a.putTask(w, r)
	case http.MethodDelete:
		a.deleteTask(w, r)
	}
}

func (a Api) postTask(w http.ResponseWriter, r *http.Request) {
	var task models.Task
	var buf bytes.Buffer
	var err error
	var id int64

	write := func() {
		w.Header().Set("Content-Type", "application/json; charset=UTF-8")
		if err != nil {
			writeErr(err, w)
			return
		} else {
			idResp := map[string]string{
				"id": strconv.Itoa(int(id)),
			}
			resp, err := json.Marshal(idResp)
			if err != nil {
				log.Println(err)
			}
			w.WriteHeader(http.StatusCreated)
			_, err = w.Write(resp)
			if err != nil {
				log.Println(err)
			}
			return
		}

	}

	_, err = buf.ReadFrom(r.Body)
	if err != nil {
		write()
		return
	}

	if err = json.Unmarshal(buf.Bytes(), &task); err != nil {
		write()
		return
	}

	task, err = a.app.FormatTask(task)
	if err != nil {
		write()
		return
	}

	id, err = a.app.AddTask(task)
	write()
}

// PutTaskHandler обрабатывает запрос с методом PUT.
// Если пользователь авторизован и задача существует, и отправлена в корректном формате, обновляет поля задачи в базе данных.
// Возвращает пустой JSON {} или JSON {"error": error} в случае ошибки.
func (a Api) putTask(w http.ResponseWriter, r *http.Request) {
	var updatedTask models.Task
	var buf bytes.Buffer
	var err error

	write := func() {
		w.Header().Set("Content-Type", "application/json; charset=UTF-8")
		if err != nil {
			writeErr(err, w)
			return
		} else {
			writeEmptyJson(w)
			return
		}

	}

	_, err = buf.ReadFrom(r.Body)
	if err != nil {
		write()
		return
	}

	if err = json.Unmarshal(buf.Bytes(), &updatedTask); err != nil {
		write()
		return
	}

	updatedTask, err = a.app.FormatTask(updatedTask)
	if err != nil {
		write()
		return
	}

	err = a.app.PutTask(updatedTask)
	write()

}

// GetTaskHandler обрабатывает запрос с методом GET.
// Если пользователь авторизован, возвращает задачу с указанным ID.
// Возвращает JSON {"task":Task}, или JSON {"error": error} при ошибке.
func (a Api) getTask(w http.ResponseWriter, r *http.Request) {
	var err error
	var task models.Task

	write := func() {
		w.Header().Set("Content-Type", "application/json; charset=UTF-8")
		var resp []byte
		if err != nil {
			writeErr(err, w)
			return
		} else {
			resp, err = json.Marshal(task)
		}

		if err != nil {
			log.Println(err)
		}
		w.WriteHeader(http.StatusOK)
		_, err = w.Write(resp)
		if err != nil {
			log.Println(err)
		}

	}

	q := r.URL.Query()
	id := q.Get("id")

	task, err = a.app.GetTaskByID(id)
	if err != nil {
		log.Println(err)
	}
	write()

}

// DeleteTaskHandler обрабатывает запрос к api/task с методом DELETE.
// Если пользователь авторизован и id существует, удаляет задачу.
// При успешном выполнение возвращает пустой JSON {}. Иначе возвращает JSON {"error":error}.
func (a Api) deleteTask(w http.ResponseWriter, r *http.Request) {
	var err error

	q := r.URL.Query()
	id := q.Get("id")
	isID := isID(id)
	if !isID {
		writeErr(fmt.Errorf("некорректный формат id"), w)
		return
	}

	err = a.app.DeleteTask(id)
	if err != nil {
		writeErr(err, w)
		return
	}
	writeEmptyJson(w)

}

// GetTasksHandler обрабатывает запросы к /api/tasks с методом GET.
// Если пользователь авторизован, возвращает JSON {"tasks": Task} содержащий последние добавленные задачи, или
// последние добавленные задачи соответствующие поисковому запросу search. В случае ошибки возвращает JSON {"error": error}.
func (a Api) GetTasksHandler(w http.ResponseWriter, r *http.Request) {
	var tasks []models.Task
	var err error

	// write отправляет клиенту ответ либо ошибку, в формате json
	write := func() {
		w.Header().Set("Content-Type", "application/json; charset=UTF-8")
		var resp []byte
		if err != nil {
			writeErr(err, w)
			return
		} else {
			if len(tasks) == 0 {
				tasksResp := map[string][]models.Task{
					"tasks": {},
				}
				resp, err = json.Marshal(tasksResp)
			} else {
				tasksResp := map[string][]models.Task{
					"tasks": tasks,
				}
				resp, err = json.Marshal(tasksResp)

			}

			if err != nil {
				log.Println(err)
			}
			w.WriteHeader(http.StatusCreated)
			_, err = w.Write(resp)
			if err != nil {
				log.Println(err)
			}
			return
		}
	}

	// Проверяем есть ли поисковой зарпос
	q := r.URL.Query()
	search := q.Get("search")

	tasks, err = a.app.SearchTasks(search)
	if err != nil {
		log.Println(err)
	}

	write()

}

// PostTaskDoneHandler обрабатывает запросы к /api/task/done с методом POST.
// Если пользователь авторизован, удаляет задачи не имеющих правил повторения repeat, или обновляет дату выполнения задач, имеющих правило repeat.
// Возвращает пустой JSON {} в случае успеха, или JSON {"error": error} при возникновение ошибки.
func (a Api) PostTaskDoneHandler(w http.ResponseWriter, r *http.Request) {
	var err error

	q := r.URL.Query()
	id := q.Get("id")
	isID := isID(id)
	if !isID {
		writeErr(fmt.Errorf("некорректный формат id"), w)
		return
	}
	task, err := a.app.GetTaskByID(id)
	if err != nil {
		writeErr(err, w)
		return
	}
	if len(task.Repeat) == 0 {
		err = a.app.DeleteTask(id)
		if err != nil {
			writeErr(err, w)
			return
		}
		writeEmptyJson(w)
		return
	} else {
		nextDate, err := nd.NextDate(time.Now(), task.Date, task.Repeat)
		if err != nil {
			writeErr(err, w)
			return
		}
		task.Date = nextDate
	}
	err = a.app.PutTask(task)
	if err != nil {
		writeErr(err, w)
		return
	}
	writeEmptyJson(w)

}

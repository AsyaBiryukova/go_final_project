package api

import (
	"bytes"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"strconv"

	"github.com/AsyaBiryukova/go_final_project/internal/db"
)

// task.go содержит обработчики запросов к api/task

// PostTaskHandler обрабатывает запрос с методом POST.
// Если пользователь авторизован и задача отправлена в корректном формате, добавляет новую задачу в базу данных.
// Возвращает JSON {"id": string} или JSON {"error": error} в случае ошибки.
func TaskHandler(w http.ResponseWriter, r *http.Request) {
	method := r.Method
	switch method {
	case http.MethodGet:
		getTask(w, r)
	case http.MethodPost:
		postTask(w, r)
	case http.MethodPut:
		putTask(w, r)
	case http.MethodDelete:
		deleteTask(w, r)
	}
}

func writeError(w http.ResponseWriter, err error, statusCode int) {
	w.Header().Set("Content-Type", "application/json; charset=UTF-8")
	w.WriteHeader(statusCode)
	resp, _ := json.Marshal(map[string]string{"error": err.Error()})
	w.Write(resp)
}

func writeResponse(w http.ResponseWriter, data interface{}, statusCode int) {
	w.Header().Set("Content-Type", "application/json; charset=UTF-8")
	w.WriteHeader(statusCode)
	resp, err := json.Marshal(data)
	if err != nil {
		writeError(w, err, http.StatusInternalServerError)
		return
	}
	w.Write(resp)
}

// Пример использования в postTask
func postTask(w http.ResponseWriter, r *http.Request) {
	var task db.Task
	var buf bytes.Buffer
	var id int64

	_, err := buf.ReadFrom(r.Body)
	if err != nil {
		writeError(w, err, http.StatusBadRequest)
		return
	}

	if err = json.Unmarshal(buf.Bytes(), &task); err != nil {
		writeError(w, err, http.StatusBadRequest)
		return
	}

	task, err = task.FormatTask()
	if err != nil {
		writeError(w, err, http.StatusBadRequest)
		return
	}

	if task.Title == "" {
		writeError(w, fmt.Errorf("missing required task fields"), http.StatusBadRequest)
		return
	}

	id, err = dbs.AddTask(task)
	if err != nil {
		writeError(w, err, http.StatusInternalServerError)
		return
	}

	writeResponse(w, map[string]string{"id": strconv.Itoa(int(id))}, http.StatusCreated)
}

// PutTaskHandler обрабатывает запрос с методом PUT.
// Если пользователь авторизован и задача существует, и отправлена в корректном формате, обновляет поля задачи в базе данных.
// Возвращает пустой JSON {} или JSON {"error": error} в случае ошибки.
func putTask(w http.ResponseWriter, r *http.Request) {
	var updatedTask db.Task
	var buf bytes.Buffer
	var err error

	_, err = buf.ReadFrom(r.Body)
	if err != nil {
		writeError(w, err, http.StatusBadRequest)
		return
	}

	if err = json.Unmarshal(buf.Bytes(), &updatedTask); err != nil {
		writeError(w, err, http.StatusBadRequest)
		return
	}

	updatedTask, err = updatedTask.FormatTask()
	if err != nil {
		writeError(w, err, http.StatusBadRequest)
		return
	}

	if updatedTask.Title == "" {
		writeError(w, fmt.Errorf("title cannot be empty"), http.StatusBadRequest)
		return
	}

	err = dbs.PutTask(updatedTask)
	if err != nil {
		writeError(w, err, http.StatusInternalServerError)
		return
	}

	writeResponse(w, struct{}{}, http.StatusOK)
}

// GetTaskHandler обрабатывает запрос с методом GET.
// Если пользователь авторизован, возвращает задачу с указанным ID.
// Возвращает JSON {"task":Task}, или JSON {"error": error} при ошибке.
func getTask(w http.ResponseWriter, r *http.Request) {
	var task db.Task

	q := r.URL.Query()
	id := q.Get("id")

	task, err := dbs.GetTaskByID(id)
	if err != nil {
		writeError(w, err, http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json; charset=UTF-8")
	w.WriteHeader(http.StatusOK)
	resp, err := json.Marshal(task)
	if err != nil {
		log.Println(err)
	}
	_, err = w.Write(resp)
	if err != nil {
		log.Println(err)
	}
}

// DeleteTaskHandler обрабатывает запрос к api/task с методом DELETE.
// Если пользователь авторизован и id существует, удаляет задачу.
// При успешном выполнение возвращает пустой JSON {}. Иначе возвращает JSON {"error":error}.
func deleteTask(w http.ResponseWriter, r *http.Request) {
	q := r.URL.Query()
	id := q.Get("id")

	err := dbs.DeleteTask(id)
	if err != nil {
		writeError(w, err, http.StatusInternalServerError)
		return
	}

	writeResponse(w, struct{}{}, http.StatusOK)
}

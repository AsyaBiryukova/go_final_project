package server

import (
	"net/http"

	"github.com/AsyaBiryukova/go_final_project/internal/api"
	"github.com/AsyaBiryukova/go_final_project/internal/auth"
	"github.com/go-chi/chi/v5"
)

func NewRouter(Api Api) *chi.Mux {
	r := chi.NewRouter()

	r.Handle("/*", http.FileServer(http.Dir("./web")))

	r.Get("/api/nextdate", api.GetNextDateHandler)
	r.Get("/api/tasks", auth.Auth(api.GetTasksHandler))
	r.Post("/api/task/done", auth.Auth(api.PostTaskDoneHandler))
	r.Post("/api/signin", auth.Auth(api.PostSigninHandler))
	r.Handle("/api/task", auth.Auth(api.TaskHandler))

	return r
}

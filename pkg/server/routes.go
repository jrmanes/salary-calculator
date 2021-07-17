package server

import (
	"net/http"

	"github.com/go-chi/chi"

	"github.com/jrmanes/my-salary/pkg/salary"
)

func createRouter() http.Handler {
	r := chi.NewRouter()

	var s = salary.Salary{}

	r.Post("/", s.CreateHandler)

	return r
}

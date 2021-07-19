package server

import (
	"net/http"

	"github.com/go-chi/chi"

	"github.com/jrmanes/salary-calculator/pkg/salary"
)

func createRouter() http.Handler {
	r := chi.NewRouter()

	var s = salary.Salary{}

	r.Post("/", s.CalculateSalaryHandler)

	return r
}

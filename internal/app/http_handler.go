package app

import (
	"encoding/json"
	"errors"
	"github.com/SomchaiSPB/user-auth/internal/service"
	"io"
	"log"
	"net/http"
	"strconv"
)

const (
	defaultPage    = 1
	defaultPerPage = 1000
)

type ErrorResponse struct {
	Message string `json:"message"`
}

func respondWithErr(w http.ResponseWriter, err error, code int) {
	e := ErrorResponse{
		Message: err.Error(),
	}

	byteResp, mErr := json.Marshal(e)

	if mErr != nil {
		log.Printf("error response json marshal error: %v", mErr)
		return
	}

	w.WriteHeader(code)
	w.Write(byteResp)
}

func (a *App) HandleCreateUser(w http.ResponseWriter, r *http.Request) {
	data, err := io.ReadAll(r.Body)

	if err != nil {
		respondWithErr(w, err, http.StatusInternalServerError)
		return
	}

	u, err := a.userSvc.Create(data)

	if err != nil {
		code := http.StatusInternalServerError

		if errors.Is(err, service.ErrUserNameExists) || errors.Is(err, service.ErrValidation) {
			code = http.StatusBadRequest
		}

		respondWithErr(w, err, code)
		return
	}

	w.WriteHeader(http.StatusCreated)
	w.Write(u)
}

func (a *App) HandleAuthUser(w http.ResponseWriter, r *http.Request) {
	data, err := io.ReadAll(r.Body)

	if err != nil {
		respondWithErr(w, err, http.StatusInternalServerError)
		return
	}

	response, err := a.userSvc.Authenticate(data, a.config.AuthJwtSecret())

	if err != nil {
		code := http.StatusInternalServerError

		if errors.Is(err, service.ErrWrongCredentials) || errors.Is(err, service.ErrValidation) {
			code = http.StatusBadRequest
		}

		respondWithErr(w, err, code)
		return
	}

	w.WriteHeader(http.StatusOK)
	w.Write(response)
}

func (a *App) HandleGetProduct(w http.ResponseWriter, r *http.Request) {
	productName := r.URL.Query().Get("name")

	product, err := a.productSvc.GetProduct(productName)

	if err != nil {
		code := http.StatusInternalServerError

		if errors.Is(err, service.ErrEmptyRequestName) {
			code = http.StatusBadRequest
		}
		if errors.Is(err, service.ErrProductNotFound) {
			code = http.StatusNotFound
		}

		respondWithErr(w, err, code)
		return
	}

	w.WriteHeader(http.StatusOK)
	w.Write(product)
}

func (a *App) HandleGetProducts(w http.ResponseWriter, r *http.Request) {
	query := r.URL.Query()

	pageStr := query.Get("page")
	perPageStr := query.Get("perPage")

	page, err := strconv.Atoi(pageStr)
	if err != nil || page <= 0 {
		page = defaultPage
	}

	perPage, err := strconv.Atoi(perPageStr)
	if err != nil || perPage <= 0 {
		perPage = defaultPerPage
	}

	product, err := a.productSvc.GetProducts(page, perPage)

	if err != nil {
		code := http.StatusInternalServerError

		if errors.Is(err, service.ErrEmptyRequestName) {
			code = http.StatusBadRequest
		}
		if errors.Is(err, service.ErrProductNotFound) {
			code = http.StatusNotFound
		}

		respondWithErr(w, err, code)
		return
	}

	w.WriteHeader(http.StatusOK)
	w.Write(product)
}

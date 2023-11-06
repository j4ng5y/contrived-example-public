package app

import (
	"encoding/json"
	"io/ioutil"
	"net/http"
)

type (
	RootRequest struct {
		Name string `json:"name"`
	}

	RootResponse struct {
		Message string `json:"message"`
		Error   string `json:"error"`
	}
)

func (app *App) Mux() *http.ServeMux {
	router := http.NewServeMux()

	router.HandleFunc("/", app.RootHandler)

	return router
}

func (app *App) RootHandler(w http.ResponseWriter, r *http.Request) {
	var (
		request  = new(RootRequest)
		response = new(RootResponse)
	)

	w.Header().Add("Content-Type", "application/json")

	req, err := ioutil.ReadAll(r.Body)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}

	if err := json.Unmarshal(req, request); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}

	if err := app.storage.Insert(request.Name, ""); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}

	response.Message = request.Name

	responseBytes, err := json.Marshal(response)

	if _, err = w.Write(responseBytes); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}

	w.WriteHeader(http.StatusOK)
}

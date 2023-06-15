package server

import (
	"YamlApiServer/pkg/model"
	"net/http"
	"strings"

	"github.com/gorilla/mux"
	"gopkg.in/yaml.v3"
)

type Server struct {
	router   *mux.Router
	metadata map[string]model.Metadata
}

func NewServer() *Server {
	s := &Server{
		router:   mux.NewRouter(),
		metadata: make(map[string]model.Metadata),
	}
	s.routes()
	return s
}

func (s *Server) routes() {
	s.router.HandleFunc("/metadata", s.createMetadata).Methods("POST")
	s.router.HandleFunc("/metadata", s.searchMetadata).Methods("GET")
}

func (s *Server) createMetadata(w http.ResponseWriter, r *http.Request) {
	var m model.Metadata
	if err := yaml.NewDecoder(r.Body).Decode(&m); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	s.metadata[m.Title] = m
	w.WriteHeader(http.StatusCreated)
}

func (s *Server) searchMetadata(w http.ResponseWriter, r *http.Request) {
	query := r.URL.Query().Get("q")
	var result []model.Metadata
	for _, v := range s.metadata {
		if strings.Contains(v.Title, query) {
			result = append(result, v)
		}
	}

	if err := yaml.NewEncoder(w).Encode(result); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
}

func (s *Server) Run(addr string) {
	http.ListenAndServe(addr, s.router)
}

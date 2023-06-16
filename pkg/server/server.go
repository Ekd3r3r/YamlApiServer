package server

import (
	"YamlApiServer/pkg/model"
	"net/http"

	"github.com/asaskevich/govalidator"
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

	// Validate metadata
	if err := s.validateMetadata(m); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	s.metadata[m.Title] = m
	w.WriteHeader(http.StatusCreated)
}

func (s *Server) searchMetadata(w http.ResponseWriter, r *http.Request) {
	query := r.URL.Query()
	var result []model.Metadata
	for _, v := range s.metadata {
		match := true
		for key, values := range query {
			switch key {
			case "title":
				match = match && (v.Title == values[0])
			case "version":
				match = match && (v.Version == values[0])
			case "company":
				match = match && (v.Company == values[0])
			case "website":
				match = match && (v.Website == values[0])
			case "source":
				match = match && (v.Source == values[0])
			case "license":
				match = match && (v.License == values[0])
			case "description":
				match = match && (v.Description == values[0])
			case "maintainer.name":
				maintainerMatch := false
				for _, maintainer := range v.Maintainers {
					if maintainer.Name == values[0] {
						maintainerMatch = true
					}
				}
				match = match && maintainerMatch
			case "maintainer.email":
				maintainerMatch := false
				for _, maintainer := range v.Maintainers {
					if maintainer.Email == values[0] {
						maintainerMatch = true
					}
				}
				match = match && maintainerMatch
			}
		}
		if match {
			result = append(result, v)
		}

		if err := yaml.NewEncoder(w).Encode(result); err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
	}

}

func (s *Server) Run(addr string) {
	http.ListenAndServe(addr, s.router)
}

func (s *Server) validateMetadata(m model.Metadata) error {

	govalidator.SetFieldsRequiredByDefault(true)
	_, err := govalidator.ValidateStruct(m)
	if err != nil {
		return err
	}
	for _, maintainer := range m.Maintainers {
		_, err := govalidator.ValidateStruct(maintainer)
		if err != nil {
			return err
		}
	}
	return nil
}

package server

import (
	"YamlApiServer/pkg/model"
	"net/http"
	"strings"
	"sync"

	"github.com/asaskevich/govalidator"
	"github.com/gorilla/mux"
	"gopkg.in/yaml.v3"
)

type Server struct {
	router    *mux.Router
	metadata  map[string]model.Metadata
	dataMutex sync.RWMutex
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

	s.dataMutex.RLock()
	_, exists := s.metadata[m.Title]
	s.dataMutex.RUnlock()

	s.dataMutex.Lock()
	s.metadata[m.Title] = m
	s.dataMutex.Unlock()

	if exists {
		// Metadata exists, so it was updated.
		w.WriteHeader(http.StatusOK)
		w.Write([]byte("Metadata updated"))
	} else {
		// Metadata did not exist, so it was created.
		w.WriteHeader(http.StatusCreated)
		w.Write([]byte("Metadata created"))
	}

}

func (s *Server) searchMetadata(w http.ResponseWriter, r *http.Request) {
	query := r.URL.Query()

	s.dataMutex.RLock()
	data := s.metadata
	s.dataMutex.RUnlock()

	var result []model.Metadata
	matchType := query.Get("matchType")
	delete(query, "matchType") // Remove it so it doesn't interfere with our matching logic

	for _, v := range data {
		var match bool

		if len(query) == 0 {
			result = append(result, v)
			continue
		}

		if matchType == "and" {
			match = true // default to true for AND logic
		}

		for key, values := range query {
			fieldMatch := false
			switch key {
			case "title":
				fieldMatch = v.Title == values[0]
			case "version":
				fieldMatch = v.Version == values[0]
			case "company":
				fieldMatch = v.Company == values[0]
			case "website":
				fieldMatch = v.Website == values[0]
			case "source":
				fieldMatch = v.Source == values[0]
			case "license":
				fieldMatch = v.License == values[0]
			case "description":
				fieldMatch = v.Description == values[0]
			case "maintainer":
				for _, value := range values {
					parts := strings.SplitN(value, "-", 2)
					if len(parts) != 2 {
						continue
					}
					name, email := parts[0], parts[1]
					for _, maintainer := range v.Maintainers {
						if maintainer.Name == name && maintainer.Email == email {
							fieldMatch = true
							break
						}
					}
					if fieldMatch {
						break
					}
				}
			case "maintainer.name":
				for _, maintainer := range v.Maintainers {
					for _, value := range values {
						if maintainer.Name == value {
							fieldMatch = true
							break
						}
					}
					if fieldMatch {
						break
					}
				}
			case "maintainer.email":
				for _, maintainer := range v.Maintainers {
					for _, value := range values {
						if maintainer.Email == value {
							fieldMatch = true
							break
						}
					}
					if fieldMatch {
						break
					}
				}
			}

			if matchType == "and" {
				match = match && fieldMatch
			} else { // default to OR logic
				match = match || fieldMatch
			}
		}

		if match {
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

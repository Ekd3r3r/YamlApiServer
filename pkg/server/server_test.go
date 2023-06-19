package server

import (
	"YamlApiServer/pkg/model"
	"bytes"
	"io/ioutil"
	"net/http"
	"net/http/httptest"
	"testing"

	"gopkg.in/yaml.v3"
)

func TestCreateMetadata(t *testing.T) {
	server := NewServer()
	ts := httptest.NewServer(server.router)
	defer ts.Close()

	testCases := []struct {
		description    string
		sampleYamlPath string
		expectedSC     int
	}{
		{
			description:    "Testing Valid App 1",
			sampleYamlPath: "../../samples/ValidApp1.yaml",
			expectedSC:     201,
		},
		{
			description:    "Testing Valid App 2",
			sampleYamlPath: "../../samples/ValidApp2.yaml",
			expectedSC:     201,
		},
		{
			description:    "Testing Invalid Email App",
			sampleYamlPath: "../../samples/InvalidEmailApp.yaml",
			expectedSC:     400,
		},
		{
			description:    "Testing Missing Version App",
			sampleYamlPath: "../../samples/MissingVerApp.yaml",
			expectedSC:     400,
		},
	}

	for _, tc := range testCases {
		t.Logf(tc.description)

		metadata, err := metadataFromFile(tc.sampleYamlPath)
		if err != nil {
			t.Fatalf("Failed to load metadata: %v", err)
		}

		metadataStr, err := yaml.Marshal(metadata)
		if err != nil {
			t.Fatalf("Failed to marshal metadata: %v", err)
		}

		resp, err := http.Post(ts.URL+"/metadata", "application/x-yaml", bytes.NewBuffer(metadataStr))
		if err != nil {
			t.Fatalf("http.Post failed: %v", err)
		}

		if resp.StatusCode != tc.expectedSC {
			t.Errorf("expected status %v, got %v", http.StatusText(tc.expectedSC), resp.Status)
		}

	}
}

func TestSearchMetadata(t *testing.T) {

	testCases := []struct {
		description   string
		queryString   string
		respYamlCount int
		expected      []string
	}{
		{
			description:   "Testing Get by license",
			queryString:   "/metadata?license=Apache-2.0",
			respYamlCount: 2,
			expected:      []string{"Valid App 1", "ValidApp2"},
		},
		{
			description:   "Testing Get by Maintainer Email and Version",
			queryString:   "/metadata?maintainer.email=firstmaintainer@hotmail.com&version=0.0.1&matchType=and",
			respYamlCount: 1,
			expected:      []string{"Valid App 1"},
		},
		{
			description:   "Testing Get by Maintainer Email Or Version",
			queryString:   "/metadata?maintainer.email=firstmaintainer@hotmail.com&version=1.0.1",
			respYamlCount: 2,
			expected:      []string{"Valid App 1", "ValidApp2"},
		},
		{
			description:   "Testing Get by Maintainer Email and Name Combo",
			queryString:   "/metadata?maintainer.email=firstmaintainer@hotmail.com&maintainer.name=firstmaintainer%20app1&matchType=and",
			respYamlCount: 1,
			expected:      []string{"Valid App 1"},
		},
		{
			description:   "Testing Get by multiple Maintainer Email and Name Combo",
			queryString:   "/metadata?maintainer=firstmaintainer%20app1-firstmaintainer@hotmail.com&maintainer=secondmaintainer%20app1-secondmaintainer@hotmail.com&matchType=and",
			respYamlCount: 1,
			expected:      []string{"Valid App 1"},
		},
		{
			description:   "Testing Get by Invalid Maintainer Email and Name Combo",
			queryString:   "/metadata?maintainer=firstmaintainer%20app1-secondmaintainer@hotmail.com",
			respYamlCount: 0,
			expected:      []string{},
		},
		{
			description:   "Testing No query",
			queryString:   "/metadata",
			respYamlCount: 2,
			expected:      []string{"Valid App 1", "ValidApp2"},
		},
	}

	for _, tc := range testCases {
		server := NewServer()
		ts := httptest.NewServer(server.router)
		defer ts.Close()

		inputs := []string{
			"../../samples/ValidApp1.yaml",
			"../../samples/ValidApp2.yaml",
		}

		for _, filePath := range inputs {

			metadata, err := metadataFromFile(filePath)
			if err != nil {
				t.Fatalf("Failed to load metadata: %v", err)
			}

			metadataStr, err := yaml.Marshal(metadata)
			if err != nil {
				t.Fatalf("Failed to marshal metadata: %v", err)
			}

			_, err = http.Post(ts.URL+"/metadata", "application/x-yaml", bytes.NewBuffer(metadataStr))
			if err != nil {
				t.Fatalf("http.Post failed: %v", err)
			}

		}

		t.Logf(tc.description)

		//Make a request to the search endpoint
		resp, err := http.Get(ts.URL + tc.queryString)
		if err != nil {
			t.Fatalf("http.Get failed: %v", err)
		}
		if resp.StatusCode != http.StatusOK {
			t.Errorf("expected status 200 OK, got %v", resp.Status)
		}
		// Read the response body
		body, err := ioutil.ReadAll(resp.Body)
		if err != nil {
			t.Fatalf("could not read response: %v", err)
		}

		// Unmarshal the response (assuming it's a slice of Metadata)
		var metadatas []model.Metadata
		if err := yaml.Unmarshal(body, &metadatas); err != nil {
			t.Fatalf("could not unmarshal response: %v", err)
		}

		if len(metadatas) != tc.respYamlCount {
			t.Errorf("expected %d metadata; got %d", tc.respYamlCount, len(metadatas))
		}

		for i := range metadatas {
			if metadatas[i].Title != tc.expected[0] {
				if len(tc.expected) > 1 && metadatas[i].Title != tc.expected[1] {
					t.Errorf("expected metadata not found; got %s", metadatas[i].Title)
				}

			}

		}
	}

}

func metadataFromFile(path string) (model.Metadata, error) {
	var metadata model.Metadata
	data, err := ioutil.ReadFile(path)
	if err != nil {
		return metadata, err
	}

	err = yaml.Unmarshal(data, &metadata)
	return metadata, err
}

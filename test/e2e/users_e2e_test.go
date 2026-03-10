package e2e_test

import (
	"bytes"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"gost/src/app"
	"gost/src/modules/users/dto"

	"github.com/stretchr/testify/assert"
)

// TestUsersEndpoint_E2E validates Full Integration (Controller -> Pipe -> Service -> Repository -> Database)
func TestUsersEndpoint_E2E(t *testing.T) {
	// Setup the whole application Network Engine (using real middlewares, routes and connections)
	router := app.SetupApp()

	t.Run("Create User E2E Success", func(t *testing.T) {
		payload := dto.CreateUserDto{
			Name:  "Jonatas E2E Player",
			Email: "e2e_player@test.com",
		}
		body, _ := json.Marshal(payload)

		// Creates raw Request exactly like a Frontend / CUrl would
		req, _ := http.NewRequest(http.MethodPost, "/api/v1/users", bytes.NewBuffer(body))
		req.Header.Set("Content-Type", "application/json")

		// Http Recorder acts as our Client fetching the response
		w := httptest.NewRecorder()
		router.ServeHTTP(w, req)

		// Pipe validations passed and controller returned 201 created!
		assert.Equal(t, http.StatusCreated, w.Code)

		var response map[string]interface{}
		err := json.Unmarshal(w.Body.Bytes(), &response)

		assert.NoError(t, err)
		assert.Equal(t, "Jonatas E2E Player", response["name"])
		assert.Equal(t, "e2e_player@test.com", response["email"])

		// Ensure DB gave it an ID
		assert.NotEmpty(t, response["ID"])
	})

	t.Run("Create User Fails by Generic Payload Validator (Pipe) on Name Size", func(t *testing.T) {
		// Missing 'email' entirely and breaking 'min=3' rule on Name
		payload := dto.CreateUserDto{
			Name: "A",
		}
		body, _ := json.Marshal(payload)

		req, _ := http.NewRequest(http.MethodPost, "/api/v1/users", bytes.NewBuffer(body))
		req.Header.Set("Content-Type", "application/json")

		w := httptest.NewRecorder()
		router.ServeHTTP(w, req)

		// Fails and returns our exact JSON structure defined globally
		assert.Equal(t, http.StatusBadRequest, w.Code)
		assert.Contains(t, w.Body.String(), "error\"")
		assert.Contains(t, w.Body.String(), "Bad Request")
	})

}

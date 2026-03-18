package main

import (
	"bytes"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
)

func TestPostGroup(t *testing.T) {
	gin.SetMode(gin.TestMode)

	// reset global state
	group = []People{}

	router := gin.Default()
	router.POST("/group", postGroup)

	body := People{
		ID:   "4",
		Name: "John kaisen",
	}

	jsonValue, err := json.Marshal(body)
	assert.NoError(t, err)

	req, err := http.NewRequest(http.MethodPost, "/group", bytes.NewBuffer(jsonValue))
	assert.NoError(t, err)
	req.Header.Set("Content-Type", "application/json")

	rec := httptest.NewRecorder()
	router.ServeHTTP(rec, req)

	assert.Equal(t, http.StatusCreated, rec.Code)
	assert.JSONEq(t, `{"id":"4","name":"John kaisen"}`, rec.Body.String())

	assert.Len(t, group, 1)
	assert.Equal(t, body, group[0])
}

package main_test

import (
	"net/http"
	"net/http/httptest"
	"testing"
	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
)


func TestPingPong(t *testing.T) {
	router := gin.Default()
	setupRoutes(router, nil)

	w := httptest.NewRecorder()
	req, _ := http.NewRequest("GET", "/", nil)
	router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusOK, w.Code)
	expectedBody := `{"ping":"pong"}`
	assert.Equal(t, expectedBody, w.Body.String())

}


package server

import (
	"github.com/stretchr/testify/assert"
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestHandleIndex(t *testing.T) {
	server := New(NewCofig())
	recorder := httptest.NewRecorder()
	request, _ := http.NewRequest(http.MethodGet, "/", nil)

	server.handleIndex().ServeHTTP(recorder, request)
	assert.Equal(t, "test", recorder.Body.String())
}

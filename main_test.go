package main

import (
	"fmt"
	"net/http"
	"net/http/httptest"
	"os"
	"strings"
	"testing"

	"github.com/auyer/fastgate/db"
	"github.com/labstack/echo"
	"github.com/stretchr/testify/assert"
)

var endpointJSON = `{"address" : "http:/localhost:8080","resource": "localapi"}`

func TestNewResource(t *testing.T) {
	var err error
	database, err = db.Init("./main.unittest.fastgate.db")
	if err != nil {
		t.Error(err)
		fmt.Println("Failed to create test Database")
	}
	defer database.Close()
	defer os.RemoveAll("./main.unittest.fastgate.db")
	// Setup
	e := echo.New()
	req := httptest.NewRequest(echo.POST, "/fastgate/", strings.NewReader(endpointJSON))
	req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
	rec := httptest.NewRecorder()
	c := e.NewContext(req, rec)
	// Assertions
	if assert.NoError(t, postNewEndpoint(c)) {
		assert.Equal(t, http.StatusCreated, rec.Code)
		assert.Equal(t, " ", rec.Body.String())
	}

	// TODO: Body is returning empty in Test, but not in real app
	// req2 := httptest.NewRequest(echo.GET, "http:/localhost:8080/test", strings.NewReader(""))
	// req2.Header.Set("X-fastgate-resource", "localapi")
	// rec2 := httptest.NewRecorder()
	// c2 := e.NewContext(req2, rec2)
	// // Assertions
	// if assert.NoError(t, getAllEndpoints(c2)) {
	// 	assert.Equal(t, 201, rec.Code)
	// 	assert.Equal(t, `[
	// 		{
	// 		  "address": "localapi",
	// 		  "resource": "http:/localhost:8080"
	// 		}
	// 	  ]`, rec.Body.String())
	// }

	// TODO: code is returning 201 instead of expected 307
	// req3 := httptest.NewRequest(echo.GET, "http:/localhost:8080/test", strings.NewReader(""))
	// req3.Header.Set("X-fastgate-resource", "localapi")
	// rec3 := httptest.NewRecorder()
	// c3 := e.NewContext(req3, rec3)
	// // Assertions
	// if assert.NoError(t, redirectToEndpoint(c3)) {
	// 	assert.Equal(t, http.StatusTemporaryRedirect, rec.Code)
	// 	assert.Equal(t, " ", rec.Body.String())
	// }
}

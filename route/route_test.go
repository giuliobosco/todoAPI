package route

import (
	"bytes"
	"encoding/json"
	"log"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/giuliobosco/todoAPI/config"

	"github.com/gin-gonic/gin"
	mocket "github.com/selvatico/go-mocket"
	"github.com/stretchr/testify/assert"
)

func TestBaseRoute(t *testing.T) {
	gin.SetMode(gin.TestMode)
	router := SetupRoutes()

	w := httptest.NewRecorder()
	req, _ := http.NewRequest("GET", "/", nil)
	router.ServeHTTP(w, req)

	assert.Equal(t, 200, w.Code)
	assert.Equal(t, config.SWelcome, w.Body.String())
}

func testV1LoginRoute(dbD []map[string]interface{}, httpD map[string]string) *httptest.ResponseRecorder {
	gin.SetMode(gin.TestMode)
	mocket.Catcher.Logging = true
	config.TestInit()
	router := SetupRoutes()

	// setup database
	mocket.Catcher.Reset().NewMock().WithArgs(httpD["username"], httpD["password"]).WithReply(dbD)

	// setup request
	w := httptest.NewRecorder()
	requestBody, err := json.Marshal(httpD)
	if err != nil {
		log.Fatal(err)
	}

	// execute request
	req, err := http.NewRequest("POST", "/v1/login", bytes.NewBuffer(requestBody))
	if err != nil {
		log.Fatal(err)
	}

	// serve request
	router.ServeHTTP(w, req)

	return w
}

func TestV1LoginRoute200(t *testing.T) {
	u := "T_Username"
	p := "T_Password"

	dbD := []map[string]interface{}{{"id": 1, "username": u, "password": p}}
	httpD := map[string]string{"username": u, "password": p}

	w := testV1LoginRoute(dbD, httpD)

	assert.Equal(t, 200, w.Code)
	assert.Contains(t, w.Body.String(), config.SToken)
	assert.Contains(t, w.Body.String(), config.SExpire)
}

func TestV1LoginRoute401(t *testing.T) {
	dbD := []map[string]interface{}{{"id": 0}}
	httpD := map[string]string{"username": "u1", "password": "p1"}

	w := testV1LoginRoute(dbD, httpD)

	assert.Equal(t, 401, w.Code)
	assert.Contains(t, w.Body.String(), config.SMessage)
}

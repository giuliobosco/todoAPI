package route

import (
	"io/ioutil"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/giuliobosco/todoAPI/config"
	"github.com/stretchr/testify/assert"

	"github.com/gin-gonic/gin"
)

func TestBaseRoute(t *testing.T) {
	gin.SetMode(gin.TestMode)
	r := gin.Default()
	r.GET("/", baseRoute)

	req, _ := http.NewRequest("GET", "/", nil)
	w := httptest.NewRecorder()

	r.ServeHTTP(w, req)

	statusOK := w.Code == http.StatusOK
	assert.True(t, statusOK)

	p, err := ioutil.ReadAll(w.Body)
	assert.Nil(t, err)
	assert.Contains(t, string(p), config.SWelcome)
}

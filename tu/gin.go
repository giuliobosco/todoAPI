package tu

import (
	"bytes"
	"encoding/json"
	"net/http"
	"net/http/httptest"

	"github.com/gin-gonic/gin"
)

// GetRecorderContext builds a http recorder and a gin test context
func GetRecorderContext() (*httptest.ResponseRecorder, *gin.Context) {
	w := httptest.NewRecorder()

	gin.SetMode(gin.TestMode)
	c, _ := gin.CreateTestContext(w)

	return w, c
}

// GetContext builds a gin test context with an http recorder
func GetContext() *gin.Context {
	_, c := GetRecorderContext()

	return c
}

// GetRequestPost builds a post request
func GetRequestPost(i interface{}, p string) (*http.Request, error) {
	s, err := json.Marshal(i)
	if err != nil {
		return nil, err
	}

	b := []byte(s)

	return http.NewRequest(http.MethodPost, p, bytes.NewReader(b))
}

package auth

import (
	"encoding/json"
	"errors"
	"io/ioutil"
	"log"
	"net/http"
	"os"

	"github.com/gin-gonic/gin"
	"github.com/giuliobosco/todoAPI/config"
	"github.com/giuliobosco/todoAPI/model"

	"golang.org/x/oauth2"
	"golang.org/x/oauth2/google"
)

// Credentials json configuration file rappresentation
type Credentials struct {
	Cid     string `json:"cid"`
	Csecret string `json:"csecret"`
}

// OAuthUser is the OAuth user rappresentation
type OAuthUser struct {
	Sub           string `json:"sub"`
	Name          string `json:"name"`
	GivenName     string `json:"given_name"`
	FamilyName    string `json:"family_name"`
	Picture       string `json:"picture"`
	Email         string `json:"email"`
	EmailVerified bool   `json:"email_verified"`
	Locale        string `json:"locale"`
}

var cred Credentials
var conf *oauth2.Config
var state string

// SetupOAuth configures the oauth authentication
func SetupOAuth(credsFilePath string) {
	file, err := ioutil.ReadFile(credsFilePath)
	if err != nil {
		log.Printf("File error %v\n", err)
		os.Exit(1)
	}
	json.Unmarshal(file, &cred)

	var ru string = os.Getenv("URL") + "auth/oauth?type=google"

	conf = &oauth2.Config{
		ClientID:     cred.Cid,
		ClientSecret: cred.Csecret,
		RedirectURL:  ru,
		Scopes: []string{
			"https://www.googleapis.com/auth/userinfo.email",
			"https://www.googleapis.com/auth/userinfo.profile",
		},
		Endpoint: google.Endpoint,
	}
}

// OAuthURL returns the oauth URL
func OAuthURL(c *gin.Context) {
	var url string = conf.AuthCodeURL(state)

	c.JSON(http.StatusOK, gin.H{"url": url})
}

// OAuthAuthenticator authenticate users via google oauth
func OAuthAuthenticator(c *gin.Context) (interface{}, error) {
	tok, err := conf.Exchange(oauth2.NoContext, c.Query("code"))
	if err != nil {
		return nil, err
	}

	client := conf.Client(oauth2.NoContext, tok)
	userinfo, err := client.Get("https://www.googleapis.com/oauth2/v3/userinfo")
	if err != nil {
		return nil, err
	}

	defer userinfo.Body.Close()
	data, _ := ioutil.ReadAll(userinfo.Body)
	var ou OAuthUser
	if err := json.Unmarshal(data, &ou); err != nil {
		return nil, errors.New(err.Error())
	}

	var u model.User
	config.GetDB().Where("email = ?", ou.Email).First(&u)

	if u.ID > 0 {
		return &u, nil
	}

	u.Email = ou.Email
	u.Firstname = ou.GivenName
	u.Lastname = ou.FamilyName
	u.Active = ou.EmailVerified

	config.GetDB().Save(&u)

	return &u, nil
}

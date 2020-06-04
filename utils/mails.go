package utils

import (
	"log"
	"net/smtp"
	"os"

	"github.com/giuliobosco/todoAPI/config"
	"github.com/giuliobosco/todoAPI/model"
)

var (
	smtpServer   = os.Getenv("SMTP_SERVER")
	smtpPort     = os.Getenv("SMTP_PORT")
	smtpUsername = os.Getenv("SMTP_USERNAME")
	smtpPassword = os.Getenv("SMTP_PASSWORD")
)

func UserConfirmationSendMail(user *model.User) {
	auth := smtp.PlainAuth("", smtpUsername, smtpPassword, smtpServer)

	// Here we do it all: connect to our server, set up a message and send it
	to := []string{user.Email}
	msg := config.BuildConfirmEmail(user, smtpUsername)
	err := smtp.SendMail(smtpServer+":"+smtpPort, auth, user.Email, to, msg)
	if err != nil {
		log.Fatal(err)
	}
}

func UserPasswordRecoverySendMail(user *model.User) {
	auth := smtp.PlainAuth("", smtpUsername, smtpPassword, smtpServer)

	// Here we do it all: connect to our server, set up a message and send it
	to := []string{user.Email}
	msg := config.BuildPasswordRecovery(user, smtpUsername)
	err := smtp.SendMail(smtpServer+":"+smtpPort, auth, user.Email, to, msg)
	if err != nil {
		log.Fatal(err)
	}
}

package main

import (
	"fmt"
	"io/ioutil"
	"log"
	"os"
	"strconv"
	"strings"
	"time"

	"github.com/calvarado2004/bookings/internal/models"
	mail "github.com/xhit/go-simple-mail/v2"
)

//listnForMail listens for the mail channel

func listenForMail() {

	//anonymous function listening to the channel
	go func() {
		for {
			msg := <-app.MailChan
			sendMsg(msg)
		}
	}()

}

func sendMsg(m models.MailData) {

	portNumber, err := strconv.Atoi(os.Getenv("MAILHOG_PORT"))
	if err != nil {
		errorLog.Println(err)
	}

	server := mail.NewSMTPClient()
	server.Host = os.Getenv("MAILHOG_HOST")
	server.Port = portNumber
	server.KeepAlive = false
	server.ConnectTimeout = 10 * time.Second
	server.SendTimeout = 10 * time.Second

	client, err := server.Connect()
	if err != nil {
		errorLog.Println(err)
	}

	email := mail.NewMSG()
	email.SetFrom(m.From).AddTo(m.To).SetSubject(m.Subject)
	if m.Template == "" {
		email.SetBody(mail.TextHTML, m.Content)

	} else {
		data, err := ioutil.ReadFile(fmt.Sprintf("./email-templates/%s", m.Template))
		if err != nil {
			app.ErrorLog.Println(err)
		}

		mailTemplate := string(data)

		msgToSend := strings.Replace(mailTemplate, "[%body%]", m.Content, 1)

		email.SetBody(mail.TextHTML, msgToSend)
	}

	err = email.Send(client)
	if err != nil {
		errorLog.Println(err)
	} else {
		log.Println("Email sent!")
	}

}

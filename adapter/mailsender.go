package adapter

import (
	"encoding/base64"
	"fmt"
	"github.com/guru-invest/guru.api.preregistry/configuration"
	"github.com/scorredoira/email"
	"io/ioutil"
	"log"
	"net/http"
	"net/mail"
	"net/smtp"
	"strconv"
)

type MailData struct {
	Subject string
	Body    string
	IsHTML  bool
	From    string
	To      []string
	CC      []string
	Bcc     []string
	Attach  []string
}

type MailSender struct {
	smtp.Auth
	Host string
	Port int
}

func NewMailSender() *MailSender {
	pArray, _ := base64.StdEncoding.DecodeString(configuration.CONFIGURATION.MAIL.SMTPPassword)
	smtpPassword := string(pArray)
	m := &MailSender{}
	m.Auth = smtp.PlainAuth("", configuration.CONFIGURATION.MAIL.SMTPUser, smtpPassword, configuration.CONFIGURATION.MAIL.SMTPServer)
	m.Host = configuration.CONFIGURATION.MAIL.SMTPServer
	port,_ := strconv.Atoi(configuration.CONFIGURATION.MAIL.SMTPPort)
	m.Port = port

	return m
}

func (m *MailSender) GetAddress() string {
	return fmt.Sprintf("%s:%d", m.Host, m.Port)
}

func (m *MailSender) CreateMessage(data MailData) *email.Message {
	if data.IsHTML {
		return email.NewHTMLMessage(data.Subject, data.Body)
	}

	return email.NewMessage(data.Subject, data.Body)
}

func (m *MailSender) Send(mailData MailData) error {
	item := m.CreateMessage(mailData)
	item.From = mail.Address{Name: "Guru", Address: configuration.CONFIGURATION.MAIL.SMTPUser}
	item.To = mailData.To
	item.Cc = mailData.CC
	item.Bcc = mailData.Bcc
	mailData.Attach = []string{"",""}

	var err error
	if mailData.Attach[0] != "" {
		for _, att := range mailData.Attach {
			err = item.AttachBuffer("Anexo.pdf", GetAttachment(att), false)
			if err != nil {
				log.Println(err)
			}
		}
	}
	err = email.Send(m.GetAddress(), m.Auth, item)

	if err != nil {
		println(err.Error())
	}

	return err
}

func GetAttachment(link string) []byte {
	resp, err := http.Get(link)
	if err != nil {
		//
	}

	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		//
	}

	return body
}
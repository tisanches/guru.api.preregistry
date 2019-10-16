package adapter

import (
	"fmt"
	"github.com/guru-invest/guru.api.preregistry/configuration"
	"github.com/mailjet/mailjet-apiv3-go"
	"log"
	"math/rand"
	"strings"
)

type EmailWorkflow struct{
	EmailBody string
	Template string
	To string
	From string
	Name string
}

func (e *EmailWorkflow) BuildWelComeEmail(link string, template []byte) {
	workflowByte := template
	name := strings.Split(e.Name, " ")
	strHtml := string(workflowByte)
	strHtml = strings.Replace(strHtml, "{User}", name[0], 1)
	strHtml = strings.Replace(strHtml, "{Link}", link, 1)
	e.EmailBody = strHtml
}


func (e *EmailWorkflow) SendEmail(subject string) {
	mailData := parseMailData(e, subject)
	sendEmail(mailData, e.To, e.Name)
}

func parseMailData(e *EmailWorkflow, subject string) MailData {
	mailData := MailData{}
	to := strings.Split(e.To, ",")
	for i := range to {
		if to[i] != "" {
			mailData.To = append(mailData.To, strings.TrimSpace(to[i]))
		}
	}
	mailData.From = e.From
	mailData.Subject = subject
	mailData.IsHTML = true
	mailData.Body = e.EmailBody
	return mailData
}

func sendEmail(mailData MailData, To string, customerName string) {
	//mailSender := NewMailSender()
	//mailSender.Send(mailData)

	mailjetClient := mailjet.NewMailjetClient(configuration.CONFIGURATION.MAIL.MailjetApiKeyPublic, configuration.CONFIGURATION.MAIL.MailjetApiKeyPrivate)
	messagesInfo := []mailjet.InfoMessagesV31 {
		{
			From: &mailjet.RecipientV31{
				Email: configuration.CONFIGURATION.MAIL.MailjetUsername,
				Name:  configuration.CONFIGURATION.MAIL.MailjetName,
			},
			To: &mailjet.RecipientsV31{
				mailjet.RecipientV31{
					Email: To,
					Name:  customerName,
				},
			},
			Subject:  mailData.Subject,
			HTMLPart: mailData.Body,
			CustomID: generateGUID(),
		},
	}
	messages := mailjet.MessagesV31{Info: messagesInfo }
	res, err := mailjetClient.SendMailV31(&messages)
	if err != nil {
		log.Fatal(err)
	}
	fmt.Printf("Data: %+v\n", res)
}

func generateGUID() string{
	b := make([]byte, 16)
	_, err := rand.Read(b)
	if err != nil {
		log.Fatal(err)
	}
	uuid := fmt.Sprintf("%x-%x-%x-%x-%x",
		b[0:4], b[4:6], b[6:8], b[8:10], b[10:])
	return uuid
}
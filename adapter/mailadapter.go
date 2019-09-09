package adapter

import (
	"io/ioutil"
	"strings"
)

type EmailWorkflow struct{
	EmailBody string
	Template string
	To string
	From string
	Name string
}

func (e *EmailWorkflow) BuildWelComeEmail() {
	workflowByte, _ := ioutil.ReadFile(e.Template)
	name := strings.Split(e.Name, " ")
	strHtml := string(workflowByte)
	strHtml = strings.Replace(strHtml, "{User}", name[0], 1)
	e.EmailBody = strHtml
}


func (e *EmailWorkflow) SendEmail(subject string) {
	mailData := parseMailData(e, subject)
	sendEmail(mailData)
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

func sendEmail(mailData MailData) {
	mailSender := NewMailSender()
	mailSender.Send(mailData)
}

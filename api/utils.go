package api

import (
	"encoding/json"
	"github.com/gin-gonic/gin"
	"github.com/guru-invest/guru.api.preregistry/adapter"
	"github.com/guru-invest/guru.api.preregistry/configuration"
	"github.com/guru-invest/guru.api.preregistry/domain"
	"io/ioutil"
	"log"
	"net/http"
	"regexp"
	"strings"
)

type authenticationType int

const (
	email authenticationType = iota
	document_number
	contact
	unknow
)

func validate(authentication string) authenticationType{
	if validateEmail(authentication){
		return email
	}else if validateDocument(authentication){
		return document_number
	}else if validateContact(authentication){
		return contact
	}else{
		return unknow
	}
}

func validateEmail(email string) bool{
	rEmail := regexp.MustCompile("^[a-zA-Z0-9.!#$%&'*+/=?^_`{|}~-]+@[a-zA-Z0-9](?:[a-zA-Z0-9-]{0,61}[a-zA-Z0-9])?(?:\\.[a-zA-Z0-9](?:[a-zA-Z0-9-]{0,61}[a-zA-Z0-9])?)*$")
	res := rEmail.MatchString(email)
	return res
}

func validateDocument(document_number string)bool{
	rDocument := regexp.MustCompile("^[0-9]*$")
	res := rDocument.MatchString(document_number)
	if res && strings.Count(document_number, "") >= 12 {
		return res
	}else{
		return false
	}
}

func validateContact(contact string)bool{
	rContact := regexp.MustCompile(`^(?:(?:\(?(?:00|\+)([1-4]\d\d|[1-9]\d?)\)?)?[\-\.\ \\\/]?)?((?:\(?\d{1,}\)?[\-\.\ \\\/]?){0,})(?:[\-\.\ \\\/]?(?:#|ext\.?|extension|x)[\-\.\ \\\/]?(\d+))?$`)
	res := rContact.MatchString(contact)
	if res && strings.Count(contact, "") <= 12 {
		return res
	}else{
		return false
	}
}

func getAuthentication(email string) map[string]interface{}{
	client := &http.Client{}
	req, _ := http.NewRequest("GET", configuration.CONFIGURATION.OTHER.Authentication + email, nil)
	res, err := client.Do(req)
	if err != nil {
		log.Println(err)
	}else{
		if res.Status == "200 OK"{
			reqBody, err := ioutil.ReadAll(res.Body)
			if err != nil{
				return make(map[string]interface{})
			}
			resp := make(map[string]interface{})
			err = json.Unmarshal(reqBody, &resp)
			if err != nil{
				return make(map[string]interface{})
			}
			return resp
		}
	}
	return make(map[string]interface{})
}

func getAuthorization(token string) bool{
	client := &http.Client{}
	req, _ := http.NewRequest("GET", configuration.CONFIGURATION.OTHER.Authorization, nil)
	req.Header.Set("Authorization", "bearer " + token)
	res, err := client.Do(req)
	if err != nil {
		log.Println(err)
	}else{
		if res.Status == "200 OK"{
			return true
		}
	}
	return false
}


type mailType int

const (
	authorization mailType = iota
	welcome
)

func sendEmail(to string, name string, link string, mtype mailType){
	subject := ""
	mail := adapter.EmailWorkflow{}
	switch mtype {
	case authorization:
		mail.Template = "templates/authorizeTemplate.html"
		subject = "Autorização de login"
	default:
		mail.Template = "templates/welcomeTemplate.html"
		subject = "Bem-vindo ao Guru"
	}
	mail.To = to
	mail.Name = name
	mail.From = configuration.CONFIGURATION.MAIL.SMTPUser
	mail.BuildWelComeEmail()
	mail.SendEmail(subject)
}

func sendCredentials(customer_code string, c *gin.Context){
	position := domain.Position{}
	position.Get(customer_code)
	if position.Customer_Code != ""{
		msg := make(map[string]interface{})
		m := getAuthentication(position.Email)
		msg["customer_code"] = position.Customer_Code
		msg["token"] = m["token"].(string)
		c.AbortWithStatusJSON(200, msg)
	}else{
		msg := make(map[string]interface{})
		msg["error"] = "User not foud"
		c.AbortWithStatusJSON(404, msg)
	}
}

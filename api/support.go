package api

import (
	"bytes"
	"encoding/json"
	"errors"
	"github.com/gin-gonic/gin"
	"github.com/guru-invest/guru.api.preregistry/adapter"
	"github.com/guru-invest/guru.api.preregistry/configuration"
	"github.com/guru-invest/guru.api.preregistry/domain"
	"github.com/guru-invest/guru.api.preregistry/logger"
	"github.com/guru-invest/guru.framework/api"
	"io/ioutil"
	"log"
	"net/http"
	"regexp"
	"strconv"
	"strings"
)

type authenticationType int

const (
	email authenticationType = iota
	document_number
	contact
	unknow
)

// Esperando o Antoine pedir pra liberar a api de validação de dados =)
func validate(authentication string) authenticationType {
	if validateEmail(authentication) {
		return email
	} else if validateDocument(authentication) {
		return document_number
	} else if validateContact(authentication) {
		return contact
	} else {
		return unknow
	}
}

func validateEmail(email string) bool {
	rEmail := regexp.MustCompile("^[a-zA-Z0-9.!#$%&'*+/=?^_`{|}~-]+@[a-zA-Z0-9](?:[a-zA-Z0-9-]{0,61}[a-zA-Z0-9])?(?:\\.[a-zA-Z0-9](?:[a-zA-Z0-9-]{0,61}[a-zA-Z0-9])?)*$")
	res := rEmail.MatchString(email)
	return res
}

func validateDocument(document_number string) bool {
	rDocument := regexp.MustCompile("^[0-9]*$")
	res := rDocument.MatchString(document_number)
	if res && strings.Count(document_number, "") >= 12 {
		return res
	} else {
		return false
	}
}

func validateContact(contact string) bool {
	rContact := regexp.MustCompile(`^(?:(?:\(?(?:00|\+)([1-4]\d\d|[1-9]\d?)\)?)?[\-\.\ \\\/]?)?((?:\(?\d{1,}\)?[\-\.\ \\\/]?){0,})(?:[\-\.\ \\\/]?(?:#|ext\.?|extension|x)[\-\.\ \\\/]?(\d+))?$`)
	res := rContact.MatchString(contact)
	if res && strings.Count(contact, "") <= 12 {
		return res
	} else {
		return false
	}
}

func getAuthentication(email string) map[string]interface{} {
	client := &http.Client{}
	req, _ := http.NewRequest("GET", configuration.CONFIGURATION.OTHER.Authentication+email, nil)
	res, err := client.Do(req)
	if err != nil {
		log.Println(err)
	} else {
		if res.Status == "200 OK" {
			reqBody, err := ioutil.ReadAll(res.Body)
			if err != nil {
				return make(map[string]interface{})
			}
			resp := make(map[string]interface{})
			err = json.Unmarshal(reqBody, &resp)
			if err != nil {
				return make(map[string]interface{})
			}
			return resp
		}
	}
	return make(map[string]interface{})
}

func getAuthorization(token string) bool {
	client := &http.Client{}
	req, _ := http.NewRequest("GET", configuration.CONFIGURATION.OTHER.Authorization, nil)
	req.Header.Set("Authorization", "bearer "+token)
	res, err := client.Do(req)
	if err != nil {
		log.Println(err)
	} else {
		if res.Status == "200 OK" {
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

func sendEmail(to string, name string, link string, mtype mailType) {
	subject := ""
	mail := adapter.EmailWorkflow{}
	reqBody := []byte{}
	switch mtype {
	case authorization:
		client := &http.Client{}
		req, _ := http.NewRequest("GET", "https://guruimages.s3.us-east-2.amazonaws.com/authorizeTemplate.html", nil)
		res, err := client.Do(req)
		if err != nil {
			log.Println(err)
		} else {
			if res.Status == "200 OK" {
				reqBody, err = ioutil.ReadAll(res.Body)
				if err != nil {
					log.Println(err)
				}

			}
		}

		subject = "Autorização de login"
	default:
		client := &http.Client{}
		req, _ := http.NewRequest("GET", "https://guruimages.s3.us-east-2.amazonaws.com/welcomeTemplate.html", nil)
		res, err := client.Do(req)
		if err != nil {
			log.Println(err)
		} else {
			if res.Status == "200 OK" {
				reqBody, err = ioutil.ReadAll(res.Body)
				if err != nil {
					log.Println(err)
				}

			}
		}
		subject = "Bem-vindo ao Guru"
	}
	mail.To = to
	mail.Name = name
	mail.From = configuration.CONFIGURATION.MAIL.SMTPUser
	mail.BuildWelComeEmail(link, reqBody)
	mail.SendEmail(subject)
}

func sendCredentials(customer_code string, c *gin.Context) {
	position := domain.Position{}
	err := position.Get(customer_code)
	checkErr(err, c)
	if position.Customer_Code != "" {
		msg := make(map[string]interface{})
		m := getAuthentication(position.Email)
		msg["customer_code"] = position.Customer_Code
		msg["email"] = position.Email
		msg["token"] = m["token"].(string)
		c.AbortWithStatusJSON(200, msg)
	} else {
		msg := make(map[string]interface{})
		msg["error"] = "User not foud"
		c.AbortWithStatusJSON(404, msg)
	}
}

func checkErr(err error, c *gin.Context) bool {
	if err != nil {
		logger.LOG.Error("error on executing. stack: " + err.Error())
		return true
	}
	return false
}

func insertCustomer(customer domain.Customer, c *gin.Context) {
	referral_position := domain.Position{}
	referral_position.Get(getCustomerByReferralCode(customer.Referral_Code))
	oldposition := referral_position.Position
	err := customer.Insert()
	if checkErr(err, c) {
		api.Error400(errors.New("invalid customer."), c)
	} else {
		position := domain.Position{}
		err = position.Get(customer.Customer_Code)
		if checkErr(err, c) {
			api.Error400(errors.New("invalid customer."), c)
		} else {
			if position.DocumentNumber == "" {
				msg := make(map[string]interface{})
				msg["msg"] = "Step saved."
				c.AbortWithStatusJSON(200, msg)
			} else {
				buildNotification(customer, referral_position, oldposition)
				sendEmail(position.Email, position.Name, "", welcome)
				sendCredentials(customer.Customer_Code, c)
			}
		}
	}
}

func buildNotification(customer domain.Customer, referral_position domain.Position, oldposition int64) {
	if customer.Referral_Code != "" {
		referral_position.Get(getCustomerByReferralCode(customer.Referral_Code))
		newPosition := referral_position.Position
		originCustomer := getCustomerByReferralCode(customer.Referral_Code)
		res := oldposition - newPosition
		sendNotification(originCustomer, customer.Email, customer.Referral_Code, strconv.Itoa(int(res)), customer.Name)
	}
}

func insertCustomerLanding(customer domain.Customer, c *gin.Context) {
	err := customer.Insert()
	if checkErr(err, c) {
		api.Error400(errors.New("invalid customer."), c)
	} else {
		if customer.Email != "" || customer.Referral_Code != "" &&
			customer.Name == "" && customer.DocumentNumber == "" &&
			customer.Contact == "" && customer.Customer_Code == "" {
			sendEmail(customer.Email, customer.Name, "", welcome)
		}
		msg := make(map[string]interface{})
		msg["msg"] = "Step saved."
		c.AbortWithStatusJSON(200, msg)
	}

}

func getCustomerByReferralCode(referral string) string {
	ref := domain.Referrals{}
	ref.Get(referral)
	return ref.Origin_Code
}

func sendNotification(customer_code string, email string, referral_code string, position string, friendname string) {
	m := make(map[string]interface{})
	m["customer_codes"] = []string{customer_code}
	m["title"] = configuration.CONFIGURATION.MESSAGES.NewReferralTitle
	message := strings.Replace(configuration.CONFIGURATION.MESSAGES.NewReferalMessage, "{email}", email, 1)
	m["message"] = message
	m["deeplink"] = configuration.CONFIGURATION.MESSAGES.NewReferralDeeplink + referral_code + "?position=" + position + "&friendname=" + friendname
	bytesRepresentation, _ := json.Marshal(m)
	resp, err := http.Post(configuration.CONFIGURATION.OTHER.Notification, "application/json", bytes.NewBuffer(bytesRepresentation))
	if err != nil {
		println(err)
	}
	println(resp.Status)
}

func updateCustomer(customer domain.Customer, c *gin.Context) {
	err := customer.Update()
	if checkErr(err, c) {
		api.Error400(errors.New("invalid customer."), c)
	} else {
		msg := make(map[string]interface{})
		msg["msg"] = "customer updated."
		c.AbortWithStatusJSON(200, msg)
	}
}

func treatCustomer(customer domain.Customer, ePosition domain.Position, c *gin.Context) {
	if validateEmail(customer.Email) {
		if ePosition.Customer_Code != "" {
			customer.Customer_Code = ePosition.Customer_Code
			if customer.Contact != "" {
				updateCustomer(customer, c)
			} else {
				//api.Error400(errors.New("invalid customer."), c)
				insertCustomer(customer, c)
			}
		} else {
			sCustomer := customer
			if sCustomer.DocumentNumber != "" {
				sCustomer.GetByEmail(sCustomer.Email)
				if ((sCustomer.DocumentNumber != customer.DocumentNumber) || (sCustomer.Email != customer.Email)) &&
					(sCustomer.DocumentNumber != "") {
					api.Error400(errors.New("user already exists."), c)
				} else {
					insertCustomer(customer, c)
				}
			} else {
				insertCustomer(sCustomer, c)
			}
		}
	}else{
		api.Error400(errors.New("invalid customer."), c)
	}
}

func treatCustomerLanding(customer domain.Customer, c *gin.Context) {
	sPosition := domain.Position{}
	if validateEmail(customer.Email) {
		err := sPosition.GetByEmail(customer.Email)
		if err == nil {
			if sPosition.Email == customer.Email {
				api.Error400(errors.New("user already exists."), c)
			} else {
				insertCustomerLanding(customer, c)
			}
		} else {

			insertCustomerLanding(customer, c)
		}
	} else {
		api.Error400(errors.New("invalid customer"), c)
	}
}

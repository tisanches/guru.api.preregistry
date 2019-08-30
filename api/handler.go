package api

import (
	"encoding/json"
	"github.com/gin-gonic/gin"
	"github.com/guru-invest/guru.api.preregistry/configuration"
	"github.com/guru-invest/guru.api.preregistry/domain"
	"github.com/guru-invest/guru.framework/api"
	"io/ioutil"
)

func InitializeApi(){
	createRoutes()
	api.InitRoutering(configuration.CONFIGURATION.API.Port, "v1", false)
}

func createRoutes(){
	api.AddRoute(api.POST, configuration.CONFIGURATION.API.Route + "/new", createCustomer)
	api.AddRoute(api.GET, configuration.CONFIGURATION.API.Route + "/position/:document_number", getPosition)
	api.AddRoute(api.GET, configuration.CONFIGURATION.API.Route + "/referrals/:referral_code", getReferrals)
}

func validate(c *gin.Context, field string){
	if field == ""{
		c.AbortWithStatusJSON(400, "Missing key: " + field)
	}
}

func  createCustomer(c *gin.Context){
	customer := domain.Customer{}
	reqBody, err := ioutil.ReadAll(c.Request.Body)
	if err != nil{
		c.AbortWithStatusJSON(400, "Invalid format")
	}
	err = json.Unmarshal(reqBody, &customer)
	if err != nil{
		c.AbortWithStatusJSON(400, "Invalid format")
	}
	validate(c, customer.Name)
	validate(c, customer.Contact)
	validate(c, customer.Email)
	validate(c, customer.DocumentNumber)
	validate(c, customer.Password)
	customer.Insert()
	position := domain.Position{}
	position.Get(customer.DocumentNumber)
	c.AbortWithStatusJSON(200, position)
}

func getPosition(c *gin.Context){
	document_number := c.Param("document_number")
	if document_number == ""{
		c.AbortWithStatusJSON(400, "Missing key: document_number")
	}else{
		position := domain.Position{}
		position.Get(document_number)
		c.AbortWithStatusJSON(200, position)
	}
}

func getReferrals(c *gin.Context){
	referral_code := c.Param("referral_code")
	if referral_code == ""{
		c.AbortWithStatusJSON(400, "Missing key: referral_code")
	}else{
		referrals := domain.Referrals{}
		referrals.Get(referral_code)
		c.AbortWithStatusJSON(200, referrals)
	}
}
package api

import (
	"encoding/json"
	"github.com/gin-gonic/gin"
	"github.com/guru-invest/guru.api.preregistry/configuration"
	"github.com/guru-invest/guru.api.preregistry/domain"
	"github.com/guru-invest/guru.framework/api"
	"io/ioutil"
	"strings"
)

func InitializeApi(){
	createRoutes()
	api.InitRoutering(configuration.CONFIGURATION.API.Port, "v1", true)
}

func createRoutes(){
	api.AddRoute(api.POST, configuration.CONFIGURATION.API.Route + "/new", createCustomer)
	api.AddRoute(api.GET, configuration.CONFIGURATION.API.Route + "/position/:customer_code", getPosition)
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
	position.Get(customer.Customer_Code)
	c.AbortWithStatusJSON(200, position)
}

func getPosition(c *gin.Context){
	customer_code := c.Param("customer_code")
	if customer_code == ""{
		c.AbortWithStatusJSON(400, "Missing key: customer_code")
	}else{
		position := domain.Position{}
		position.Get(customer_code)
		c.AbortWithStatusJSON(200, position)
	}
}

func getReferrals(c *gin.Context){
	referral_code := c.Param("referral_code")
	if referral_code == ""{
		c.AbortWithStatusJSON(400, "Missing key: referral_code")
	}else{
		if strings.Contains("https://seja.guru/", referral_code){
			referral_code = strings.Replace(referral_code, "https://seja.guru/", "",1 )
		}
		referrals := domain.Referrals{}
		referrals.Get(referral_code)
		c.AbortWithStatusJSON(200, referrals)
	}
}
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
	api.AddRoute(api.GET, configuration.CONFIGURATION.API.Route + "/customer/:param", getCustomer)
	api.AddRoute(api.POST, configuration.CONFIGURATION.API.Route + "/authorize/device", setDeviceAuthorization)
	api.AddRoute(api.GET, configuration.CONFIGURATION.API.Route + "/position/:customer_code", getPosition)
	api.AddRoute(api.GET, configuration.CONFIGURATION.API.Route + "/authentication/:customer_code", getAuthenticationHandler)
	api.AddRoute(api.GET, configuration.CONFIGURATION.API.Route + "/authorize/:token/:customer_code", getPositionWithToken)
	api.AddRoute(api.GET, configuration.CONFIGURATION.API.Route + "/referrals/:referral_code", getReferrals)
}

func  createCustomer(c *gin.Context){
	customer := domain.Customer{}
	reqBody, err := ioutil.ReadAll(c.Request.Body)
	if err != nil{
		msg := make(map[string]interface{})
		msg["error"] = "Invalid format"
		c.AbortWithStatusJSON(400, msg)
	}
	err = json.Unmarshal(reqBody, &customer)
	if err != nil{
		msg := make(map[string]interface{})
		msg["error"] = "Invalid format"
		c.AbortWithStatusJSON(400, msg)
	}
	customer.Insert()
	position := domain.Position{}
	position.Get(customer.Customer_Code)
	if position.Customer_Code == ""{
		msg := make(map[string]interface{})
		msg["msg"] = "Step saved."
		c.AbortWithStatusJSON(200, msg)
	}else {
		sendEmail(position.Email, position.Name, "", welcome)
		msg := make(map[string]interface{})
		msg["msg"] = "User notified"
		msg["customer_code"] = position.Customer_Code
		c.AbortWithStatusJSON(200, msg)
	}
}

func getPosition(c *gin.Context){
	token := c.GetHeader("Authorization")
	token = strings.Replace(token, "Bearer ", "", 1)
	token = strings.Replace(token, "bearer ", "", 1)
	if token != "" && getAuthorization(token) {
		customer_code := c.Param("customer_code")
		if customer_code == "" {
			msg := make(map[string]interface{})
			msg["error"] = "Missing Key: customer_code"
			c.AbortWithStatusJSON(400, msg)
		} else {
			position := domain.Position{}
			position.Get(customer_code)
			c.AbortWithStatusJSON(200, position)
		}
	}else{
		msg := make(map[string]interface{})
		msg["error"] = "Invalid Token"
		c.AbortWithStatusJSON(400, msg)
	}
}

func getCustomer(c *gin.Context){
	param := c.Param("param")
	if param == ""{
		msg := make(map[string]interface{})
		msg["error"] = "Missing Key: email"
		c.AbortWithStatusJSON(400, msg)
	}else{
		customer := domain.Customer{}
		position := domain.Position{}
		customer.GetByEmail(param)
		position.GetByEmail(param)
		if position.Customer_Code != ""{
			msg := make(map[string]interface{})
			msg["customer_code"] = position.Customer_Code
			c.AbortWithStatusJSON(200, msg)
		}else if customer.Email != ""{
			c.AbortWithStatusJSON(200, customer)
		}else{
			msg := make(map[string]interface{})
			msg["error"] = "User not foud"
			c.AbortWithStatusJSON(404, msg)
		}
	}
}

func setDeviceAuthorization(c *gin.Context){
	customer := domain.Customer{}
	reqBody, err := ioutil.ReadAll(c.Request.Body)
	if err != nil{
		msg := make(map[string]interface{})
		msg["error"] = "Invalid format"
		c.AbortWithStatusJSON(400, msg)
	}
	err = json.Unmarshal(reqBody, &customer)
	if err != nil{
		msg := make(map[string]interface{})
		msg["error"] = "Invalid format"
		c.AbortWithStatusJSON(400, msg)
	}
	if customer.Customer_Code == ""{
		msg := make(map[string]interface{})
		msg["error"] = "Missing Key: customer_code"
		c.AbortWithStatusJSON(400, msg)
	}else {
		position := domain.Position{}
		position.Get(customer.Customer_Code)
		if position.Customer_Code != "" {
			m := getAuthentication(position.Email)
			link := configuration.CONFIGURATION.OTHER.PositionPrefix + m["token"].(string) + "/" + m["customer_code"].(string)
			sendEmail(position.Email, position.Name, link, authorization)
			msg := make(map[string]interface{})
			msg["msg"] = "Email sent to the user."
			c.AbortWithStatusJSON(200, msg)
		}else{
			msg := make(map[string]interface{})
			msg["error"] = "Customer not foud"
			c.AbortWithStatusJSON(404, msg)
		}
	}
}

func getPositionWithToken(c *gin.Context) {
	customer_code := c.Param("customer_code")
	token := c.Param("token")
	if token != "" && getAuthorization(token) {
		if customer_code != "" {
			sendCredentials(customer_code, c)
		} else {
			msg := make(map[string]interface{})
			msg["error"] = "Invalid token"
			c.AbortWithStatusJSON(400, msg)
		}
	} else {
		msg := make(map[string]interface{})
		msg["error"] = "Invalid token"
		c.AbortWithStatusJSON(400, msg)
	}
}

func getReferrals(c *gin.Context){
	referral_code := c.Param("referral_code")
	if referral_code == ""{
		msg := make(map[string]interface{})
		msg["error"] = "Missing Key: referral_code"
		c.AbortWithStatusJSON(400, msg)
	}else{
		if strings.Contains(configuration.CONFIGURATION.OTHER.DeepLinkPrefix, referral_code){
			referral_code = strings.Replace(referral_code, configuration.CONFIGURATION.OTHER.DeepLinkPrefix, "",1 )
		}
		referrals := domain.Referrals{}
		referrals.Get(referral_code)
		if referrals.Referral_Code == ""{
			msg := make(map[string]interface{})
			msg["error"] = "Referral code not found"
			c.AbortWithStatusJSON(404, msg)
		}
		c.AbortWithStatusJSON(200, referrals)
	}
}

func getAuthenticationHandler(c *gin.Context) {
	customer_code := c.Param("customer_code")
	if customer_code == ""{
		msg := make(map[string]interface{})
		msg["error"] = "Missing Key: customer_code"
		c.AbortWithStatusJSON(400, msg)
	}else{
		sendCredentials(customer_code, c)
	}
}



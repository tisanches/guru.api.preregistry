package api

import (
	"encoding/json"
	"github.com/gin-gonic/gin"
	"github.com/guru-invest/guru.api.preregistry/configuration"
	"github.com/guru-invest/guru.api.preregistry/domain"
	"github.com/guru-invest/guru.framework/api"
	"io/ioutil"
	"log"
	"net/http"
	"strings"
)

func InitializeApi(){
	createRoutes()
	api.InitRoutering(configuration.CONFIGURATION.API.Port, "v1", false)
}

func createRoutes(){
	api.AddRoute(api.POST, configuration.CONFIGURATION.API.Route + "/new", createCustomer)
	api.AddRoute(api.GET, configuration.CONFIGURATION.API.Route + "/customer/:email", getPreRegistryStep)
	//api.AddRoute(api.GET, configuration.CONFIGURATION.API.Route + "/position/:customer_code", getPosition)
	api.AddRoute(api.GET, configuration.CONFIGURATION.API.Route + "/position/:token/:customer_code", getPositionWithToken)
	api.AddRoute(api.GET, configuration.CONFIGURATION.API.Route + "/referrals/:referral_code", getReferrals)
}

func validate(c *gin.Context, field string){
	if field == ""{
		msg := make(map[string]interface{})
		msg["error"] = "Missing key: " + field
		c.AbortWithStatusJSON(400, msg)
	}
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
		c.AbortWithStatusJSON(200, position)
	}
}

func getPosition(c *gin.Context){
	customer_code := c.Param("customer_code")
	if customer_code == ""{
		msg := make(map[string]interface{})
		msg["error"] = "Missing Key: customer_code"
		c.AbortWithStatusJSON(400, msg)
	}else{
		position := domain.Position{}
		position.Get(customer_code)
		c.AbortWithStatusJSON(200, position)
	}
}

func getPreRegistryStep(c *gin.Context){
	email := c.Param("email")
	if email == ""{
		msg := make(map[string]interface{})
		msg["error"] = "Missing Key: email"
		c.AbortWithStatusJSON(400, msg)
	}else{
		customer := domain.Customer{}
		position := domain.Position{}
		customer.GetByEmail(email)
		position.GetByEmail(email)
		if position.Customer_Code != ""{
			msg := make(map[string]interface{})
			msg["customer_code"] = position.Customer_Code
			msg["position_at"] = configuration.CONFIGURATION.OTHER.PositionPrefix + getAuthentication(email)
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

func getPositionWithToken(c *gin.Context) {
	customer_code := c.Param("customer_code")
	token := c.Param("token")
	if token != "" && getAuthorization(token) {
		if customer_code != "" {
			position := domain.Position{}
			position.Get(customer_code)
			c.AbortWithStatusJSON(200, position)
		} else {
			msg := make(map[string]interface{})
			msg["error"] = "Missing Key: customer_code"
			c.AbortWithStatusJSON(400, msg)
		}
	} else {
		msg := make(map[string]interface{})
		msg["error"] = "Invalid Token"
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
		c.AbortWithStatusJSON(200, referrals)
	}
}

func getAuthentication(email string) string{
	client := &http.Client{}
	req, _ := http.NewRequest("GET", configuration.CONFIGURATION.OTHER.Authentication + email, nil)
	res, err := client.Do(req)
	if err != nil {
		log.Println(err)
	}else{
		if res.Status == "200 OK"{
			reqBody, err := ioutil.ReadAll(res.Body)
			if err != nil{
				return ""
			}
			resp := make(map[string]interface{})
			err = json.Unmarshal(reqBody, &resp)
			if err != nil{
				return ""
			}
			return resp["token"].(string) + "/" + resp["customer_code"].(string)
		}
	}
	return ""
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
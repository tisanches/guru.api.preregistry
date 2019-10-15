package api

import (
	"encoding/json"
	"github.com/gin-gonic/gin"
	"github.com/guru-invest/guru.api.preregistry/configuration"
	"github.com/guru-invest/guru.api.preregistry/domain"
	"github.com/guru-invest/guru.api.preregistry/logger"
	"github.com/guru-invest/guru.framework/api"
	"github.com/pkg/errors"
	"strings"
)

func InitializeApi(){
	logger.LOG.Info("Creating routes")
	createRoutes()
	logger.LOG.Info("Initializing application server")
	api.InitRoutering(configuration.CONFIGURATION.API.Port, "v1", true)
}

func createRoutes(){
	//region /add route
	logger.LOG.Debug("Adding route " + configuration.CONFIGURATION.API.Route +
		"/add")
	api.AddRoute(api.POST, configuration.CONFIGURATION.API.Route +
		"/add", createCustomerHandler)
	//endregion
	//region /customer/:param route
	logger.LOG.Debug("Adding route " + configuration.CONFIGURATION.API.Route +
		"/customer/:param")
	api.AddRoute(api.GET, configuration.CONFIGURATION.API.Route +
		"/customer/:param", getCustomerHandler)
	//endregion
	//region /authorize/device route
	logger.LOG.Debug("Adding route " + configuration.CONFIGURATION.API.Route +
		"/authorize/devicee")
	api.AddRoute(api.POST, configuration.CONFIGURATION.API.Route + "/authorize/device", setDeviceAuthorizationHandler)
	//endregion
	//region /position/:customer_code route
	logger.LOG.Debug("Adding route " + configuration.CONFIGURATION.API.Route +
		"/position/:customer_code")
	api.AddRoute(api.GET, configuration.CONFIGURATION.API.Route +
		"/position/:customer_code", getPositionHandler)
	//endregion
	//region /authentication/:customer_code route
	logger.LOG.Debug("Adding route " + configuration.CONFIGURATION.API.Route +
		"/authentication/:customer_code")
	api.AddRoute(api.GET, configuration.CONFIGURATION.API.Route +
		"/authentication/:customer_code", getAuthenticationHandler)
	//endregion
	//region /authorize/a/:token/:customer_code route
	logger.LOG.Debug("Adding route " + configuration.CONFIGURATION.API.Route +
		"/authorize/a/:token/:customer_code")
	api.AddRoute(api.POST, configuration.CONFIGURATION.API.Route +
		"/authorize/a/:token/:customer_code", getAuthorizationByEmailHandler)
	//endregion
	//region /referrals/:referral_code route
	logger.LOG.Debug("Adding route " + configuration.CONFIGURATION.API.Route +
		"/referrals/:referral_code")
	api.AddRoute(api.GET, configuration.CONFIGURATION.API.Route +
		"/referrals/:referral_code", getReferralsHandler)
	//endregion
}

func  createCustomerHandler(c *gin.Context) {
	customer := domain.Customer{}
	m := api.Extract(customer, c)
	err := json.Unmarshal(m, &customer)
	if err != nil {
		api.Error400(err, c)
	}
	ePosition := domain.Position{}
	err = ePosition.GetByEmail(customer.Email)
	if checkErr(err, c) {
		api.Error400(errors.New("invalid customer."), c)
	} else {
		treatCustomer(customer, ePosition, c)
	}
}



func getPositionHandler(c *gin.Context){
	token := c.GetHeader("Authorization")
	token = strings.Replace(token, "Bearer ", "", 1)
	token = strings.Replace(token, "bearer ", "", 1)
	if token != "" && getAuthorization(token) {
		customer_code := c.Param("customer_code")
		if customer_code == "" {
			api.Error400(nil, c)
		} else {
			position := domain.Position{}
			err := position.Get(customer_code)
			checkErr(err, c)
			if position.DocumentNumber != ""{
				c.AbortWithStatusJSON(200, position)
			}else{
				api.Error404(errors.New("customer_code not found"), c)
			}

		}
	}else{
		api.Error400(errors.New("authentication error"), c)
	}
}

func getCustomerHandler(c *gin.Context){
	param := c.Param("param")
	if param == ""{
		api.Error400(errors.New("missing key: email"), c)
	}else{
		customer := domain.Customer{}
		position := domain.Position{}
		err := customer.GetByEmail(param)
		checkErr(err, c)
		err = position.GetByEmail(param)
		checkErr(err, c)
		if position.Customer_Code != ""{
			msg := make(map[string]interface{})
			msg["customer_code"] = position.Customer_Code
			c.AbortWithStatusJSON(200, msg)
		}else if customer.Email != ""{
			c.AbortWithStatusJSON(200, customer)
		}else{
			api.Error404(errors.New("customer not found"), c)
		}
	}
}

func setDeviceAuthorizationHandler(c *gin.Context){
	customer := domain.Customer{}
	m := api.Extract(customer, c)
	err := json.Unmarshal(m, &customer)
	if err != nil {
		api.Error400(err, c)
	}
	if customer.Customer_Code == ""{
		api.Error400(errors.New("missing key: customer_code"), c)
	}else {
		position := domain.Position{}
		err := position.Get(customer.Customer_Code)
		checkErr(err, c)
		if position.Customer_Code != "" {
			m := getAuthentication(position.Email)
			link := configuration.CONFIGURATION.OTHER.AuthorizationPrefix + m["token"].(string) + "/" + m["customer_code"].(string)
			sendEmail(position.Email, position.Name, link, authorization)
			msg := make(map[string]interface{})
			msg["msg"] = "Email sent to the user."
			c.AbortWithStatusJSON(200, msg)
		}else{
			api.Error404(errors.New("customer not found"), c)
		}
	}
}

func getAuthorizationByEmailHandler(c *gin.Context) {
	customer_code := c.Param("customer_code")
	token := c.Param("token")
	if token != "" && getAuthorization(token) {
		if customer_code != "" {
			sendCredentials(customer_code, c)
		} else {
			api.Error400(errors.New("missing key: customer_code"), c)
		}
	} else {
		api.Error400(errors.New("invalid format"), c)
	}
}

func getReferralsHandler(c *gin.Context){
	referral_code := c.Param("referral_code")
	if referral_code == ""{
		api.Error400(errors.New("missing key: referral_code"), c)
	}else{
		if strings.Contains(configuration.CONFIGURATION.OTHER.DeepLinkPrefix, referral_code){
			referral_code = strings.Replace(referral_code, configuration.CONFIGURATION.OTHER.DeepLinkPrefix, "",1 )
		}
		referrals := domain.Referrals{}
		err := referrals.Get(referral_code)
		checkErr(err, c)
		if referrals.Referral_Code == ""{
			m := make(map[string]interface{})
			m["msg"] = "No referrals found"
			c.AbortWithStatusJSON(200, m)
		}else {
			c.AbortWithStatusJSON(200, referrals)
		}
	}
}

func getAuthenticationHandler(c *gin.Context) {
	customer_code := c.Param("customer_code")
	if customer_code == ""{
		api.Error400(errors.New("missing key: customer_code"), c)
	}else{
		sendCredentials(customer_code, c)
	}
}



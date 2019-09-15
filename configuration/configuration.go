package configuration

import (
	b "encoding/base64"
	"encoding/json"
	"io/ioutil"
	"log"
	"net/http"
)

type Api struct{
	Port string `json:"port,omitempty"`
	Route string `json:"route-prefix,omitempty"`
}

type Database struct{
	Port string `json:"port,omitempty"`
	Url string `json:"url,omitempty"`
	Username string `json:"username,omitempty"`
	Password string `json:"password,omitempty"`
	Database string `json:"database,omitempty"`
}

type Mail struct{
	SMTPUser string `json:"smtpuser,omitempty"`
	SMTPPassword string `json:"smtppassword,omitempty"`
	SMTPServer string `json:"smtpserver,omitempty"`
	SMTPPort string `json:"smtpport,omitempty"`
}

type Other struct{
	DeepLinkPrefix string `json:"deeplink-prefix,omitempty"`
	Authorization string `json:"authorization,omitempty"`
	Authentication string `json:"authentication,omitempty"`
	PositionPrefix string `json:"position-prefix,omitempty"`
	LogLevel string `json:"loglevel,omitempty"`
}

type Configuration struct{
	API Api `json:"api,omitempty"`
	DATABASE Database `json:"database,omitempty"`
	OTHER Other `json:"other,omitempty"`
	MAIL Mail `json:"mail,omitempty"`
}

const CONFIGURATION_SERVER_KEY  = "aHR0cHM6Ly9jb25maWd1cmF0aW9uLmd1cnUuY29tLnZjL2FwaS92MS9wcm9kL3ByZXJlZ2lzdHJ5"

var CONFIGURATION Configuration

func (c *Configuration) Load(){
	getConfiguration()
	url,_ := b.StdEncoding.DecodeString(c.DATABASE.Url)
	c.DATABASE.Url = string(url)
	passwd,_ := b.StdEncoding.DecodeString(c.DATABASE.Password)
	c.DATABASE.Password = string(passwd)
}

func getConfigurationKey()string{
	sDec, _ := b.StdEncoding.DecodeString(CONFIGURATION_SERVER_KEY)
	return string(sDec)
}

func getConfiguration() {
	resp, err := http.Get(getConfigurationKey())
	if err != nil {
		log.Fatalln(err)
	}
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		log.Fatalln(err)
	}
	json.Unmarshal(body, &CONFIGURATION)
}

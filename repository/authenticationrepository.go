package repository

import (
	"fmt"
	"github.com/guru-invest/guru.framework/security/cripto"
	"log"
)

func insertAuthentication(customer_code string, password string) error{
	connect()
	defer database.Close()
	sttmt := fmt.Sprintf(INSERTAUTHENTICATION, customer_code, cripto.EncodeSHA256([]byte(password), []byte(password)))
	_, err := database.Exec(sttmt)
	if err != nil {
		log.Println("Error on insert customer referral: %v", err)
		return err
	}
	return nil
}

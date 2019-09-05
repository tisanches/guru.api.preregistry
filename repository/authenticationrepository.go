package repository

import (
	"fmt"
	"log"
)

func insertAuthentication(customer_code string, email string) error{
	connect()
	defer database.Close()
	sttmt := fmt.Sprintf(INSERTAUTHENTICATION, customer_code, email)
	_, err := database.Exec(sttmt)
	if err != nil {
		log.Println("Error on insert customer referral: %v", err)
		return err
	}
	return nil
}

package repository

import (
	"fmt"
	"github.com/guru-invest/guru.framework/dynamic"
	"log"
)

func GetPosition(document_number string)map[string][]map[string]interface{}{
	connect()
	defer database.Close()
	rows, err := database.Query(SELECTCUSTOMERQUEUE, document_number)
	if err != nil {
		log.Println("Error on getting customer: %v", err)
	}
	if rows == nil{
		return make(map[string][]map[string]interface{})
	}else {
		return mapResult(rows, "customer")
	}
}

func GetReferrals(referral_code string)map[string][]map[string]interface{}{
	connect()
	defer database.Close()
	rows, err := database.Query(SELECTREFERRALS, referral_code)
	if err != nil {
		log.Println("Error on getting referrals: %v", err)
	}
	if rows == nil{
		return make(map[string][]map[string]interface{})
	}else {
		return mapResult(rows, "referrals")
	}
}

func InsertCustomer(documentNumber string, name string, email string, contact string, customer_code string, referral_code string, password string) error{
	connect()
	defer database.Close()
	sttmt := fmt.Sprintf(INSERTPREREGISTRY, documentNumber, name, email, contact, customer_code, referral_code)
	_, err := database.Exec(sttmt)
	if err != nil {
		log.Println("Error on insert new customer pre-registry: %v", err)
		return err
	}
	referral_code = dynamic.GenerateShortId()
	err = insertReferrals(customer_code, referral_code)
	if err != nil {
		return err
	}
	err = insertOnQueue(customer_code)
	if err != nil {
		return err
	}

	err = insertAuthentication(customer_code, password)
	if err != nil {
		return err
	}

	return nil
}

func insertReferrals(customer_code string, referral_code string) error{
	connect()
	defer database.Close()
	sttmt := fmt.Sprintf(INSERTREFERRAL, customer_code, referral_code)
	_, err := database.Exec(sttmt)
	if err != nil {
		log.Println("Error on insert customer referral: %v", err)
		return err
	}
	return nil
}

func insertOnQueue(customer_code string) error{
	connect()
	defer database.Close()
	sttmt := fmt.Sprintf(INSERTQUEUE, customer_code)
	_, err := database.Exec(sttmt)
	if err != nil {
		log.Println("Error on insert customer on queue: %v", err)
		return err
	}
	return nil
}
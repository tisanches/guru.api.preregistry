package repository

import (
	"fmt"
	"github.com/guru-invest/guru.framework/dynamic"
	"github.com/pkg/errors"
	"log"
	"strings"
)

func GetPosition(customer_code string)(map[string][]map[string]interface{}, error){
	connect()
	defer database.Close()
	rows, err := database.Query(SELECTCUSTOMERQUEUE, customer_code)
	if err != nil {
		log.Println("Error on getting customer: %v", err)
	}
	if rows == nil{
		return make(map[string][]map[string]interface{}), errors.Wrap(err, "error on getting position.")
	}else {
		return mapResult(rows, "customer"), nil
	}
}

func GetPositionByEmail(email string)(map[string][]map[string]interface{}, error){
	connect()
	defer database.Close()
	rows, err := database.Query(SELECTCUSTOMERQUEUEBYEMAIL, email)
	if err != nil {
		log.Println("Error on getting customer: %v", err)
	}
	if rows == nil{
		return make(map[string][]map[string]interface{}), errors.Wrap(err, "error on getting position by email.")
	}else {
		return mapResult(rows, "customer"), nil
	}
}

func GetReferrals(referral_code string)(map[string][]map[string]interface{}, error){
	connect()
	defer database.Close()
	rows, err := database.Query(SELECTREFERRALS, referral_code)
	if err != nil {
		log.Println("Error on getting referrals: %v", err)
	}
	if rows == nil{
		return make(map[string][]map[string]interface{}),errors.Wrap(err, "error on getting referrals.")
	}else {
		return mapResult(rows, "referrals"), nil
	}
}

func GetPreRegistryStep(email string)(map[string][]map[string]interface{}, error){
	connect()
	defer database.Close()
	rows, err := database.Query(SELECTPREREGISTRYSTEP, email)
	if err != nil {
		log.Println("Error on getting referrals: %v", err)
	}
	if rows == nil{
		return make(map[string][]map[string]interface{}), errors.Wrap(err, "error on getting preregistrystep")
	}else {
		return mapResult(rows, "registrystep"), nil
	}
}

func InsertCustomer(documentNumber string, name string, email string, contact string, customer_code string, referral_code string, password string) error{
	connect()
	defer database.Close()
	sttmt := fmt.Sprintf(INSERTPREREGISTRY, documentNumber, name, email, contact, customer_code, referral_code)
	_, err := database.Exec(sttmt)
	if err != nil {
		log.Println("Error on insert new customer pre-registry: %v", err)
		return errors.Wrap(err, "error on insert customer.")
	}
	referral_code = dynamic.GenerateShortId()
	err = insertReferrals(customer_code, referral_code)
	if err != nil {
		return errors.Wrap(err, "error on insert customer referrals.")
	}
	err = insertOnQueue(customer_code)
	if err != nil {
		return errors.Wrap(err, "error on insert customer queue.")
	}

	err = insertAuthentication(customer_code, email)
	if err != nil {
		return errors.Wrap(err, "error on insert customer authentication.")
	}
	err = deletePreRegistryStep(email)
	if err != nil {
		return errors.Wrap(err, "error on delete customer preregistrystep.")
	}

	return nil
}

func UpdateCustomer(customer_code string, contact string) error{
	connect()
	defer database.Close()
	result, err := database.Exec(UPDATEPREREGISTGRY, contact, customer_code )
	if err != nil {
		log.Println("Error on update customer pre-registry: %v", err)
		return errors.Wrap(err, "error on update customer.")
	}
	rows, _ := result.RowsAffected()
	if rows == 0{
		log.Println("0 rows are affected")
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

func InsertOnPreRegistryStep(email string, name string, document_number string, contact string, referral string) error{
	connect()
	defer database.Close()
	sttmt := strings.Replace(INSERTPREREGISTRYSTEP, "@", DefineEmptyOrNot(name, document_number, contact, referral),1)
	sttmt = fmt.Sprintf(sttmt, email, name, document_number, contact, referral)
	_, err := database.Exec(sttmt)
	if err != nil {
		log.Println("Error on insert customer on queue: %v", err)
		return err
	}
	return nil
}

func deletePreRegistryStep(email string) error{
	connect()
	defer database.Close()
	_, err := database.Exec(DELETEPREREGISTRYSTEP, email)
	if err != nil {
		log.Println("Error on delete customer in preregistrystep: %v", err)
		return err
	}
	return nil
}
package domain

import (
	"github.com/guru-invest/guru.api.preregistry/configuration"
	"github.com/guru-invest/guru.api.preregistry/repository"
	"time"
)

func mapPosition(customer_code string) Position {
	customer := Position{}
	rawData := repository.GetPosition(customer_code)
	for _,mp  := range rawData["customer"]{
		customer = Position{
			Customer_Code: mp["customer_code"].(string),
			DocumentNumber: mp["document_number"].(string),
			Name: mp["name"].(string),
			Email: mp["email"].(string),
			Referral_Code:  configuration.CONFIGURATION.OTHER.DeepLinkPrefix + mp["referral_code"].(string),
			Referral_Count: mp["referral_count"].(int64),
			Position: mp["position"].(int64),
			Behind: mp["behind"].(int64),
		}
	}
	return customer
}

func mapPositionByEmail(email string) Position {
	customer := Position{}
	rawData := repository.GetPositionByEmail(email)
	for _,mp  := range rawData["customer"]{
		customer = Position{
			Customer_Code: mp["customer_code"].(string),
			DocumentNumber: mp["document_number"].(string),
			Name: mp["name"].(string),
			Email: mp["email"].(string),
			Referral_Code:  configuration.CONFIGURATION.OTHER.DeepLinkPrefix + mp["referral_code"].(string),
			Referral_Count: mp["referral_count"].(int64),
			Position: mp["position"].(int64),
			Behind: mp["behind"].(int64),
		}
	}
	return customer
}

func mapPreRegistryStep(email string) Customer {
	customer := Customer{}
	rawData := repository.GetPreRegistryStep(email)
	for _,mp  := range rawData["registrystep"]{
		customer = Customer{
			DocumentNumber: mp["document_number"].(string),
			Name: mp["name"].(string),
			Email: mp["email"].(string),
			Referral_Code:  mp["referral"].(string),
		}
	}
	return customer
}

func mapReferrals(referral_code string) Referrals {
	customer := Referrals{}
	accepted := Referrals_Accepted{}
	rawData := repository.GetReferrals(referral_code)
	for _,mp  := range rawData["referrals"]{
		if customer.Origin_Code == "" {
			customer = Referrals{
				Referral_Code: mp["referral"].(string),
				Origin_Name:   mp["origin_name"].(string),
				Origin_Code:   mp["origin_code"].(string),
			}
		}
		accepted = Referrals_Accepted{
			Customer_Name: mp["customer_name"].(string),
			Customer_Code: mp["customer_code"].(string),
			Creation_Date: mp["creation_date"].(time.Time),
		}
		customer.Accepted = append(customer.Accepted, accepted)
	}
	return customer
}

func insert(c *Customer) error{
	if c.Email != "" && c.Name != "" && c.DocumentNumber != ""{
		err := repository.InsertCustomer(c.DocumentNumber, c.Name, c.Email, c.Contact, c.Customer_Code, c.Referral_Code, c.Password)
		if err != nil{
			return err
		}
	}else{
		err := repository.InsertOnPreRegistryStep(c.Email, c.Name, c.DocumentNumber, c.Contact, c.Referral_Code)
		if err != nil{
			return err
		}
	}
	return nil
}

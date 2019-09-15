package domain

import (
	"github.com/guru-invest/guru.api.preregistry/configuration"
	"github.com/guru-invest/guru.api.preregistry/repository"
	"github.com/pkg/errors"
	"time"
)

type POSITIONBYTYPE int

const(
	BYEMAIL POSITIONBYTYPE = iota
	BYCUSTOMERCODE
)

func mapPosition(param string, byType POSITIONBYTYPE) (Position, error) {
	position := Position{}
	rawData, err := repository.GetPosition(param)
	switch byType {
		case BYEMAIL:
			rawData, err = repository.GetPositionByEmail(param)
		case BYCUSTOMERCODE:
			rawData, err = repository.GetPosition(param)
	}
	if err != nil{
		return position, errors.Wrap(err, "error on mapping position")
	}else {
		for _, mp := range rawData["customer"] {
			position = Position{
				Customer_Code:  mp["customer_code"].(string),
				DocumentNumber: mp["document_number"].(string),
				Name:           mp["name"].(string),
				Email:          mp["email"].(string),
				Referral_Code:  configuration.CONFIGURATION.OTHER.DeepLinkPrefix + mp["referral_code"].(string),
				Referral_Count: mp["referral_count"].(int64),
				Position:       mp["position"].(int64),
				Behind:         mp["behind"].(int64),
			}
		}
	}
	return position, nil
}

func mapPreRegistryStep(email string) (Customer, error) {
	customer := Customer{}
	rawData,err := repository.GetPreRegistryStep(email)
	if err != nil{
		return customer, errors.Wrap(err, "error on mapping pre registry step")
	}else {
		for _, mp := range rawData["registrystep"] {
			customer = Customer{
				DocumentNumber: mp["document_number"].(string),
				Name:           mp["name"].(string),
				Email:          mp["email"].(string),
				Referral_Code:  mp["referral"].(string),
			}
		}
	}
	return customer, nil
}

func mapReferrals(referral_code string) (Referrals, error) {
	referrals := Referrals{}
	accepted := Referrals_Accepted{}
	rawData, err := repository.GetReferrals(referral_code)
	if err != nil{
		return referrals, errors.Wrap(err, "error on mapping referrals")
	}else {
		for _, mp := range rawData["referrals"] {
			if referrals.Origin_Code == "" {
				referrals = Referrals{
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
			referrals.Accepted = append(referrals.Accepted, accepted)
		}
	}
	return referrals, nil
}

func insert(c *Customer) error{
	if c.Email != "" && c.Name != "" && c.DocumentNumber != ""{
		err := repository.InsertCustomer(c.DocumentNumber, c.Name, c.Email, c.Contact, c.Customer_Code, c.Referral_Code, c.Password)
		if err != nil{
			return errors.Wrap(err, "error on inserting customer on pre registry.")
		}
	}else{
		err := repository.InsertOnPreRegistryStep(c.Email, c.Name, c.DocumentNumber, c.Contact, c.Referral_Code)
		if err != nil{
			return errors.Wrap(err, "error on inserting customer on pre registry step.")
		}
	}
	return nil
}

func update(c *Customer) error{
	if c.Customer_Code != "" && c.Contact != ""{
		err := repository.UpdateCustomer(c.Customer_Code, c.Contact)
		if err != nil{
			return errors.Wrap(err, "error on update customer.")
		}
	}
	return nil
}

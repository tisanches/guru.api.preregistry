package domain

import (
	"github.com/guru-invest/guru.framework/dynamic"
	"github.com/pkg/errors"
)

type Customer struct{
	DocumentNumber string `json:"document_number,omitempty"`
	Name string `json:"name,omitempty"`
	Email string `json:"email,omitempty"`
	Contact string `json:"contact,omitempty"`
	Password string `json:"password,omitempty"`
	Customer_Code string `json:"customer_code,omitempty"`
	Referral_Code string `json:"referral_code,omitempty"`
}

func (c *Customer)Insert() error{
	if c.Customer_Code == "" {
		c.generateCustomerCode()
	}
	err := insert(c)
	if err != nil{
		return errors.Wrap(err, "Error on insert customer.")
	}
	return nil
}

func (c *Customer)generateCustomerCode(){
	c.Customer_Code = dynamic.GenerateCustomerCode()
}

func (c *Customer) GetByEmail(email string) error{
	err := errors.New("")
	*c, err = mapPreRegistryStep(email)
	if err != nil{
		return errors.Wrap(err, "error on getting customer")
	}
	return nil
}

func (c *Customer) Update() error{
	if c.Customer_Code != "" && c.Contact != ""{
		err := update(c)
		if err != nil{
			return errors.Wrap(err, "Error on update customer.")
		}
	}else{
		err := insert(c)
		if err != nil{
			return errors.Wrap(err, "Error on update customer on preregistrystep.")
		}
	}
	return nil
}

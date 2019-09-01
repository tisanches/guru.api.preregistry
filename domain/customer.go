package domain

import (
	"github.com/guru-invest/guru.framework/dynamic"
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

func (c *Customer)Insert(){
	if c.Customer_Code == "" {
		c.generateCustomerCode()
	}
	insert(c)
}

func (c *Customer)generateCustomerCode(){
	c.Customer_Code = dynamic.GenerateCustomerCode()
}

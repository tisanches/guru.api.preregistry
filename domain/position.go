package domain

import "github.com/pkg/errors"

type Position struct{
	Customer_Code 	string `json:"customer_code,omitempty"`
	DocumentNumber 	string `json:"document_number,omitempty"`
	Name 			string `json:"name,omitempty"`
	Email 			string `json:"email,omitempty"`
	Referral_Code 	string `json:"referral_code,omitempty"`
	Referral_Count 	int64 `json:"referral_count,omitempty"`
	Position 		int64 `json:"position,omitempty"`
	Behind 			int64 `json:"behind,omitempty"`
}

func (p *Position)Get(customer_code string) error {
	err := errors.New("")
	*p, err = mapPosition(customer_code, BYCUSTOMERCODE)
	if err != nil{
		return errors.Wrap(err, "error on getting position")
	}
	return nil
}

func (p *Position)GetByEmail(email string) error{
	err := errors.New("")
	*p, err = mapPosition(email, BYEMAIL)
	if err != nil{
		return errors.Wrap(err, "error on getting position")
	}
	return nil
}

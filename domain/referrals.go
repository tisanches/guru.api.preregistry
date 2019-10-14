package domain

import (
	"github.com/pkg/errors"
	"time"
)

type Referrals struct{
	Referral_Code 	string `json:"referral_code,omitempty"`
	Origin_Name 	string `json:"origin_name,omitempty"`
	Origin_Code 	string `json:"origin_code,omitempty"`
	Accepted []Referrals_Accepted `json:"accepted,omitempty"`
}

type Referrals_Accepted struct{
	Customer_Name 	string `json:"customer_name,omitempty"`
	Customer_Code 	string `json:"customer_code,omitempty"`
	Creation_Date 	time.Time `json:"creation_date,omitempty"`
}

func (r *Referrals)Get(referral_code string) error {
	err := errors.New("")
	*r, err = mapReferrals(referral_code)
	if err != nil{
		return errors.Wrap(err, "error on getting referrals")
	}
	return nil
}


package domain

import (
	"time"
)

type Referrals struct{
	Referral_Code 	string `json:"referral_code,omitempty"`
	Origin_Name 	string `json:"origin_name,omitempty"`
	Origin_Code 	string `json:"origin_code,omitempty"`
	Customer_Name 	string `json:"customer_name,omitempty"`
	Customer_Code 	string `json:"customer_code,omitempty"`
	Creation_Date 	time.Time `json:"creation_date,omitempty"`
}

func (r *Referrals)Get(referral_code string){
	*r = mapReferrals(referral_code)
}
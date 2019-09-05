package domain

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

func (p *Position)Get(customer_code string){
	*p = mapPosition(customer_code)
}

func (p *Position)GetByDocumentNumber(document_number string){
	*p = mapPositionByDocumentNumber(document_number)
}

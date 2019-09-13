package repository

import "strings"

const DOCUMENT_NUMBER = "'%s'"
const NAME = "'%s'"
const EMAIL = "'%s'"
const CONTACT = "'%s'"
const CUSTOMER_CODE = "'%s'"
const REFERRAL_CODE = "'%s'"
const PASSWORD = "'%s'"
const EMPTYNAME = " NAME = ''"
const EMPTYDOCUMENTNUMBER = " , DOCUMENT_NUMBER = ''"
const EMPTYCONTACT = " , CONTACT = ''"
const EMPTYREFERRAL = " , REFERRAL = ''"
const NONEMPTYNAME = " NAME = '@'"
const NONEMPTYDOCUMENTNUMBER = " , DOCUMENT_NUMBER = '@'"
const NONEMPTYCONTACT = " , CONTACT = '@'"
const NONEMPTYREFERRAL = " , REFERRAL = '@'"


const INSERTPREREGISTRY  =
	"INSERT INTO CUSTOMER.PREREGISTRY VALUES(" +
		DOCUMENT_NUMBER + "," +
		NAME + ", " +
		EMAIL + ", " +
		CONTACT + ", " +
		CUSTOMER_CODE + "," +
		REFERRAL_CODE + ")"

const INSERTREFERRAL  =
	"INSERT INTO CUSTOMER.REFERRAL VALUES(" +
		CUSTOMER_CODE + "," +
		REFERRAL_CODE + ", " +
		"0" + ")"

const UPDATEPREREGISTGRY  =
	"UPDATE CUSTOMER.PREREGISTRY SET CONTACT = $1 WHERE CUSTOMER_CODE = $2"

const INSERTQUEUE  =
	"INSERT INTO CUSTOMER.QUEUE (CUSTOMER_CODE) VALUES(" +
		CUSTOMER_CODE + ")"

const SELECTCUSTOMERQUEUE = "SELECT * FROM CUSTOMER.CUSTOMER_VIEW v WHERE v.CUSTOMER_CODE = $1"
const SELECTCUSTOMERQUEUEBYEMAIL = "SELECT * FROM CUSTOMER.CUSTOMER_VIEW v WHERE v.EMAIL = $1"
const SELECTREFERRALS = "SELECT * FROM CUSTOMER.CUSTOMER_REFERRALS_VIEW v WHERE v.REFERRAL = $1"

const INSERTAUTHENTICATION  =
	"INSERT INTO AUTH.AUTHENTICATION VALUES(" +
		CUSTOMER_CODE + "," +
		PASSWORD + ")"

const SELECTPREREGISTRYSTEP = "SELECT EMAIL, NAME, DOCUMENT_NUMBER, CONTACT, REFERRAL FROM CUSTOMER.PREREGISTRYSTEP WHERE EMAIL = $1"

var INSERTPREREGISTRYSTEP  =
	"INSERT INTO CUSTOMER.PREREGISTRYSTEP (EMAIL, NAME, DOCUMENT_NUMBER, CONTACT, REFERRAL) VALUES ( " + EMAIL +
		"," + NAME + "," + DOCUMENT_NUMBER + "," + CONTACT + "," + REFERRAL_CODE + ") ON CONFLICT (EMAIL) DO UPDATE SET " + "@"

	func DefineEmptyOrNot(name string, document_number string, contact string, referral string) string{
		statement := ""
		if name == ""{
			statement += EMPTYNAME
		}else{
			statement +=  strings.Replace(NONEMPTYNAME, "@", name, 1)
		}
		if document_number == ""{
			statement += EMPTYDOCUMENTNUMBER
		}else{
			statement +=  strings.Replace(NONEMPTYDOCUMENTNUMBER, "@", document_number, 1)
		}
		if contact == ""{
			statement += EMPTYCONTACT
		}else{
			statement +=  strings.Replace(NONEMPTYCONTACT, "@", contact, 1)

		}
		if referral == ""{
			statement += EMPTYREFERRAL
		}else{
			statement += strings.Replace(NONEMPTYREFERRAL, "@", referral, 1)
		}
		return statement
	}

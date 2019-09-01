package repository

const DOCUMENT_NUMBER = "'%s'"
const NAME = "'%s'"
const EMAIL = "'%s'"
const CONTACT = "'%s'"
const CUSTOMER_CODE = "'%s'"
const REFERRAL_CODE = "'%s'"
const PASSWORD = "'%s'"

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

const INSERTQUEUE  =
	"INSERT INTO CUSTOMER.QUEUE (CUSTOMER_CODE) VALUES(" +
		CUSTOMER_CODE + ")"

const SELECTCUSTOMERQUEUE = "SELECT * FROM CUSTOMER.CUSTOMER_VIEW v WHERE v.CUSTOMER_CODE = $1"
const SELECTREFERRALS = "SELECT * FROM CUSTOMER.CUSTOMER_REFERRALS_VIEW v WHERE v.REFERRAL = $1"

const INSERTAUTHENTICATION  =
	"INSERT INTO AUTH.AUTHENTICATION VALUES(" +
		CUSTOMER_CODE + "," +
		PASSWORD + ")"

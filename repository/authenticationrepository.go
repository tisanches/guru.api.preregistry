package repository

import (
	"fmt"
	"github.com/pkg/errors"
)

func insertAuthentication(customer_code string, email string) error{
	connect()
	defer database.Close()
	sttmt := fmt.Sprintf(INSERTAUTHENTICATION, customer_code, email)
	_, err := database.Exec(sttmt)
	if err != nil {
		return errors.Wrap(err, "error on insert authentication token.")
	}
	return nil
}

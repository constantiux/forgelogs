package myparser

import (
	"fmt"
	"strings"
)

var FrSource string

func ValidateArgSource(args []string) error {
	argString := strings.Join(args, "")
	if strings.Compare(argString, "am") == 0 {
		FrSource = "am-everything"
	} else if strings.Compare(argString, "idm") == 0 {
		FrSource = "idm-everything"
	} else {
		return fmt.Errorf("%v is invalid, only accepts \"am\" or \"idm\"", argString)
	}
	return nil
}

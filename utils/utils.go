package utils

import (
	"fmt"
	"os"
)

// errors is a map of error output value in ErrorHandler
var errors = map[string]string{
	"web":        "😮 Oops! Something went wrong",
	"restricted": "😣 Oops! this is a restricted path.\nplease use another path.",
}

// ErrorHandler outputs errors and safely exits the program
func ErrorHandler(errType string) {
	fmt.Println(errors[errType])
	os.Exit(0)
}
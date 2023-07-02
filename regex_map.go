package main

import (
	"regexp"
)

type Schema struct {
	schema map[string]string
}

func Validate() Schema {

	schema := make(map[string]string)
	regexes := make(map[string]string)

	regexes["Only letters"] = `^[A-Za-z\s]+$`
	regexes["Only numbers"] = `^[0-9]+$`
	regexes["Only letters and numbers"] = `^[A-Za-z0-9]+$`
	regexes["Letters, numbers and spaces"] = `^[A-Za-z0-9\s]+$`
	regexes["Letters and hyphens"] = `^[A-Za-z-]+$`
	regexes["Phone number"] = `^[0-9]{3}-[0-9]{3}-[0-9]{4}$`
	regexes["Email"] = `^[a-zA-Z0-9.!#$%&'*+/=?^_{|}~-]+@[a-zA-Z0-9](?:[a-zA-Z0-9-]{0,61}[a-zA-Z0-9])?(?:\.[a-zA-Z0-9](?:[a-zA-Z0-9-]{0,61}[a-zA-Z0-9])?)*$`
	regexes["Date"] = `^[0-9]{2}\/[0-9]{2}\/[0-9]{4}$`
	regexes["SSN"] = `^[0-9]{3}-[0-9]{2}-[0-9]{4}$`
	regexes["Zip"] = `^[0-9]{5}$`

	// Add the field to the map and its regex
	schema["Email"] = regexes["Email"]
	schema["First Name"] = regexes["Letters and hyphens"]
	schema["Middle Name"] = regexes["Letters and hyphens"]
	schema["Last Name"] = regexes["Letters and hyphens"]
	schema["Phone Number"] = regexes["Phone number"]
	schema["DL Number"] = regexes["Only letters and numbers"]
	schema["DOB"] = regexes["Date"]
	schema["SSN"] = regexes["SSN"]
	schema["Address"] = regexes["Letters, numbers and spaces"]
	schema["City"] = regexes["Only letters"]
	schema["State"] = regexes["Only letters"]
	schema["Country"] = regexes["Only letters"]
	schema["Zip"] = regexes["Zip"]
	schema["Apt Number"] = regexes["Only letters and numbers"]
	schema["Rent"] = regexes["Only numbers"]
	schema["Employer"] = regexes["Only letters"]
	schema["Job Title"] = regexes["Only letters"]
	schema["Employer Address"] = regexes["Letters, numbers and spaces"]
	schema["Employer City"] = regexes["Only letters"]
	schema["Employer State"] = regexes["Only letters"]
	schema["Occupation"] = regexes["Only letters"]
	schema["Employer Zip"] = regexes["Zip"]
	schema["Employer Phone"] = regexes["Phone number"]
	schema["Monthly Income"] = regexes["Only numbers"]

	return Schema{schema}

}

func CheckParamsAgainstRegex(key string, value string, schema Schema) bool {

	regex, ok := schema.schema[key]

	if !ok {
		return false
	}

	match, error := regexp.MatchString(regex, value)

	if error != nil || !match {
		return false
	}

	return true
}

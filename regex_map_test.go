package main

import (
	"testing"
)

func TestValidateRegexesTableDriven(t *testing.T) {
	var tests = []struct {
		name  string
		value string
		want  bool
	}{
		{"Email", "dd", false},
		{"Email", "test2@mail.com", true},
		{"First Name", "dd", true},
		{"First Name", "rick-ashley", true},
		{"First Name", "dd2", false},
		{"Phone Number", "dd", false},
		{"Phone Number", "123-456-7890", true},
		{"Phone Number", "123-456-789", false},
		{"Phone Number", "123-456-78901", false},
		{"DL Number", "dd", true},
		{"DL Number", "dd2d", true},
		{"DL Number", "dd-!2dd", false},
		{"DOB", "dd", false},
		{"DOB", "12/12/1990", true},
		{"DOB", "12/12/90", false},
		{"SSN", "dd", false},
		{"SSN", "123-45-6789", true},
		{"SSN", "123-45-678", false},
		{"Address", "dd", true},
		{"Address", "1234 Main St", true},
		{"Address", "1234 Main St Apt 1", true},
		{"Address", "1234 Main St - Apt 1", false},
		{"City", "dd", true},
		{"City", "New York", true},
		{"City", "New York-City", false},
		{"State", "NY", true},
		{"State", "New York", true},
		{"State", "New-York", false},
		{"Country", "USA", true},
		{"Country", "United States", true},
		{"Country", "United-States", false},
		{"Zip", "dd", false},
		{"Zip", "12345", true},
		{"Zip", "1234", false},
		{"Apt Number", "dd", true},
		{"Apt Number", "1", true},
		{"Apt Number", "1A", true},
		{"Apt Number", "1-A", false},
		{"Rent", "dd", false},
		{"Rent", "1000", true},
		{"Rent", "1000.00", false},
		{"Employer", "McDonalds", true},
		{"Employer", "McDonalds-2", false},
		{"Job Title", "Manager", true},
		{"Job Title", "Manager-2", false},
		{"Employer Address", "1234 Main St", true},
		{"Employer Address", "1234 Main St Apt 1", true},
		{"Employer Address", "1234 Main St - Apt 1", false},
		{"Employer City", "New York", true},
		{"Employer City", "New York-City", false},
		{"Employer State", "NY", true},
		{"Employer State", "New York", true},
		{"Employer State", "New-York", false},
		{"Occupation", "dd", true},
		{"Occupation", "dd2", false},
		{"Employer Zip", "dd", false},
		{"Employer Zip", "12345", true},
		{"Employer Zip", "1234", false},
		{"Employer Phone", "dd", false},
		{"Employer Phone", "123-456-7890", true},
		{"Employer Phone", "123-456-789", false},
		{"Employer Phone", "123-456-78901", false},
		{"Monthly Income", "dd", false},
		{"Monthly Income", "1000", true},
		{"Monthly Income", "1000.00", false},
	}

	schema := Validate()

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			ans := CheckParamsAgainstRegex(test.name, test.value, schema)
			if ans != test.want {
				t.Errorf("got %v, want %v", ans, test.want)
			}
		})
	}

}

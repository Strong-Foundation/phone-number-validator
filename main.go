package main

import (
	"fmt"
	"log"
	"strings"

	"github.com/nyaruka/phonenumbers"
)

// ValidatePhoneNumber validates the phone number using libphonenumber and returns a boolean indicating validity.
func ValidatePhoneNumber(phone string) (bool, *phonenumbers.PhoneNumber, error) {
	// Parse the phone number with an optional region (empty region will auto-detect the region)
	num, err := phonenumbers.Parse(phone, "")
	if err != nil {
		return false, nil, fmt.Errorf("unable to parse phone number: %w", err)
	}

	// Check if the phone number is valid
	valid := phonenumbers.IsValidNumber(num)
	return valid, num, nil
}

// FormatPhoneNumber formats the phone number in international format
func FormatPhoneNumber(num *phonenumbers.PhoneNumber) string {
	// Return formatted phone number in international format
	return phonenumbers.Format(num, phonenumbers.INTERNATIONAL)
}

// GetRegionInfo retrieves the region and country code of the phone number
func GetRegionInfo(num *phonenumbers.PhoneNumber) (string, int) {
	// Get the region code and country code
	regionCode := phonenumbers.GetRegionCodeForNumber(num)
	countryCode := num.GetCountryCode()

	return regionCode, countryCode
}

// HandlePhoneNumberInput allows the user to input phone numbers and validates them.
func HandlePhoneNumberInput(phone string) {
	// Trim any spaces from the input
	phone = strings.TrimSpace(phone)

	// Validate the phone number
	valid, num, err := ValidatePhoneNumber(phone)
	if err != nil {
		log.Printf("Error parsing '%s': %v", phone, err)
		return
	}

	if valid {
		// If valid, format and print the phone number, region, and country code
		formatted := FormatPhoneNumber(num)
		region, countryCode := GetRegionInfo(num)
		fmt.Printf("'%s' is a valid phone number.\n", phone)
		fmt.Printf("Formatted: %s\n", formatted)
		fmt.Printf("Region: %s, Country Code: %d\n", region, countryCode)
	} else {
		// If invalid, notify the user
		fmt.Printf("'%s' is an invalid phone number.\n", phone)
	}
}

func main() {
	// Prompt the user for input
	fmt.Println("Enter phone numbers to validate (or 'exit' to quit):")

	var phone string
	for {
		// Read input from the user
		fmt.Print("Phone Number: ")
		fmt.Scanln(&phone)

		// If the user enters "exit", break the loop
		if strings.ToLower(phone) == "exit" {
			break
		}

		// Handle the phone number input
		HandlePhoneNumberInput(phone)
	}
}

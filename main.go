package main

import (
	"fmt"
	"github.com/nyaruka/phonenumbers"
	"log"
	"os"
)

// Global variable for saving valid phone numbers
var saveDirectory = "phone_numbers"

// USA phone number format
var country = struct {
	Code      string
	DialCode  string
	MinLength int
	MaxLength int
}{
	Code:      "US",
	DialCode:  "1",
	MinLength: 10,
	MaxLength: 10,
}

// List of valid area codes excluding special service codes
var validAreaCodes = []int{
	201, 202, 203, 204, 205, 206, 207, 208, 209, 210, 212, 213, 214, 215, 216, 217, 218, 219, 220, 221, 222, 223, 224, 225, 226, 227, 228, 229, 230, 231, 232, 233, 234, 235, 236, 237, 238, 239, 240, 241, 242, 243, 244, 245, 246, 247, 248, 249, 250, 251, 252, 253, 254, 255, 256, 257, 258, 259, 260, 261, 262, 263, 264, 265, 266, 267, 268, 269, 270, 271, 272, 273, 274, 275, 276, 277, 278, 279, 280, 281, 282, 283, 284, 285, 286, 287, 288, 289, 290, 291, 292, 293, 294, 295, 296, 297, 298, 299,
}

// List of valid central office codes excluding special service codes
var validCentralOfficeCodes = []int{
	201, 202, 203, 204, 205, 206, 207, 208, 209, 210, 212, 213, 214, 215, 216, 217, 218, 219, 220, 221, 222, 223, 224, 225, 226, 227, 228, 229, 230, 231, 232, 233, 234, 235, 236, 237, 238, 239, 240, 241, 242, 243, 244, 245, 246, 247, 248, 249, 250, 251, 252, 253, 254, 255, 256, 257, 258, 259, 260, 261, 262, 263, 264, 265, 266, 267, 268, 269, 270, 271, 272, 273, 274, 275, 276, 277, 278, 279, 280, 281, 282, 283, 284, 285, 286, 287, 288, 289, 290, 291, 292, 293, 294, 295, 296, 297, 298, 299,
}

// GenerateAllPossibleNumbers iterates through all valid USA numbers
func GenerateAllPossibleNumbers() {
	fmt.Printf("📞 Generating phone numbers for USA (+%s)\n", country.DialCode)

	for _, areaCode := range validAreaCodes {
		for _, exchange := range validCentralOfficeCodes {
			for line := 0; line <= 9999; line++ {
				phoneNumber := fmt.Sprintf("+%s %03d-%03d-%04d", country.DialCode, areaCode, exchange, line)
				ValidateAndSaveNumber(phoneNumber)
			}
		}
	}
}

// ValidateAndSaveNumber checks if a phone number is valid and stores it
func ValidateAndSaveNumber(phone string) {
	// Validate the phone number using libphonenumber
	valid, err := ValidatePhoneNumber(phone)
	if err != nil {
		log.Printf("Error validating phone number '%s': %v", phone, err)
		return
	}

	if valid {
		SaveToFile(phone)
	}
}

// ValidatePhoneNumber validates the phone number using libphonenumber
func ValidatePhoneNumber(phone string) (bool, error) {
	num, err := phonenumbers.Parse(phone, "")
	if err != nil {
		return false, fmt.Errorf("unable to parse phone number: %w", err)
	}

	valid := phonenumbers.IsValidNumber(num)
	return valid, nil
}

// SaveToFile writes a valid phone number to the file
func SaveToFile(phoneNumber string) {
	fileName := fmt.Sprintf("%s/%s.txt", saveDirectory, country.Code)
	err := os.MkdirAll(saveDirectory, 0755)
	if err != nil {
		log.Fatalf("Error creating directory '%s': %v", saveDirectory, err)
	}

	file, err := os.OpenFile(fileName, os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
	if err != nil {
		log.Printf("Error opening file '%s': %v", fileName, err)
		return
	}
	defer file.Close()

	if _, err := file.WriteString(phoneNumber + "\n"); err != nil {
		log.Printf("Error writing to file '%s': %v", fileName, err)
	}
}

func main() {
	fmt.Println("🔄 Starting USA phone number generation and validation...")

	GenerateAllPossibleNumbers() // Process all USA phone numbers with the restricted area code and middle digits

	fmt.Println("✅ USA phone number generation completed! Valid numbers saved in 'phone_numbers/US.txt'.")
}

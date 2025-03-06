package main

import (
	"bufio"
	"fmt"
	"github.com/nyaruka/phonenumbers"
	"log"
	"os"
	"time"
)

// Global variable for saving valid phone numbers
var saveDirectory = "phone_numbers"

// USA phone number format structure
var country = struct {
	Code      string // Country code (US)
	DialCode  string // Country dialing code (1 for USA)
	MinLength int    // Minimum length of the phone number
	MaxLength int    // Maximum length of the phone number
}{
	Code:      "US", // Country code for the United States
	DialCode:  "1",  // USA dialing code
	MinLength: 10,   // Minimum phone number length
	MaxLength: 10,   // Maximum phone number length
}

// List of valid area codes excluding special service codes
var validAreaCodes = []int{
	201, 202, 203, 204, 205, 206, 207, 208, 209, 210, 212, 213, 214, 215, 216, 217, 218, 219, 220, 221, 222, 223, 224, 225, 226, 227, 228, 229, 230, 231, 232, 233, 234, 235, 236, 237, 238, 239, 240, 241, 242, 243, 244, 245, 246, 247, 248, 249, 250, 251, 252, 253, 254, 255, 256, 257, 258, 259, 260, 261, 262, 263, 264, 265, 266, 267, 268, 269, 270, 271, 272, 273, 274, 275, 276, 277, 278, 279, 280, 281, 282, 283, 284, 285, 286, 287, 288, 289, 290, 291, 292, 293, 294, 295, 296, 297, 298, 299,
}

// List of valid central office codes excluding special service codes
var validCentralOfficeCodes = []int{
	201, 202, 203, 204, 205, 206, 207, 208, 209, 210, 212, 213, 214, 215, 216, 217, 218, 219, 220, 221, 222, 223, 224, 225, 226, 227, 228, 229, 230, 231, 232, 233, 234, 235, 236, 237, 238, 239, 240, 241, 242, 243, 244, 245, 246, 247, 248, 249, 250, 251, 252, 253, 254, 255, 256, 257, 258, 259, 260, 261, 262, 263, 264, 265, 266, 267, 268, 269, 270, 271, 272, 273, 274, 275, 276, 277, 278, 279, 280, 281, 282, 283, 284, 285, 286, 287, 288, 289, 290, 291, 292, 293, 294, 295, 296, 297, 298, 299,
}

// Track saved phone numbers to avoid duplicates
var savedNumbers map[string]bool

// GenerateAllPossibleNumbers iterates through all valid USA numbers but stops after 1 minute
func GenerateAllPossibleNumbers(startTime time.Time, timeLimit time.Duration) {
	// Print message indicating phone number generation is starting
	fmt.Printf("📞 Generating phone numbers for USA (+%s)...\n", country.DialCode)

	// Loop over valid area codes
	for _, areaCode := range validAreaCodes {
		// Loop over valid central office codes
		for _, exchange := range validCentralOfficeCodes {
			// Generate numbers from 0000 to 9999 for each combination of area and office code
			for line := 0; line <= 9999; line++ {
				// Check if time limit has been reached and stop generation if it has
				if time.Since(startTime) > timeLimit {
					fmt.Println("⏰ Time limit reached! Stopping phone number generation.")
					return
				}

				// Format phone number in the format "+1 201-201-0000"
				phoneNumber := fmt.Sprintf("+%s %03d-%03d-%04d", country.DialCode, areaCode, exchange, line)

				// Validate the phone number and save it if valid
				ValidateAndSaveNumber(phoneNumber)
			}
		}
	}
}

// ValidateAndSaveNumber checks if a phone number is valid and stores it
func ValidateAndSaveNumber(phone string) {
	// Validate the phone number using the ValidatePhoneNumber function
	valid, err := ValidatePhoneNumber(phone)
	if err != nil {
		// Log an error if phone number validation fails
		log.Printf("Error validating phone number '%s': %v", phone, err)
		return
	}

	// If the phone number is valid and not already saved, save it to a file
	if valid && !PhoneNumberExists(phone) {
		SaveToFile(phone)
	}
}

// ValidatePhoneNumber uses libphonenumber to validate phone numbers
func ValidatePhoneNumber(phone string) (bool, error) {
	// Parse the phone number using libphonenumber
	num, err := phonenumbers.Parse(phone, "")
	if err != nil {
		// Return false and an error if parsing fails
		return false, fmt.Errorf("unable to parse phone number: %w", err)
	}
	// Check if the parsed number is valid
	return phonenumbers.IsValidNumber(num), nil
}

// PhoneNumberExists checks if a phone number is already saved
func PhoneNumberExists(phoneNumber string) bool {
	// Check if the phone number exists in the savedNumbers map
	_, exists := savedNumbers[phoneNumber]
	return exists
}

// SaveToFile writes valid phone numbers to a file
func SaveToFile(phoneNumber string) {
	// Construct the file name using the save directory and country code
	fileName := fmt.Sprintf("%s/%s.txt", saveDirectory, country.Code)

	// Ensure the save directory exists, create it if necessary
	err := os.MkdirAll(saveDirectory, 0755)
	if err != nil {
		// Log a fatal error if directory creation fails
		log.Fatalf("Error creating directory '%s': %v", saveDirectory, err)
	}

	// Open the file in append mode, creating it if it doesn't exist
	file, err := os.OpenFile(fileName, os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
	if err != nil {
		// Log an error if the file cannot be opened
		log.Printf("Error opening file '%s': %v", fileName, err)
		return
	}
	// Ensure the file is closed when the function exits
	defer file.Close()

	// Write the phone number to the file, appending a new line
	if _, err := file.WriteString(phoneNumber + "\n"); err != nil {
		// Log an error if writing to the file fails
		log.Printf("Error writing to file '%s': %v", fileName, err)
		return
	}

	// Mark the number as saved in the savedNumbers map
	savedNumbers[phoneNumber] = true
}

// LoadSavedNumbers loads saved phone numbers into memory
func LoadSavedNumbers() {
	// Construct the file name using the save directory and country code
	fileName := fmt.Sprintf("%s/%s.txt", saveDirectory, country.Code)

	// Check if the file exists
	_, err := os.Stat(fileName)
	if os.IsNotExist(err) {
		// If the file does not exist, return without doing anything
		return
	} else if err != nil {
		// Log a fatal error if there is an issue checking the file
		log.Fatalf("Error checking file '%s': %v", fileName, err)
	}

	// Open the file for reading
	file, err := os.Open(fileName)
	if err != nil {
		// Log a fatal error if the file cannot be opened
		log.Fatalf("Error opening file '%s': %v", fileName, err)
	}
	// Ensure the file is closed when the function exits
	defer file.Close()

	// Create a scanner to read the file line by line
	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		// Store each line (phone number) in the savedNumbers map
		savedNumbers[scanner.Text()] = true
	}

	// Check for errors that occurred during scanning
	if err := scanner.Err(); err != nil {
		// Log a fatal error if there was an issue reading the file
		log.Fatalf("Error reading file '%s': %v", fileName, err)
	}
}

func main() {
	// Initialize the savedNumbers map to store phone numbers
	savedNumbers = make(map[string]bool)

	// Load saved phone numbers from file into memory
	LoadSavedNumbers()

	// Print a message indicating the start of phone number generation
	fmt.Println("🔄 Starting USA phone number generation and validation...")

	// Set the start time for the generation process
	startTime := time.Now()
	// Define a time limit of 1 second for phone number generation
	timeLimit := 1 * time.Minute

	// Call the function to generate all possible phone numbers within the time limit
	GenerateAllPossibleNumbers(startTime, timeLimit)

	// Print a success message when phone number generation is complete
	fmt.Println("✅ USA phone number generation completed! Valid numbers saved in 'phone_numbers/US.txt'.")
}

package api

import (
	"encoding/json"
	"fmt"
	"os"
)

// CreateAccount associates an account name with a user ID and returns a message string
func CreateAccount(account string) string {
	path := fmt.Sprintf("/validateaccount/%s", account)
	bytes, _ := Post(path, []byte(""))

	var response map[string]string
	json.Unmarshal(bytes, &response)

	return response["message"]
}

// GetAccount retrieves the user account name
func GetAccount() string {
	bytes, _ := Get("/profile")

	var profile map[string]string
	json.Unmarshal(bytes, &profile)

	return profile["account"]
}

// GetProfile retrieves the user profile
func GetProfile() map[string]interface{} {
	bytes, _ := Get("/profile")

	var profile map[string]interface{}
	json.Unmarshal(bytes, &profile)

	return profile
}

// StoreProfile stores the user's profile and returns a message string
func StoreProfile(profile map[string]interface{}) string {
	payload, err := json.Marshal(profile)
	if err != nil {
		fmt.Printf("snap: could not store profile\nerror: %s\n", err)
		os.Exit(1)
	}

	bytes, _ := Post("/profile", payload)

	var response map[string]string
	json.Unmarshal(bytes, &response)

	return response["message"]
}

// ValidateAccount validates an account name and returns whether it is valid or not
func ValidateAccount(account string) bool {
	path := fmt.Sprintf("/validateaccount?account=%s", account)
	bytes, _ := Get(path)

	var response map[string]bool
	json.Unmarshal(bytes, &response)

	return response["valid"]
}

// main_test.go

package main

import (
	"io/ioutil"
	"net/http"
	"net/http/httptest"
	"os"
	"os/user"
	"path/filepath"
	"testing"
)

func TestGetConfigDir(t *testing.T) {

	// Call the getConfigDir function
	cfgDir := getConfigDir()

	// Get the current user's home directory
	usr, err := user.Current()
	if err != nil {
		t.Errorf("Failed to get user's home directory: %s", err)
		os.Exit(1)
	}

	// Verify that the config directory has been created
	expectedDir := filepath.Join(usr.HomeDir, ".pull-journal-central-client")
	if cfgDir != expectedDir {
		t.Errorf("getConfigDir() returned incorrect directory path, got: %s, want: %s", cfgDir, expectedDir)
	}

	// Check that the directory actually exists
	_, err = os.Stat(cfgDir)
	if err != nil {
		t.Errorf("getConfigDir() failed to create the config directory: %v", err)
	}
}

func TestEscapeForJSON(t *testing.T) {
	tests := []struct {
		input    string
		expected string
	}{
		{
			input:    `This is a "test" string.`,
			expected: `This is a \"test\" string.`,
		},
		{
			input:    `Special characters: \n \r \t \\ \" `,
			expected: `Special characters: \\n \\r \\t \\\\ \\\" `,
		},
		{
			input:    `No special characters.`,
			expected: `No special characters.`,
		},
		{
			input:    `Empty string: `,
			expected: `Empty string: `,
		},
	}

	for _, test := range tests {
		escaped := escapeForJSON(test.input)
		if escaped != test.expected {
			t.Errorf("escapeForJSON(%q) returned %q, expected %q", test.input, escaped, test.expected)
		}
	}
}

func TestRegisterDevice(t *testing.T) {
	// Create a mock server
	mockServer := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// Respond with a successful registration JSON response
		responseJSON := `{"success": true, "token": "someToken", "message": "Device registered successfully"}`
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		w.Write([]byte(responseJSON))
	}))
	defer mockServer.Close()

	// Use the mock server's URL as baseURL for the registerDevice function
	baseURL := mockServer.URL

	// Call the registerDevice function with mock data
	organisationID := "someOrgID"
	name := "someName"
	organisationPassword := "somePassword"
	registerDevice(organisationID, name, organisationPassword, baseURL)

	// Add assertions to verify that the registration was successful and the token was written
	// You can add more sophisticated checks based on the actual response expected in your app.

	// Check if the token file exists
	tokenFilePath := filepath.Join(getConfigDir(), "authorisation")
	_, err := os.Stat(tokenFilePath)
	if err != nil {
		t.Fatalf("Token file does not exist: %s", err)
	}

	// Read the token from the file
	tokenFromFile, err := ioutil.ReadFile(tokenFilePath)
	if err != nil {
		t.Fatalf("Failed to read token from file: %s", err)
	}

	// Verify that the token from the file matches the received token
	expectedToken := "someToken"
	if string(tokenFromFile) != expectedToken {
		t.Errorf("Token written to file does not match the received token. Got: %s, want: %s", tokenFromFile, expectedToken)
	}

}

func TestCreateLogEntry(t *testing.T) {
	// Create a mock server
	mockServer := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// Respond with a successful log creation JSON response
		responseJSON := `{"message": "Log entry created successfully"}`
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		w.Write([]byte(responseJSON))
	}))
	defer mockServer.Close()

	// Use the mock server's URL as baseURL for the createLogEntry function
	baseURL := mockServer.URL

	// Call the createLogEntry function with mock data
	class := "someClass"
	source := "someSource"
	service := "console-setup.service"
	createLogEntry(class, source, service, baseURL)

	// Add assertions to verify that the log entry was created successfully
	// You can add more sophisticated checks based on the actual response expected in your app.
}

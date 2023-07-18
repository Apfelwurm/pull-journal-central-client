package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"os"

	"github.com/spf13/cobra"
)

const baseURL = "http://localhost"

type ApiResponse struct {
	Success bool   `json:"success"`
	Token   string `json:"token"`
	Message string `json:"message"`
}

type ApiError struct {
	Message string            `json:"message"`
	Errors  map[string][]string `json:"errors"`
}

func main() {
	var rootCmd = &cobra.Command{Use: "app"}
	var organisationID, name, organisationPassword string

	registerCmd := &cobra.Command{
		Use:   "register",
		Short: "Register a device",
		Run: func(cmd *cobra.Command, args []string) {
			registerDevice(organisationID, name, organisationPassword)
		},
	}

	rootCmd.AddCommand(registerCmd)

	registerCmd.Flags().StringVar(&organisationID, "organisationID", "", "Organisation ID")
	registerCmd.Flags().StringVar(&name, "name", "", "Name")
	registerCmd.Flags().StringVar(&organisationPassword, "organisationpassword", "", "Organisation Password")

	if err := rootCmd.Execute(); err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}

func registerDevice(organisationID, name, organisationPassword string) {
	// Read device identifier from the file
	deviceIdentifier, err := ioutil.ReadFile("/etc/machine-id")
	if err != nil {
		fmt.Println("Failed to read device identifier from file:", err)
		os.Exit(1)
	}

	// Create the URL with query parameters
	url := fmt.Sprintf("%s/api/devices/register/%s?name=%s&organisationpassword=%s&deviceidentifier=%s",
		baseURL, organisationID, name, organisationPassword, string(bytes.TrimSpace(deviceIdentifier)))

	// Create the HTTP GET request
    req, err := http.NewRequest("GET", url, nil)
    if err != nil {
        fmt.Println("Failed to create HTTP request:", err)
        os.Exit(1)
    }

    // Add headers to the request
    req.Header.Add("Content-Type", "application/json")
    req.Header.Add("Accept", "application/json")

    // Send the HTTP request
    client := http.Client{}
    resp, err := client.Do(req)
    if err != nil {
        fmt.Println("Failed to send HTTP request:", err)
        os.Exit(1)
    }
    defer resp.Body.Close()

	// Check if the response is successful
	if resp.StatusCode >= 200 && resp.StatusCode < 300 {
		var apiResponse ApiResponse
		err := json.NewDecoder(resp.Body).Decode(&apiResponse)
		if err != nil {
			fmt.Println("Failed to decode JSON response:", err)
			os.Exit(1)
		}
		fmt.Printf("Bearer %+v\n", apiResponse.Token)
	} else {
		var apiError ApiError
		err := json.NewDecoder(resp.Body).Decode(&apiError)
		if err != nil {
			fmt.Println("Failed to decode JSON error response:", err)
			os.Exit(1)
		}
		fmt.Println("Error:", apiError.Message)
	}
}

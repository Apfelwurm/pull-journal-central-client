package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"os"
	"os/exec"
	"os/user"
	"path/filepath"
	"strings"

	"github.com/spf13/cobra"
)

const baseURL = "http://localhost"
const configdir = ".pull-journal-central-client"
const conttype = "application/json"

type ApiResponse struct {
	Success bool   `json:"success"`
	Token   string `json:"token"`
	Message string `json:"message"`
}

type ApiError struct {
	Message string              `json:"message"`
	Errors  map[string][]string `json:"errors"`
}

func main() {
	var rootCmd = &cobra.Command{Use: "app"}
	var organisationID, name, organisationPassword string
	var class, source, service string

	registerCmd := &cobra.Command{
		Use:   "register",
		Short: "Register a device",
		Run: func(cmd *cobra.Command, args []string) {
			registerDevice(organisationID, name, organisationPassword)
		},
	}

	logCmd := &cobra.Command{
		Use:   "log",
		Short: "Create a log entry",
		Run: func(cmd *cobra.Command, args []string) {
			createLogEntry(class, source, service)
		},
	}

	rootCmd.AddCommand(registerCmd)
	rootCmd.AddCommand(logCmd)

	registerCmd.Flags().StringVar(&organisationID, "organisationID", "", "Organisation ID")
	registerCmd.Flags().StringVar(&name, "name", "", "Name")
	registerCmd.Flags().StringVar(&organisationPassword, "organisationpassword", "", "Organisation Password")

	logCmd.Flags().StringVar(&class, "class", "", "class of the Log Entry")
	logCmd.Flags().StringVar(&source, "source", "", "source of the log Entry")
	logCmd.Flags().StringVar(&service, "service", "", "service name")

	if err := rootCmd.Execute(); err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}

func getConfigDir() string {

	// Get the current user's home directory
	usr, err := user.Current()
	if err != nil {
		fmt.Println("Failed to get user's home directory:", err)
		os.Exit(1)
	}

	homeDir := usr.HomeDir
	fulcfgDir := filepath.Join(homeDir, configdir)

	err = os.MkdirAll(fulcfgDir, 0700)
	if err != nil {
		fmt.Println("Failed to create config directory:", err)
		os.Exit(1)
	}
	return fulcfgDir
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
	req.Header.Add("Content-Type", conttype)
	req.Header.Add("Accept", conttype)

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

		// Write the token to the file in the home directory
		authToken := []byte(apiResponse.Token)

		filePath := filepath.Join(getConfigDir(), "authorisation")

		err = ioutil.WriteFile(filePath, authToken, 0600)
		if err != nil {
			fmt.Println("Failed to write token to file:", err)
			os.Exit(1)
		}
		fmt.Println("Token written successfully")

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

func createLogEntry(class, source, service string) {

	// Create the URL
	url := fmt.Sprintf("%s/api/logEntries/create", baseURL)

	// Read the output of the "ls" command and escape it for JSON
	content, err := executeServiceCommand(service)
	if err != nil {
		fmt.Println("Failed to read the output of service command:", err)
		os.Exit(1)
	}

	// Create the HTTP POST request body
	data := map[string]string{
		"source":  source,
		"class":   class,
		"content": content,
	}

	// Convert the data map to JSON
	jsonData, err := json.Marshal(data)
	if err != nil {
		fmt.Println("Failed to marshal JSON data:", err)
		os.Exit(1)
	}

	// Create the HTTP POST request
	req, err := http.NewRequest("POST", url, bytes.NewBuffer(jsonData))
	if err != nil {
		fmt.Println("Failed to create HTTP request:", err)
		os.Exit(1)
	}

	// Add headers to the request
	req.Header.Add("Content-Type", conttype)
	req.Header.Add("Accept", conttype)

	// Read the authorization token from the file
	token, err := ioutil.ReadFile(filepath.Join(getConfigDir(), "authorisation"))
	if err != nil {
		fmt.Println("Failed to read authorization token from config:", err)
		os.Exit(1)
	}

	// Add the Authorization header with the token
	req.Header.Add("Authorization", "Bearer "+string(token))

	// Send the HTTP request
	client := http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		fmt.Println("Failed to send HTTP request:", err)
		os.Exit(1)
	}
	defer resp.Body.Close()

	// Check the response status code
	if resp.StatusCode >= 200 && resp.StatusCode < 300 {
		fmt.Println("Log entry created successfully")
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

func executeServiceCommand(service string) (string, error) {

	invocid := exec.Command("systemctl", "show", "-p", "InvocationID", "--value", service)
	// Capture the command output
	var stdoutid, stderrid bytes.Buffer
	invocid.Stdout = &stdoutid
	invocid.Stderr = &stderrid

	// Run the command
	iderr := invocid.Run()
	if iderr != nil {
		return "", fmt.Errorf("command execution failed: %v, stderr: %v", iderr, stderrid.String())
	}

	// Get the output as a string
	idoutput := stdoutid.String()

	fmt.Println("invocid:", idoutput)

	// Execute the "ls" command
	cmd := exec.Command("journalctl", "_SYSTEMD_INVOCATION_ID="+string(strings.ReplaceAll(idoutput, "\n", "")), "--no-pager")

	// Capture the command output
	var stdout, stderr bytes.Buffer
	cmd.Stdout = &stdout
	cmd.Stderr = &stderr

	// Run the command
	err := cmd.Run()
	if err != nil {
		return "", fmt.Errorf("command execution failed: %v, stderr: %v", err, stderr.String())
	}

	// Get the output as a string
	output := stdout.String()

	fmt.Println("output:", output)

	// Escape the output for JSON
	escapedOutput := escapeForJSON(output)

	return escapedOutput, nil
}

func escapeForJSON(str string) string {
	str = strings.ReplaceAll(str, "\\", "\\\\")
	str = strings.ReplaceAll(str, "\n", "\\n")
	str = strings.ReplaceAll(str, "\r", "\\r")
	str = strings.ReplaceAll(str, "\t", "\\t")
	str = strings.ReplaceAll(str, "\"", "\\\"")
	return str
}

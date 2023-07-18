// main.rs

use std::fs;
use clap::{App, Arg, SubCommand};
use reqwest::header::{HeaderMap, ACCEPT, CONTENT_TYPE};
use serde::{Serialize, Deserialize};
use serde_json::json;


const BASE_URL: &str = "http://localhost";

#[derive(Debug, Serialize, Deserialize)]
struct ApiResponse {
    success: bool,
    token: String,
    message: String,
}

#[derive(Debug, Serialize, Deserialize)]
struct ApiError {
    message: String,
    errors: Option<serde_json::Value>,
}

fn main() {
    let matches = App::new("API Client")
        .subcommand(SubCommand::with_name("register")
            .about("Register a device")
            .arg(Arg::with_name("organisationID")
                .required(true)
                .takes_value(true)
                .help("The organisation ID"))
            .arg(Arg::with_name("name")
                .required(true)
                .takes_value(true)
                .help("The name parameter"))
            .arg(Arg::with_name("organisationpassword")
                .required(true)
                .takes_value(true)
                .help("The organisation password parameter")))
        .get_matches();

    if let Some(matches) = matches.subcommand_matches("register") {
        let organisation_id = matches.value_of("organisationID").unwrap();
        let name = matches.value_of("name").unwrap();
        let organisation_password = matches.value_of("organisationpassword").unwrap();

        // Read device identifier from the file
        let device_identifier = fs::read_to_string("/root/somefile")
            .expect("Failed to read device identifier from file");

        // Create the request body as JSON
        let request_body = json!({
            "name": name,
            "organisationpassword": organisation_password,
            "deviceidentifier": device_identifier.trim(),
        });

        // Serialize the JSON data into a string
        let request_body_str = serde_json::to_string(&request_body)
            .expect("Failed to serialize request body to JSON");

        // Create headers with required content-type and accept
        let mut headers = HeaderMap::new();
        headers.insert(CONTENT_TYPE, "application/json".parse().unwrap());
        headers.insert(ACCEPT, "application/json".parse().unwrap());

        // Send the POST request
        let url = format!("{}/api/devices/register/{}", BASE_URL, organisation_id);
        let client = reqwest::blocking::Client::new();
        let response = client.post(&url)
            .headers(headers)
            .body(request_body_str) // Set the serialized JSON data as the request body
            .send()
            .expect("Failed to send request");

        // Check if the response is successful
        if response.status().is_success() {
            let api_response: ApiResponse = response.json().expect("Failed to deserialize response");
            println!("Response: {:?}", api_response);
        } else {
            let text = response.text().expect("Failed to read response text");
            let api_error: ApiError = serde_json::from_str(&text).expect("Failed to deserialize error response");
            println!("Error: {:?}", api_error.message);
        }
    }
}
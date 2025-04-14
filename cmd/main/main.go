package main

import (
	"app"
	"app/internal/application"
	"encoding/json"
	"fmt"
	"os"
	"time"
)

func init() {
	app.Import()
}

func main() {

	defer func() {
		if err := recover(); err != nil {
			fmt.Println("Error:", err)
		}

		// print stack trace
		println(app.ColorRed)
		println("-----------------------------------------------------------------------")
		println("End of program!")
		println(app.ColorReset)
	}()

	fmt.Println("Welcome to the Command Line Interface of Application\n" +
		"\r\nSelect a command:" +
		"\r\n1. Start the application with a config file" +
		"\r\n2. Do no thing",
	)

	var input string
	fmt.Print("\r\nPlease enter a command: ")
	fmt.Scanln(&input)
	// TODO: add a switch case for the input
	switch input {
	case "1":
		fmt.Printf("-> Execute a command: (%s)\n", input)
		appN1()
		return
	case "2":
		fmt.Printf("-> Execute a command: (%s)\n", input)
		return
	default:
		fmt.Printf("-> Command `%s` is not defined.\nPlease enter a command from the list (1 >> 2).\n", input)
	}
}

func appN1() {

	type Config struct {
		Method             string                 `json:"method"`
		URL                string                 `json:"url"`
		MaxLifeTime        string                 `json:"maxLifeTime"`
		RequestTimeout     string                 `json:"requestTimeout"`
		InsecureSkipVerify bool                   `json:"insecureSkipVerify"`
		Headers            map[string]string      `json:"headers"`
		Body               map[string]interface{} `json:"body"`
	}

	var (
		inputWithJSON string
		tps           int

		config = Config{
			Method:             "GET",
			URL:                "https://www.google.com",
			MaxLifeTime:        "3s",
			RequestTimeout:     "15s",
			InsecureSkipVerify: false,
			Headers:            nil,
			Body:               nil,
		}
	)

	// read user input for json path
	fmt.Print("Enter the path to file of JSON config (default is `config.json`): ")
	fmt.Scanln(&inputWithJSON)
	if inputWithJSON == "" {
		inputWithJSON = "config.json"
	}

	// read json config
	if b, err := os.ReadFile(inputWithJSON); err != nil {
		// if file not found, create a new one
		if os.IsNotExist(err) {
			fmt.Println("File `config.json` not found, please create it first.")
			// marshal default config
			data, _ := json.Marshal(config)
			// create a new file
			if err := os.WriteFile(inputWithJSON, data, 0644); err != nil {
				fmt.Println("Error:", err)
			}
		} else {
			fmt.Println("Error:", err)
			return
		}
	} else {
		// parse json config
		if err := json.Unmarshal(b, &config); err != nil {
			fmt.Println("Error:", err)
		}
	}

	// read user input for tps
	fmt.Print("1000 Transactions Per Minute equals 1000/60 = 16.67 TPS; Please enter the TPS (default is 5): ")
	fmt.Scanln(&tps)

	// with config
	impl := application.New(config.Method, config.URL)

	// set max life time
	if d, err := time.ParseDuration(config.MaxLifeTime); err != nil {
		fmt.Println("Error:", err)
	} else {
		impl.SetMaxLifeTime(d)
	}
	// set request timeout
	if d, err := time.ParseDuration(config.RequestTimeout); err != nil {
		fmt.Println("Error:", err)
	} else {
		impl.SetRequestTimeout(d)
	}
	// set insecure skip verify
	if config.InsecureSkipVerify {
		impl.SetInsecureSkipVerify(config.InsecureSkipVerify)
	}
	// set headers
	if config.Headers != nil {
		if err := impl.SetHeaders(config.Headers); err != nil {
			fmt.Println("Error:", err)
		}
	}
	// set body if json config has body
	if config.Body != nil {
		// marshal body to json
		data, _ := json.Marshal(config.Body)
		// set body
		if err := impl.SetBody(data); err != nil {
			fmt.Println("Error:", err)
		}
		println("Body is already set:", string(data))
	}

	if tps == 0 {
		// set default tps is 5 transactions per second
		tps = 5
	}

	// run with tps custom
	if err := impl.RunWithTPS(tps); err != nil {
		fmt.Println("Error:", err)
	}
}

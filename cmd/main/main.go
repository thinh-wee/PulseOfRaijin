package main

import (
	"app"
	"app/internal/application"
	"encoding/json"
	"flag"
	"fmt"
	"net/http"
	"os"
	"time"
)

var (
	setConfigPath                     string
	setMaxLifeTime, setRequestTimeout string
	setMethod, setURL                 string
	setInsecureSkipVerify             = false
)

type Config struct {
	Method             string                 `json:"method"`
	URL                string                 `json:"url"`
	MaxLifeTime        string                 `json:"maxLifeTime"`
	RequestTimeout     string                 `json:"requestTimeout"`
	InsecureSkipVerify bool                   `json:"insecureSkipVerify"`
	Headers            map[string]string      `json:"headers"`
	Body               map[string]interface{} `json:"body"`
}

func (cfg *Config) WithFlag() (ok bool) {
	if setMaxLifeTime != "" {
		// set max life time
		if _, err := time.ParseDuration(setMaxLifeTime); err != nil {
			fmt.Printf("flag --life-time invalid."+
				" A duration string is a possibly signed sequence of decimal numbers, each with optional fraction and a unit suffix,"+
				" such as \"300ms\", \"-1.5h\" or \"2h45m\"."+
				" Valid time units are \"ns\", \"us\" (or \"µs\"), \"ms\", \"s\", \"m\", \"h\".\n%s\n", err)
		} else {
			cfg.MaxLifeTime = setMaxLifeTime
			ok = true
		}
	}
	if setRequestTimeout != "" {
		if _, err := time.ParseDuration(setRequestTimeout); err != nil {
			fmt.Printf("flag --timeout invalid."+
				" A duration string is a possibly signed sequence of decimal numbers, each with optional fraction and a unit suffix,"+
				" such as \"300ms\", \"-1.5h\" or \"2h45m\"."+
				" Valid time units are \"ns\", \"us\" (or \"µs\"), \"ms\", \"s\", \"m\", \"h\".\n%s\n", err)
		} else {
			cfg.RequestTimeout = setRequestTimeout
			ok = true
		}
	}
	if setMethod != "" {
		cfg.Method = setMethod
		ok = true
	}
	if setURL != "" {
		cfg.URL = setURL
		ok = true
	}
	return
}

func loadConfig() (*Config, error) {
	var (
		inputWithJSON string
		config        = Config{
			Method:             http.MethodPost,
			URL:                "http://localhost:8080/healthcheck",
			RequestTimeout:     "30s",
			MaxLifeTime:        "10s",
			InsecureSkipVerify: false,
			Headers:            nil,
			Body:               nil,
		}
	)

	if setConfigPath != "" {
		inputWithJSON = setConfigPath
	} else {
		// read user input for json path
		fmt.Print("Enter the path to file of JSON config (default is `config.json`): ")
		fmt.Scanln(&inputWithJSON)
		if inputWithJSON == "" {
			inputWithJSON = "config.json"
		}
	}

	// read json config
	if b, err := os.ReadFile(inputWithJSON); err != nil {
		// if file not found, create a new one
		if os.IsNotExist(err) {
			fmt.Printf("File `%s` not found, please create it first.\n", inputWithJSON)
			if ok := config.WithFlag(); ok {
				fmt.Println("Configuration applied by flags")
			}
			// marshal default config
			data, _ := json.Marshal(config)
			// create a new file
			if err := os.WriteFile(inputWithJSON, data, 0644); err != nil {
				fmt.Println("Error:", err)
			}
		} else {
			return nil, err
		}
	} else {
		// parse json config
		if err := json.Unmarshal(b, &config); err != nil {
			fmt.Println("Error:", err)
		}
		if ok := config.WithFlag(); ok {
			fmt.Println("Configuration applied by flags")
		}
	}
	return &config, nil
}

func init() {
	flag.StringVar(&setConfigPath, "config", "", "set path to configuration file")
	flag.StringVar(&setMaxLifeTime, "life-time", "", "set max life time")
	flag.StringVar(&setRequestTimeout, "timeout", "", "set timeout to request")
	flag.StringVar(&setMethod, "method", "", "set method")
	flag.StringVar(&setURL, "url", "", "set URL endpoint")
	flag.BoolVar(&setInsecureSkipVerify, "insecure-skip", false, "skip insecure verify SSL")
	flag.Parse()

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
		"\r\n1. Start the application with mode Pulse N per Minute (default)" +
		"\r\n2. Start the application with mode Pulse N per Second" +
		"\r\n3. Do no thing",
	)

	var input string
	fmt.Print("\r\nPlease enter a command: ")
	fmt.Scanln(&input)
	// add a switch case for the input
	if input == "" {
		input = "1" // set default command by (1)
	}
	//
	switch input {
	case "1":
		fmt.Printf("-> Execute a command: (%s) - Pulse Mode: N per Minute\n", input)
		appPulseNperMinute()
		return
	case "2":
		fmt.Printf("-> Execute a command: (%s) - Pulse Mode: N per Second\n", input)
		appPulseNperSecond()
		return
	default:
		fmt.Printf("-> Command `%s` is not defined.\nPlease enter a command from the list (1 >> 2).\n", input)
	}
}

func appPulseNperMinute() {

	var (
		tpm int
	)
	config, err := loadConfig()
	if err != nil {
		fmt.Println("Error:", err)
		return
	}

	// read user input for tps
	fmt.Print("Please enter the TPM (minimum is 60, default is 300): ")
	fmt.Scanln(&tpm)

	// with config
	impl := application.New(config.Method, config.URL)

	// set max life time
	if d, err := time.ParseDuration(config.MaxLifeTime); err != nil {
		fmt.Println("Error:", err)
	} else if err := impl.SetMaxLifeTime(d); err != nil {
		fmt.Println("Error:", err)
	}
	// set request timeout
	if d, err := time.ParseDuration(config.RequestTimeout); err != nil {
		fmt.Println("Error:", err)
	} else if err := impl.SetRequestTimeout(d); err != nil {
		fmt.Println("Error:", err)
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

	if tpm < 60 {
		// set default tpm is 300 transactions per minute ~ 5 transactions per second
		tpm = 300
	}

	// run with tps custom
	if err := impl.RunWithTPS(tpm / 60); err != nil {
		fmt.Println("Error:", err)
	}
}

func appPulseNperSecond() {

	var (
		tps int
	)
	config, err := loadConfig()
	if err != nil {
		fmt.Println("Error:", err)
		return
	}

	// read user input for tps
	fmt.Print("1000 Transactions Per Minute equals 1000/60 = 16.67 TPS; Please enter the TPS (default is 5): ")
	fmt.Scanln(&tps)

	// with config
	impl := application.New(config.Method, config.URL)

	// set max life time
	if d, err := time.ParseDuration(config.MaxLifeTime); err != nil {
		fmt.Println("Error:", err)
	} else if err := impl.SetMaxLifeTime(d); err != nil {
		fmt.Println("Error:", err)
	}
	// set request timeout
	if d, err := time.ParseDuration(config.RequestTimeout); err != nil {
		fmt.Println("Error:", err)
	} else if err := impl.SetRequestTimeout(d); err != nil {
		fmt.Println("Error:", err)
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

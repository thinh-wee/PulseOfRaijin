package main

import (
	"encoding/json"
	"fmt"
	"net/http"
	"os"
	"time"
)

type Config struct {
	Method             string            `json:"method"`
	URL                string            `json:"url"`
	MaxLifeTime        string            `json:"maxLifeTime"`
	RequestTimeout     string            `json:"requestTimeout"`
	InsecureSkipVerify bool              `json:"insecureSkipVerify"`
	Verbose            bool              `json:"verbose"`
	Headers            map[string]string `json:"headers"`
	Body               any               `json:"body"`
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
	if setInsecureSkipVerify {
		cfg.InsecureSkipVerify = true
	}
	if setVerbose {
		cfg.Verbose = true
	}
	if setMethod != "" {
		cfg.Method = setMethod
		ok = true
	}
	if setURL != "" {
		cfg.URL = setURL
		ok = true
	}
	if setReadFile != "" {
		//
		if b, err := os.ReadFile(setReadFile); err != nil {
			fmt.Printf("Read file error: %+v\n", err.Error())
			os.Exit(1)
		} else {
			cfg.Body = b
		}
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

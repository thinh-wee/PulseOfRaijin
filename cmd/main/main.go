package main

import (
	"app"
	"app/internal/application"
	"encoding/json"
	"flag"
	"fmt"
	"time"
)

var (
	setConfigPath                     string
	setMaxLifeTime, setRequestTimeout string
	setMethod, setURL                 string
	setReadFile                       string
	setInsecureSkipVerify             = false
	setVerbose                        = false
)

func init() {
	flag.StringVar(&setConfigPath, "config", "", "set path to configuration file")
	flag.StringVar(&setMaxLifeTime, "life-time", "", "set max life time")
	flag.StringVar(&setRequestTimeout, "timeout", "", "set timeout to request")
	flag.StringVar(&setMethod, "method", "", "set method")
	flag.StringVar(&setURL, "url", "", "set URL endpoint")
	flag.StringVar(&setReadFile, "file", "", "set file to body")
	flag.BoolVar(&setInsecureSkipVerify, "insecure-skip", false, "skip insecure verify SSL")
	flag.BoolVar(&setVerbose, "verbose", false, "print response")
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
	impl := application.New(config.Method, config.URL, config.Verbose)

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
	switch b := config.Body.(type) {
	case nil:
		break
	case []byte:
		// set body
		if err := impl.SetBody(b); err != nil {
			fmt.Println("Error:", err)
		}
		println("Body is already set:", len(b), "bytes")
	default:
		// marshal body to json
		data, _ := json.Marshal(b)
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
	impl := application.New(config.Method, config.URL, config.Verbose)

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
	switch b := config.Body.(type) {
	case nil:
		break
	case []byte:
		// set body
		if err := impl.SetBody(b); err != nil {
			fmt.Println("Error:", err)
		}
		println("Body is already set:", len(b), "bytes")
	default:
		// marshal body to json
		data, _ := json.Marshal(b)
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

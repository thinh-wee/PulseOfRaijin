package app

import (
	"fmt"
	"os"
	"path/filepath"
)

const (
	OriginalURL = "https://github.com/thinh-wee/pulse-of-raijin"
	LicenseURL  = "https://github.com/thinh-wee/pulse-of-raijin/blob/main/LICENSE"

	Author  = "thinh-wee"
	License = "MIT"

	ColorReset   = "\033[0m"
	ColorRed     = "\033[31m"
	ColorGreen   = "\033[32m"
	ColorYellow  = "\033[33m"
	ColorBlue    = "\033[34m"
	ColorMagenta = "\033[35m"
	ColorCyan    = "\033[36m"
)

var (
	Version   = "0.0.1"
	BuildDate = "2025-04-12T00:00:00Z"
	BuildUser = "Unknown"

	RepoURL     = "https://github.com/thinh-wee/pulse-of-raijin"
	LicenseText = `
	MIT License

	Copyright (c) 2025 thinh-wee

	Permission is hereby granted, free of charge, to any person obtaining a copy
	of this software and associated documentation files (the "Software"), to deal
	in the Software without restriction, including without limitation the rights
	to use, copy, modify, merge, publish, distribute, sublicense, and/or sell
	copies of the Software, and to permit persons to whom the Software is
	furnished to do so, subject to the following conditions:

	The above copyright notice and this permission notice shall be included in all
	copies or substantial portions of the Software.

	THE SOFTWARE IS PROVIDED "AS IS", WITHOUT WARRANTY OF ANY KIND, EXPRESS OR
	IMPLIED, INCLUDING BUT NOT LIMITED TO THE WARRANTIES OF MERCHANTABILITY,
	FITNESS FOR A PARTICULAR PURPOSE AND NONINFRINGEMENT. IN NO EVENT SHALL THE
	AUTHORS OR COPYRIGHT HOLDERS BE LIABLE FOR ANY CLAIM, DAMAGES OR OTHER
	LIABILITY, WHETHER IN AN ACTION OF CONTRACT, TORT OR OTHERWISE, ARISING FROM,
	OUT OF OR IN CONNECTION WITH THE SOFTWARE OR THE USE OR OTHER DEALINGS IN THE
	SOFTWARE.
	`
)

func init() {
	println(ColorYellow)
	println("=== PulseOfRaijin ===")
	println("Author:", Author)
	println("License:", LicenseURL)
	println("URL:", RepoURL)
	if RepoURL != OriginalURL {
		println("\r\t"+
			"Warning: The repository URL is different from the original URL",
			fmt.Sprintf("(%s).", OriginalURL))
		println("\r\n" + LicenseText)
	}
	println("-----------------------------------------------------------------------")
	println(ColorGreen)
	println(">> Version:", Version)
	println(">> Build by", BuildUser)
	println(">> Build at", BuildDate)
	println("=======================================================================")
	println(ColorReset)
}

func Import() {
	if _, err := os.Stat("LICENSE"); os.IsNotExist(err) {
		if err := os.WriteFile("LICENSE", []byte(LicenseText), 0644); err != nil {
			println("Error: Failed to create LICENSE file.")
		}
	}
}

func GetFullPathFromRoot(path string) string {
	dir, _ := os.Getwd()
	return filepath.Join(dir, path)
}

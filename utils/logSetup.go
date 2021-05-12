package utils

import (
	"strconv"

	"github.com/fatih/color"
)

// Log the env variables for the reverse proxy
func LogSetup(args []string) {
	color.Cyan("Server will run on: %s\n")

	for i := 0; i < len(args); i++ {
		color.Magenta("Redirecting to " + strconv.Itoa(i+1) + " url: %s\n", args[i])
	}
}
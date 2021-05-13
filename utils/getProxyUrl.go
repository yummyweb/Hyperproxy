package utils

import (
	"strings"
	"strconv"
)

// Get the url for a given proxy condition
func GetProxyUrl(args []string, proxyConditionRaw string) string {
	proxyCondition := strings.ToUpper(proxyConditionRaw)

	// url_a := GetEnv("URL_A", "http://localhost:1331")
	// url_b := GetEnv("URL_B", "http://localhost:1332")
	// default_url := GetEnv("DEFAULT_URL", "http://localhost:1333")

	for i := 0; i < len(args); i++ {
		if proxyCondition == strconv.Itoa(i) {
			return args[i]
		}
	}

	// TODO: Make this dynamic by taking --default flag input
	return "http://localhost:1333"
}
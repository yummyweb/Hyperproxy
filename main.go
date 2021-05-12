package main

import (
	"bytes"
	"encoding/json"
	"io/ioutil"
	"log"
	"net/http"
	"net/http/httputil"
	"net/url"
	"os"
	"strings"

	"github.com/fatih/color"
)

// Get env var or use default value
func getEnv(key, fallbackValue string) string {
	if value, ok := os.LookupEnv(key); ok {
		return value
	}
	return fallbackValue
}

// Get the port to listen on
func getListenAddress() string {
	port := getEnv("PORT", "1338")
	return ":" + port
}

// Log the env variables for the reverse proxy
func logSetup() {
	url_a := getEnv("URL_A", "http://localhost:1331")
	url_b := getEnv("URL_B", "http://localhost:1332")
	default_url := getEnv("DEFAULT_URL", "http://localhost:1333")

	color.Cyan("Server will run on: %s\n", getListenAddress())
	color.Magenta("Redirecting to A url: %s\n", url_a)
	color.Magenta("Redirecting to B url: %s\n", url_b)
	color.Magenta("Redirecting to Default url: %s\n", default_url)
}

type requestPayloadStruct struct {
	ProxyCondition string `json:"proxy_condition"`
}

// Log the typeform payload and redirect url
func logRequestPayload(requestionPayload requestPayloadStruct, proxyUrl string) {
	log.Printf("proxy_condition: %s, proxy_url: %s\n", requestionPayload.ProxyCondition, proxyUrl)
}

// Get the url for a given proxy condition
func getProxyUrl(proxyConditionRaw string) string {
	proxyCondition := strings.ToUpper(proxyConditionRaw)

	url_a := getEnv("URL_A", "http://localhost:1331")
	url_b := getEnv("URL_B", "http://localhost:1332")
	default_url := getEnv("DEFAULT_URL", "http://localhost:1333")

	if proxyCondition == "A" {
		return url_a
	}

	if proxyCondition == "B" {
		return url_b
	}

	return default_url
}

// Serve a reverse proxy for a given url
func serveReverseProxy(target string, res http.ResponseWriter, req *http.Request) {
	// Parse the url
	url, _ := url.Parse(target)

	// create the reverse proxy
	proxy := httputil.NewSingleHostReverseProxy(url)

	// Update the headers to allow for SSL redirection
	req.URL.Host = url.Host
	req.URL.Scheme = url.Scheme
	req.Header.Set("X-Forwarded-Host", req.Header.Get("Host"))
	req.Host = url.Host

	// Note that ServeHttp is non blocking and uses a go routine under the hood
	proxy.ServeHTTP(res, req)
}

// Get a json decoder for a given requests body
func requestBodyDecoder(request *http.Request) *json.Decoder {
	// Read body to buffer
	body, err := ioutil.ReadAll(request.Body)
	if err != nil {
		log.Printf("Error reading body: %v", err)
		panic(err)
	}

	request.Body = ioutil.NopCloser(bytes.NewBuffer(body))

	return json.NewDecoder(ioutil.NopCloser(bytes.NewBuffer(body)))
}

// Parse the requests body
func parseRequestBody(request *http.Request) requestPayloadStruct {
	decoder := requestBodyDecoder(request)

	var requestPayload requestPayloadStruct
	err := decoder.Decode(&requestPayload)

	if err != nil {
		panic(err)
	}

	return requestPayload
}

// Redirect a request to the appropriate url
func handleRequestAndRedirect(res http.ResponseWriter, req *http.Request) {
	requestPayload := parseRequestBody(req)
	url := getProxyUrl(requestPayload.ProxyCondition)
	logRequestPayload(requestPayload, url)

	serveReverseProxy(url, res, req)
}

func main() {
	// Log setup values
	logSetup()

	// Start proxy server
	http.HandleFunc("/", handleRequestAndRedirect)
	if err := http.ListenAndServe(getListenAddress(), nil); err != nil {
		panic(err)
	}
}

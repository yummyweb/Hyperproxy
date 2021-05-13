package cmd

import (
	"bytes"
	"encoding/json"
	"io/ioutil"
	"log"
	"net/http"
	"net/http/httputil"
	"net/url"
	"strconv"

	"github.com/yummyweb/Hyperproxy/utils"

	"github.com/spf13/cobra"
)

var port int

func init() {
	startCmd.Flags().IntVarP(&port, "port", "p", 1338, "the port to run the proxy server")

	rootCmd.AddCommand(startCmd)
}

var startCmd = &cobra.Command{
	Use:   "start",
	Short: "Starts proxy server",
	Long: `start is a sub-command for hyperproxy.
	It starts the reverse proxy server at specified port,
	address and other proxy urls. If no flags are specified, the port will be 1338
	Usage:
		hyperproxy start http://localhost:1331 http://localhost:1332
		hyperproxy start -p 4090 http://localhost:5454
	`,
	Args: cobra.MinimumNArgs(1),
	Run:  executeStartCmd,
}

// Get the port to listen on
func getListenAddress(port int) string {
	return ":" + strconv.Itoa(port)
}

type requestPayloadStruct struct {
	ProxyCondition string `json:"proxy_condition"`
}

// Log the typeform payload and redirect url
func logRequestPayload(requestionPayload requestPayloadStruct, proxyUrl string) {
	log.Printf("proxy_condition: %s, proxy_url: %s\n", requestionPayload.ProxyCondition, proxyUrl)
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
func handleRequestAndRedirect(args []string, res http.ResponseWriter, req *http.Request) {
	requestPayload := parseRequestBody(req)
	url := utils.GetProxyUrl(args, requestPayload.ProxyCondition)
	logRequestPayload(requestPayload, url)

	serveReverseProxy(url, res, req)
}

func executeStartCmd(cmd *cobra.Command, args []string) {
	// Log setup values
	utils.LogSetup(args)

	// Start proxy server
	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		handleRequestAndRedirect(args, w, r)
	})
	if err := http.ListenAndServe(getListenAddress(port), nil); err != nil {
		panic(err)
	}
}

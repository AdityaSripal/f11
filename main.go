package main

import (
	"encoding/json"
	"flag"
	"fmt"
	"github.com/aws/aws-lambda-go/events"
	"github.com/aws/aws-lambda-go/lambda"
	"github.com/awslabs/aws-lambda-go-api-proxy/gorillamux"
	sdkversion "github.com/cosmos/cosmos-sdk/version"
	"github.com/gorilla/mux"
	"github.com/greg-szabo/f11/defaults"
	"github.com/greg-szabo/f11/endpoints"
	tendermintversion "github.com/tendermint/tendermint/version"
	"log"
	"net/http"
	"os"
	"time"
)

var initialized = false
var muxLambda *gorillamux.GorillaMuxAdapter

type MainMessage struct {
	Message string `json:"message"`
	Name    string `json:"name"`
	Version string `json:"version"`
}

func MainHandler(w http.ResponseWriter, r *http.Request) {
	defaults.Headers(w)
	defaults.AddStatusOK(w)
	json.NewEncoder(w).Encode(MainMessage{"", defaults.TestnetName, defaults.Version})
}

func loggingMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		log.Println(r.RequestURI)
		// Call the next handler, which can be another middleware in the chain, or the final handler.
		next.ServeHTTP(w, r)
	})
}

func AddRoutes() (r *mux.Router) {

	// Root
	r = mux.NewRouter()
	r.HandleFunc("/", MainHandler)

	// Routes
	endpoints.AddRoutesV1(r)

	// Finally
	r.Use(loggingMiddleware)
	http.Handle("/", r)

	return
}

func Initialization() {
	endpoints.InitializeV1()
}

func LambdaHandler(req events.APIGatewayProxyRequest) (events.APIGatewayProxyResponse, error) {

	if !initialized {
		// stdout and stderr are sent to AWS CloudWatch Logs
		log.Println("Cold start")
		r := AddRoutes()

		muxLambda = gorillamux.New(r)
		initialized = true
	}

	//Todo: Add lambda timeout function so a response is made before the function times out in AWS

	return muxLambda.Proxy(req)

}

func LocalExecution() {
	log.Println("Local Execution start")

	r := AddRoutes()

	srv := &http.Server{
		Addr: "127.0.0.1:3000",
		// Good practice to set timeouts to avoid Slowloris attacks.
		WriteTimeout: time.Second * 15,
		ReadTimeout:  time.Second * 15,
		IdleTimeout:  time.Second * 60,
		Handler:      r, // Pass our instance of gorilla/mux in.
	}

	// Run our server in a goroutine so that it doesn't block.
	if err := srv.ListenAndServe(); err != nil {
		log.Println(err)
	}
}

func main() {
	var localExecution bool
	var versionSwitch bool

	flag.BoolVar(&versionSwitch, "version", false, "Return version number and exit.")
	flag.BoolVar(&localExecution, "local", false, "run a local web-server instead of as an AWS Lambda function - for development and troubleshooting purposes")
	flag.Parse()

	if versionSwitch {
		fmt.Println(defaults.Version)
		fmt.Printf("Testnet: %s\n", defaults.TestnetName)
		fmt.Printf("SDK: %v\n", sdkversion.Version)
		fmt.Printf("Tendermint: %v\n", tendermintversion.Version)
		os.Exit(0)
	}

	if localExecution {
		LocalExecution()
	} else {
		lambda.Start(LambdaHandler)
	}

}

package main

import (
	"context"
	"encoding/json"
	"flag"
	"fmt"
	"github.com/aws/aws-lambda-go/events"
	"github.com/aws/aws-lambda-go/lambda"
	"github.com/awslabs/aws-lambda-go-api-proxy/gorillamux"
	"github.com/gorilla/mux"
	"github.com/greg-szabo/f11/defaults"
	"github.com/greg-szabo/f11/endpoints"
	"github.com/greg-szabo/f11/testnet"
	"github.com/greg-szabo/f11/version"
	"log"
	"net/http"
	"os"
	"os/signal"
	"time"
)

var initialized = false
var muxLambda *gorillamux.GorillaMuxAdapter

type MainMessage struct {
	Name    string `json:"name"`
	Version string `json:"version"`
}

func MainHandler(w http.ResponseWriter, r *http.Request) {
	defaults.Headers(w)
	json.NewEncoder(w).Encode(MainMessage{testnet.TestnetName, version.Version})
}

func loggingMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		log.Println(r.RequestURI)
		// Call the next handler, which can be another middleware in the chain, or the final handler.
		next.ServeHTTP(w, r)
	})
}

func throttlingMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		log.Println("Throttle-throttle")
		//Todo: implement throttling
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
	r.Use(throttlingMiddleware)
	http.Handle("/", r)

	return
}

func LambdaHandler(req events.APIGatewayProxyRequest) (events.APIGatewayProxyResponse, error) {

	if !initialized {
		// stdout and stderr are sent to AWS CloudWatch Logs
		log.Println("GorillaMux cold start")
		r := AddRoutes()

		muxLambda = gorillamux.New(r)
		initialized = true
	}

	//Todo: Add lambda timeout function so a response is made before the function stops

	return muxLambda.Proxy(req)

}

func main() {
	var timeout time.Duration
	var localExecution bool
	var versionSwitch bool

	flag.BoolVar(&versionSwitch, "version", false, "Return version number and exit.")
	flag.BoolVar(&localExecution, "local", false, "run a local web-server instead of as an AWS Lambda function - for development and troubleshooting purposes")
	flag.DurationVar(&timeout, "timeout", time.Second*3, "the duration for which the server gracefully timeout for existing connections to finish - e.g. 15s or 1m")
	flag.Parse()

	if versionSwitch {
		fmt.Println(version.Version)
		fmt.Println(testnet.TestnetName)
		os.Exit(0)
	}

	if localExecution {
		r := AddRoutes()
		srv := &http.Server{
			Addr: "0.0.0.0:3000",
			// Good practice to set timeouts to avoid Slowloris attacks.
			WriteTimeout: time.Second * 15,
			ReadTimeout:  time.Second * 15,
			IdleTimeout:  time.Second * 60,
			Handler:      r, // Pass our instance of gorilla/mux in.
		}

		// Run our server in a goroutine so that it doesn't block.
		go func() {
			if err := srv.ListenAndServe(); err != nil {
				log.Println(err)
			}
		}()

		c := make(chan os.Signal, 1)
		// We'll accept graceful shutdowns when quit via SIGINT (Ctrl+C)
		// SIGKILL, SIGQUIT or SIGTERM (Ctrl+/) will not be caught.
		signal.Notify(c, os.Interrupt)

		// Block until we receive our signal.
		<-c

		// Create a deadline to timeout for.
		ctx, cancel := context.WithTimeout(context.Background(), timeout)
		defer cancel()
		// Doesn't block if no connections, but will otherwise timeout
		// until the timeout deadline.
		srv.Shutdown(ctx)
		// Optionally, you could run srv.Shutdown in a goroutine and block on
		// <-ctx.Done() if your application should timeout for other services
		// to finalize based on context cancellation.
		log.Println("shutting down")
		os.Exit(0)
	}

	lambda.Start(LambdaHandler)

}

package main

import (
	"log/slog"
	"net/http"
	"os"

	"github.com/arkag/cope-and-hope-weather/cache"
	"github.com/arkag/cope-and-hope-weather/client"
	"github.com/arkag/cope-and-hope-weather/handler"

	"github.com/aws/aws-lambda-go/lambda"
	"github.com/awslabs/aws-lambda-go-api-proxy/httpadapter"
)

func main() {
	opts := &slog.HandlerOptions{
		Level: slog.LevelDebug,
	}
	logger := slog.New(slog.NewJSONHandler(os.Stdout, opts))
	slog.SetDefault(logger)

	apiKey := os.Getenv("OWM_API_KEY")
	if apiKey == "" {
		slog.Error("missing environment variable", "env_var", "OWM_API_KEY")
		os.Exit(1)
	}

	realClient := client.Client{
		APIKey:  apiKey,
		BaseURL: "https://api.openweathermap.org",
	}

	s := &handler.Server{
		WeatherClient: &cache.CachedClient{
			RealFetcher: realClient,
			TableName:   "WeatherCache",
			Dynamo:      nil, // We will initialize the AWS client later in Phase 2
		},
	}

	http.HandleFunc("/weather", s.HandleWeather)

	if os.Getenv("AWS_LAMBDA_FUNCTION_NAME") != "" {
		slog.Info("starting AWS Lambda handler")
		adapter := httpadapter.New(http.DefaultServeMux)
		lambda.Start(adapter.ProxyWithContext)
	} else {
		slog.Info("starting local server", "port", 8080)
		if err := http.ListenAndServe(":8080", nil); err != nil {
			slog.Error("server crashed", "error", err)
			os.Exit(1)
		}
	}
}

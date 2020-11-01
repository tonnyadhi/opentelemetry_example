package main

import (
	"context"
	"encoding/json"
	"io/ioutil"
	"net/http"
	"os"

	"go.opentelemetry.io/contrib/instrumentation/net/http/httptrace/otelhttptrace"
	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/api/global"
	"go.opentelemetry.io/otel/label"
	"go.opentelemetry.io/otel/propagators"
)

//RequestWeatherForecast represent return from weatherclient service / adapter
type RequestWeatherForecast struct {
	Condition   string  `json:"Condition"`
	Temperature float64 `json:"Temperature"`
	Humidity    int32   `json:"Humidity"`
}

//GetWeatherForecast return WeatherForecast
func getWeatherForecast(ctx context.Context, city string) (*RequestWeatherForecast, error) {
	client := http.DefaultClient
	adapterHost := os.Getenv("ADAPTER_HOST")
	requestPath := "getweather/owm"
	requestParam := city

	tracer := global.Tracer("GetWeatherForecastFunction")

	global.SetTextMapPropagator(propagators.Baggage{})

	distributedctx := otel.ContextWithBaggageValues(ctx, label.Any("requestWeatherForecast", "OWM"))

	ctx, span := tracer.Start(distributedctx, "inside adapter for owm weather forecast")

	span.AddEvent(ctx, "Inside Forecast Function", label.Any("Forecast Function Called", "Forwarded to Adapter Server"))

	req, _ := http.NewRequest("GET", adapterHost+"/"+requestPath+"/"+requestParam, nil)
	ctx, req = otelhttptrace.W3C(ctx, req)
	otelhttptrace.Inject(ctx, req)

	defer span.End()

	res, err := client.Do(req)
	if err != nil {
		return nil, err
	}

	body, err := ioutil.ReadAll(res.Body)
	if err != nil {
		return nil, err
	}

	defer res.Body.Close()

	rwf, err := parseResponse(body)
	if err != nil {
		return nil, err
	}

	return rwf, nil

}

func parseResponse(body []byte) (*RequestWeatherForecast, error) {
	rwf := &RequestWeatherForecast{}

	err := json.Unmarshal(body, rwf)
	if err != nil {
		return nil, err
	}

	return rwf, nil

}

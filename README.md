# Open Telemetry Example in Golang

## Structure

   This example contains  two services 
   - HTTP Weather Service
   - HTTP Weahter Adapter Service, fetching weather forecast from [openweathermap](https://openweathermap.org)

## Stack

   - Golang v1.15
   - [OpenTelemetry Go](https://github.com/open-telemetry/opentelemetry-go) v0.13
   - [Go Chi](https://github.com/go-chi/chi) v4.1.2

## Running

## Running Without Kubernetes

   - Get yourself an API Key from openweathermap
   - Download and run [Jaeger All In One](https://hub.docker.com/r/jaegertracing/all-in-one) for tracing
     - `$./jaeger-all-in-one`
   - Run Weather Adapter Service
     - ` $export OWM_APP_ID=<OWM_APIKEY> && export JAEGER_HOST=http://localhost:14268/api/traces && go run weatherclient_adapter/openweathermap_svc.go`
   - Run Weather Service
     - `$export OWM_APP_ID=<OWM_APIKEY>  && export JAEGER_HOST=http://localhost:14268/api/traces && export ADAPTER_HOST=http://localhost:8181 && go run weather_service/weather_server.go weather_service/adapter_client.go`
   - Try to curl into Weather Service
     - `$curl localhost:9191/forecast/depok`
   - Go to Jaeger All In UI at port 16686 for observing traces
     - `$firefox localhost:16686`
    
## Running Inside Kubernetes
   - You can use provided helm chart on k8s directory
   - Assumed that you already installed Istio on your k8s for convenience on receiving tracing via Jaeger
   - Modify provided helm chart base on your need

    


package main

import (
	"fmt"
	"net/http"
	openweathermap "opentelemetry_example/openweathermap_client"
	tracing "opentelemetry_example/tracing"
	"os"

	"github.com/go-chi/chi"
	"github.com/go-chi/chi/middleware"
	"github.com/go-chi/render"
	"go.opentelemetry.io/contrib/instrumentation/net/http/httptrace/otelhttptrace"
	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/api/global"
	"go.opentelemetry.io/otel/api/trace"
	"go.opentelemetry.io/otel/label"
)

type strippedWeatherData struct {
	Condition   string
	Temperature float64
	Humidity    int
}

var (
	currentWeather *openweathermap.CurrentWeatherResponse
	err            error
)

func returnPing(w http.ResponseWriter, r *http.Request) {
	tracer := global.Tracer("WeatherAdapterService")
	ctx, span := tracer.Start(r.Context(), "ping route called")

	span.AddEvent(ctx, "Ping Route Called", label.Any("Ping Called", "GET /ping"))
	defer span.End()

	w.Write([]byte("pong"))
}

//Bind strippedWeatherData
func (swd *strippedWeatherData) Bind(r *http.Request) error {
	return nil
}

func (swd *strippedWeatherData) Render(w http.ResponseWriter, r *http.Request) error {
	return nil
}

func getWeatherByCity(w http.ResponseWriter, r *http.Request) {
	city := chi.URLParam(r, "city")
	cityWeather, err := getOwmForecastByCity(city)

	tracer := global.Tracer("OWM Function")
	attrs, entries, spanCtx := otelhttptrace.Extract(r.Context(), r)

	r = r.WithContext(otel.ContextWithBaggageValues(r.Context(), label.Any("requestWeatherForecast", "OWM")))

	ctx, span := tracer.Start(
		trace.ContextWithRemoteSpanContext(r.Context(), spanCtx),
		"owm forecast route called",
		trace.WithAttributes(attrs...),
		trace.WithAttributes(entries...),
	)

	span.AddEvent(ctx, "OWM Route Called", label.Any("OWM Forecast Called", "GET /forecast/owm/city"))

	if err != nil {
		w.WriteHeader(500)
		http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
		return
	}

	defer span.End()
	render.JSON(w, r, cityWeather)

}

func getOwmForecastByCity(city string) (*strippedWeatherData, error) {

	owm := openweathermap.OpenWeatherMap{APIKEY: os.Getenv("OWM_APP_ID")}
	currentWeather, err = owm.CurrentWeatherFromCity(city)

	if err != nil {
		return nil, err
	}

	swd := &strippedWeatherData{
		Condition:   currentWeather.Weather[0].Main,
		Temperature: currentWeather.Main.Temp,
		Humidity:    currentWeather.Main.Humidity,
	}

	return swd, nil

}

func httpTraceWrapper(h http.Handler) http.Handler {
	fn := func(w http.ResponseWriter, r *http.Request) {
		t := global.TracerProvider().Tracer("component-http")
		ctx, span := t.Start(r.Context(), r.URL.Path)
		r = r.WithContext(ctx)
		h.ServeHTTP(w, r)
		span.End()
	}
	return http.HandlerFunc(fn)
}

func main() {
	router := chi.NewRouter()

	tracerConf := &tracing.Config{
		ServiceName: "WeatherAdapterService",
		Endpoint:    os.Getenv("JAEGER_HOST"),
		Probability: 1.0,
	}

	jaeger := tracing.NewJaeger(tracerConf)

	defer jaeger()

	router.Use(middleware.Logger)
	router.Use(middleware.Recoverer)
	router.Use(render.SetContentType(render.ContentTypeJSON))

	router.Get("/ping", returnPing)
	router.Route("/getweather/owm", func(router chi.Router) {
		router.Get("/{city}", getWeatherByCity)
	})

	err := http.ListenAndServe(":8181", httpTraceWrapper(router))
	if err != nil {
		fmt.Println("ListenAndServe:", err)
	}

}

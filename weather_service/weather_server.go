package main

import (
	"fmt"
	"net/http"
	tracing "opentelemetry_example/tracing"
	"os"

	"github.com/go-chi/chi"
	"github.com/go-chi/chi/middleware"
	"github.com/go-chi/render"
	"go.opentelemetry.io/otel/api/global"
	"go.opentelemetry.io/otel/label"
)

//requestWeatherForecast represent return of WeatherForecast to service requester
type requestWeatherForecast struct {
	Condition   string  `json:"Condition"`
	Temperature float64 `json:"Temperature"`
	Humidity    int32   `json:"Humidity"`
}

//Bind RequestWeatherForecast
func (wf *requestWeatherForecast) Bind(r *http.Request) error {
	return nil
}

//Render requestWeatherForecast
func (wf *requestWeatherForecast) Render(w http.ResponseWriter, r *http.Request) error {
	return nil
}

func weatherForecast(w http.ResponseWriter, r *http.Request) {
	tracer := global.Tracer("WeatherForecastFunction")
	ctx, span := tracer.Start(r.Context(), "forecast route called")

	city := chi.URLParam(r, "city")
	wF, err := getWeatherForecast(r.Context(), city)

	span.AddEvent(ctx, "Forecast Route Called", label.Any("Forecast Called", "GET /forecast/city"))
	defer span.End()

	if err != nil {
		w.WriteHeader(500)
		http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
		return
	}

	render.JSON(w, r, wF)

}

func returnPing(w http.ResponseWriter, r *http.Request) {
	tracer := global.Tracer("PingFunction")
	ctx, span := tracer.Start(r.Context(), "ping route called")

	span.AddEvent(ctx, "Ping Route Called", label.Any("Ping Called", "GET /ping"))
	defer span.End()

	w.Write([]byte("pong"))
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
		ServiceName: "WeatherService",
		Endpoint:    os.Getenv("JAEGER_HOST"),
		Probability: 1.0,
	}

	jaeger := tracing.NewJaeger(tracerConf)

	defer jaeger()

	router.Use(middleware.Logger)
	router.Use(middleware.Recoverer)
	router.Use(render.SetContentType(render.ContentTypeJSON))

	router.Get("/ping", returnPing)

	router.Route("/forecast", func(router chi.Router) {
		router.Get("/{city}", weatherForecast)
	})

	err := http.ListenAndServe(":9191", httpTraceWrapper(router))
	if err != nil {
		fmt.Println("ListenAndServe:", err)
	}
}

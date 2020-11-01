module opentelemetry_example

go 1.15

require (
	github.com/ernesto-jimenez/httplogger v0.0.0-20150224132909-86cc44f6150a
	github.com/getsentry/sentry-go v0.7.0
	github.com/go-chi/chi v4.1.2+incompatible
	github.com/go-chi/render v1.0.1
	github.com/golang/protobuf v1.4.3
	go.opentelemetry.io/contrib/instrumentation/net/http/httptrace v0.11.0 // indirect
	go.opentelemetry.io/contrib/instrumentation/net/http/httptrace/otelhttptrace v0.13.0
	go.opentelemetry.io/otel v0.13.0
	go.opentelemetry.io/otel/exporters/trace/jaeger v0.13.0
	go.opentelemetry.io/otel/sdk v0.13.0
	golang.org/x/net v0.0.0-20201016165138-7b1cca2348c0 // indirect
	golang.org/x/sys v0.0.0-20201018230417-eeed37f84f13 // indirect
	google.golang.org/genproto v0.0.0-20201015140912-32ed001d685c // indirect
	google.golang.org/grpc v1.33.0
	google.golang.org/grpc/cmd/protoc-gen-go-grpc v1.0.0 // indirect
	google.golang.org/protobuf v1.25.0
)

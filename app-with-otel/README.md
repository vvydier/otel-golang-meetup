# shipping-service


## Run

Run the executable with the following required environment variables. Alternatively, you can run it as a dockerized container.
Access it a few times, so you can emit from telemetry

```cmd
  cd app
  export OTEL_EXPORTER_OTLP_ENDPOINT="0.0.0.0:4317"
  export OTEL_SERVICE_NAME="shipping-service"
  export OTEL_RESOURCE_ATTRIBUTES="deployment.environment=production"
  go run .
  curl http://127.0.0.1:9100/shipping/11
```

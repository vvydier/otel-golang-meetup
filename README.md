# otel-golang-meetup

## Instructions to get started and follow along

Run the executable with the following required environment variables

```cmd
  export OTEL_EXPORTER_OTLP_ENDPOINT = "0.0.0.0:4317"
  export OTEL_SERVICE_NAME = "shipping-service"
  export OTEL_RESOURCE_ATTRIBUTES = deployment.environment=production
  ```

# otel-golang-meetup
https://www.meetup.com/stl-go/events/293596247/

## The presentation deck with all the links and instructions
https://docs.google.com/presentation/d/1kElmZfLCjW8kcn5hS2y_Xx5HWIQKkeAWIBZNDV6ILUo

## Instructions to get started and follow along

### Opentelemetry collector, Prometheus and Zipkin running locally as a docker images

1. Run docker from otel-docker directory: This will start otel-collector, Zipkin and Prometheus

    ```shell script
    # from this directory
    docker-compose up
    ```

2. Run  app

    ```shell script
    # from the app directory
    go run .
    ```

3. Teardown the docker images

    ```shell script
    # from this directory
    docker-compose down
    ```

Additional instructions are at otel-docker/README.md

### Add instrumentation to your App and run

Run the app with the following required environment variables

```cmd
  export OTEL_EXPORTER_OTLP_ENDPOINT = "0.0.0.0:4317"
  export OTEL_SERVICE_NAME = "shipping-service"
  export OTEL_RESOURCE_ATTRIBUTES = deployment.environment=production
  ```

```cmd
  go run .
  curl http://127.0.0.1:9100/shipping/11
  ```
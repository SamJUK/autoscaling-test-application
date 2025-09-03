FROM golang:1.24 AS build-stage

WORKDIR /app

COPY go.mod go.sum ./
RUN go mod download

COPY *.go ./

RUN CGO_ENABLED=0 GOOS=linux go build -o /autoscaling-test-server

# Deploy the application binary into a lean image
FROM gcr.io/distroless/base-debian11 AS build-release-stage

WORKDIR /

COPY --from=build-stage /autoscaling-test-server /autoscaling-test-server

EXPOSE 80

USER nonroot:nonroot

ENTRYPOINT ["/autoscaling-test-server"]
FROM golang:1.21 AS build-stage
WORKDIR /app
COPY go.mod go.sum ./
RUN go mod download
COPY . ./
RUN go build -o /app/gateway-api ./cmd/server/gateway/main.go

FROM gcr.io/distroless/base-debian12 AS build-release-stage
WORKDIR /
COPY --from=build-stage /app/gateway-api /gateway-api
ENTRYPOINT ["/gateway-api"]
# ---------- build ----------
FROM golang:1.19-alpine AS build
WORKDIR /app

COPY go.mod go.sum ./
RUN go mod download
RUN go mod tidy
COPY . .

RUN CGO_ENABLED=0 GOOS=linux go build -o geoapi ./cmd/geoapi

# ---------- runtime ----------
FROM alpine:latest
WORKDIR /app
COPY --from=build /app/geoapi /app/geoapi

EXPOSE 8080
ENTRYPOINT ["/app/geoapi"]

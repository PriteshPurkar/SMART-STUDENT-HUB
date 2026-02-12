FROM golang:1.22-alpine AS build

WORKDIR /app
COPY go.mod go.sum ./
RUN go mod download

COPY . .
RUN CGO_ENABLED=0 GOOS=linux go build -a -installsuffix cgo -o api ./cmd/api

FROM alpine:3.20
RUN apk --no-cache add ca-certificates

WORKDIR /app
COPY --from=build /app/api /app/api

EXPOSE 8080
CMD ["/app/api"]


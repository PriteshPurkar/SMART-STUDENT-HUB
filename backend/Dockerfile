FROM golang:1.22-alpine AS build

WORKDIR /app
COPY . .
RUN go build -o api ./cmd/api

FROM alpine:3.20
WORKDIR /app
COPY --from=build /app/api /app/api
EXPOSE 8080
CMD ["/app/api"]


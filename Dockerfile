FROM golang:1.24 AS build-stage

WORKDIR /app

COPY go.mod go.sum ./

RUN go mod download

COPY . .

RUN CGO_ENABLED=0 GOOS=linux go build -o /app/main ./simple_server

# Add a test stage next

FROM alpine

WORKDIR /app

COPY --from=build-stage /app/main .

RUN adduser -D appuser
USER appuser

CMD ["./main"]

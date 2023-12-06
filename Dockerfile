# This file should be populated with the necessary configuration to build your web service's container image

FROM golang:alpine AS BUILDER

WORKDIR /app

COPY go.mod .

RUN go mod download

COPY . .

RUN CGO_ENABLED=0 GOOS=linux go build -o server .


FROM alpine AS RUNNER

WORKDIR /app

COPY --from=BUILDER /app/server /app/games.json .

EXPOSE 8080

CMD ["./server"]


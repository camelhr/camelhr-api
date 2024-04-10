FROM golang:1.22.2-alpine3.19 AS builder
WORKDIR /app
COPY . .
RUN go build -trimpath -o ./bin/api ./cmd/api

FROM alpine:3.19
WORKDIR /app
COPY --from=builder /app/bin/api .

CMD [ "/app/api"]
FROM golang:1.22.2-alpine3.19 AS builder
WORKDIR /app

# download dependencies into go's module cache
# this allows the cache to be used across builds
COPY go.mod go.sum ./
RUN go mod download -x

# copy the rest of the source code
COPY . .

# build the application
RUN go build -trimpath -o ./bin/api ./cmd/api
RUN go build -trimpath -o ./bin/dbmigrator ./cmd/dbmigrator

# create a new image with the binaries and required files
FROM alpine:3.19
WORKDIR /app
COPY --from=builder /app/bin/api .
COPY --from=builder /app/bin/dbmigrator .
COPY --from=builder /app/migrations ./migrations

CMD ["/app/api"]

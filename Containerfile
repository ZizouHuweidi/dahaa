FROM docker.io/library/golang:1.23-alpine AS build

WORKDIR /src
RUN apk add --no-cache git ca-certificates
COPY go.mod go.sum ./
RUN go mod download
COPY . .
RUN CGO_ENABLED=0 go build -o /out/server ./cmd/api

FROM docker.io/library/alpine:3.20
WORKDIR /app
RUN apk add --no-cache ca-certificates
COPY --from=build /out/server /app/server
EXPOSE 8080
CMD ["/app/server"]

FROM golang:1.25.4-alpine AS build
WORKDIR /app
COPY go.mod go.sum ./
RUN go mod download
COPY . .
RUN go build -o server ./cmd/server

FROM alpine:3.22.2
WORKDIR /app
COPY --from=build /app/server .
COPY --from=build /app/web ./web
EXPOSE 8080
CMD ["./server"]

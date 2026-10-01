FROM golang:1.24-alpine AS build
WORKDIR /src
COPY go.mod go.sum ./
RUN go mod download
COPY . .
RUN CGO_ENABLED=0 go build -trimpath -ldflags="-s -w" -o /out/nagorneva_reviews ./cmd/nagorneva_reviews

FROM alpine:3.20
RUN adduser -D -u 10001 app
USER app
COPY --from=build /out/nagorneva_reviews /usr/local/bin/nagorneva_reviews
EXPOSE 8080
ENTRYPOINT ["/usr/local/bin/nagorneva_reviews"]

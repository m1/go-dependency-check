FROM golang:1.13-alpine AS builder
ENV GOHOSTARCH="amd64"
ENV GOHOSTOS="linux"
ENV GOOS="linux"
ENV GO111MODULE=on

RUN apk add --no-cache build-base git ca-certificates
RUN git config --global core.autocrlf false

WORKDIR /src

# Download dependencies
COPY go.mod .
COPY go.sum .
RUN go mod download

# Copy everything else and build
COPY . .
RUN go build -ldflags "-linkmode external -extldflags -static" -a ./main.go

FROM scratch AS final

# Set to 0.0.0.0 to allow access to the outside world
ENV PORT="8181"
ENV REDIS="localhost:6380"
EXPOSE 8181

# Stripe requires a certificate
COPY --from=builder /etc/ssl/certs/ca-certificates.crt /etc/ssl/certs/
COPY --from=builder /src/main /main

ENTRYPOINT ["/main", "serve"]

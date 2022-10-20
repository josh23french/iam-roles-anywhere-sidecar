FROM golang:alpine AS builder

# Git is required for fetching the dependencies.
RUN apk update && apk add --no-cache git ca-certificates && update-ca-certificates

ENV USER=changeme
ENV UID=10001 
# See https://stackoverflow.com/a/55757473/12429735RUN 
RUN adduser \    
    --disabled-password \    
    --gecos "" \    
    --home "/nonexistent" \    
    --shell "/sbin/nologin" \    
    --no-create-home \    
    --uid "${UID}" \    
    "${USER}"

WORKDIR $GOPATH/src/github.com/josh23french/iam-roles-anywhere-sidecar

COPY go.mod .
COPY go.sum .

# Fetch dependencies.
RUN go mod download
RUN go mod verify

COPY . .

# Build the binary.
RUN CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -ldflags="-w -s" -o /iam-roles-anywhere-sidecar

FROM scratch

COPY --from=builder /etc/ssl/certs/ca-certificates.crt /etc/ssl/certs/
COPY --from=builder /etc/passwd /etc/passwd
COPY --from=builder /etc/group /etc/group
COPY --from=builder /iam-roles-anywhere-sidecar /usr/local/bin/iam-roles-anywhere-sidecar

USER changeme:changeme

ENTRYPOINT ["/usr/local/bin/iam-roles-anywhere-sidecar"]

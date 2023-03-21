FROM golang:alpine AS builder
ARG VERSION

# install certificates to copy it to the scratch container later
# needed to let the validation of TLS certs work
RUN apk update && apk add --no-cache ca-certificates git

WORKDIR /build
COPY . .
# 
# install go dependencies
RUN go mod tidy

# Build go
RUN go build \
    -ldflags "-s -w -X main.version=${VERSION}" \
    -o /go-myapps-sysclienttester \
    .

# make it executable, not sure if really needed when running in scratch container later on
RUN chmod 775 /go-myapps-sysclienttester

# try to start it to let the build fail in case of emergency
RUN /go-myapps-sysclienttester -h

FROM scratch 
ARG VERSION
WORKDIR /
COPY --from=builder /go-myapps-sysclienttester /go-myapps-sysclienttester
COPY --from=builder /build/identities.json /identities.json
COPY --from=builder /etc/ssl/certs/ca-certificates.crt /etc/ssl/certs/ca-certificates.crt
LABEL \
  org.opencontainers.image.vendor="Rico Schulte" \
  org.opencontainers.image.title="go-myapps-sysclienttester" \
  org.opencontainers.image.source="https://github.com/ricoschulte/go-myapps-sysclienttester" \
  org.opencontainers.image.version="${VERSION}"


VOLUME ["/data"]
ENTRYPOINT ["/go-myapps-sysclienttester"]
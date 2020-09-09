# Builder image
# incomplete

FROM golang:1.14-buster as builder

# Disable CGO so binary is fully static, then ensure that Go modules are enabled
ENV CGO_ENABLED=0 GO111MODULE=on
ENV PACKR_VERSION=2.7.1

# Fetch packr2 which we require to build the app
RUN curl -L "https://github.com/gobuffalo/packr/releases/download/v${PACKR_VERSION}/packr_${PACKR_VERSION}_linux_amd64.tar.gz" | tar -xvz \
    && mv packr2 /bin \
    && chmod +x /bin/packr2

# Copy the app
#

# Build the app
# RUN go get -d -v ./... \
#    && /bin/packr2 \


# Runtime image
#
FROM gcr.io/distroless/static

# Env vars for the runtime 
ENV BIND=0.0.0.0

# Copy the binary from the builder image
COPY --from=builder /go/bin/api-server /

# Entrypoint
CMD ["/api-server"]

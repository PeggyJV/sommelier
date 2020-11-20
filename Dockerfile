FROM golang:alpine AS build-env

# Set up dependencies
ENV PACKAGES curl make git libc-dev bash gcc linux-headers eudev-dev python3

# Set working directory for the build
WORKDIR /go/src/github.com/peggyjv/sommelier

# Add source files
COPY . .

# Install minimum necessary dependencies, build Cosmos SDK, remove packages
RUN apk add --no-cache $PACKAGES && \
    make install

# Final image
FROM alpine:edge

ENV SOMM /somm

# Install ca-certificates
RUN apk add --update ca-certificates

RUN addgroup sommuser && \
    adduser -S -G sommuser sommuser -h "$SOMM"
    
USER sommuser

WORKDIR $SOMM

# Copy over binaries from the build-env
COPY --from=build-env /go/bin/sommelier /usr/bin/sommelier

CMD ["sommelier"]

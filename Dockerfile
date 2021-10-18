FROM golang:alpine AS build-env

RUN apk add --no-cache curl make git libc-dev bash gcc linux-headers eudev-dev python3

# Set working directory for the build
WORKDIR /go/src/github.com/peggyjv/sommelier

# Get dependancies - will also be cached if we won't change mod/sum
COPY go.mod .
COPY go.sum .
RUN go mod download

# Add source files
COPY . .

# build Sommelier
RUN make install

# Final image
FROM alpine:edge

#ENV SOMM /somm

# Install ca-certificates
RUN apk add --update ca-certificates bash

# create and expose config dir
#RUN mkdir -p /somm/.sommelier/config
#RUN chmod -R 777 /somm/.sommelier/config
#RUN mkdir  /somm/.sommelier/data
#RUN chmod 777 /somm/.sommelier/data

# RUN addgroup sommuser && \
#    adduser -S -G sommuser sommuser -h "$SOMM"
    
#USER sommuser
#EXPOSE 1317 6060 6061 6062 6063 6064 6065 9090 26656 26657

#WORKDIR $SOMM

# Copy over binaries from the build-env
COPY --from=build-env /go/bin/sommelier /usr/bin/sommelier

CMD ["sommelier", "start"]

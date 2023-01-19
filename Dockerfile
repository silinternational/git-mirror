FROM golang:1.19

RUN curl -o- -L https://slss.io/install | VERSION=3.26.0 bash

# Copy in source and install deps
WORKDIR /src
COPY ./ /src/
RUN go get ./...

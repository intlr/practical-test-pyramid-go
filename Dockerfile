FROM golang:1.16

WORKDIR /go/src/github.com/alr-lab/ptp

COPY . /go/src/github.com/alr-lab/ptp

RUN apt-get update && apt-get install --yes make && make

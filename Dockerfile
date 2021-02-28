FROM golang:1.16

WORKDIR /go/src/github.com/alr-lab/practical-test-pyramid-go

COPY . /go/src/github.com/alr-lab/practical-test-pyramid-go

RUN apt-get update && apt-get install --yes make && make

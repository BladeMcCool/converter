FROM golang:1.15.8-buster
ENV GO111MODULE=on

WORKDIR /go/src
RUN git clone https://github.com/BladeMcCool/converter
WORKDIR /go/src/converter/src

RUN go get
RUN go build -o ../bin/converter
EXPOSE 5445
CMD ["converter"]

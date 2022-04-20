FROM golang:1.18.1

RUN apt-get update && apt-get install -y curl make git jq

ENV GOROOT=/usr/local/go
ENV GOPATH=/go
ENV PATH=$PATH:$GOROOT/bin:$GOPATH/bin

RUN go get -u github.com/jteeuwen/go-bindata/...

RUN mkdir -p $GOPATH/src/github.com/bhojpur/
WORKDIR $GOPATH/src/github.com/bhojpur/kernel

COPY ./ $GOPATH/src/github.com/bhojpur/kernel

CMD make -e TARGET_OS=${TARGET_OS} localbuild && mv ./kernctl /opt/build/kernctl
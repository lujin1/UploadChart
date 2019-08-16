FROM golang:1.12 as build-env

WORKDIR /go/src/uploadchart
ADD . /go/src/uploadchart

RUN go get -d -v github.com/ddliu/go-httpclient && go get -d -v github.com/urfave/cli && go get -d -v github.com/mholt/archiver
RUN go install

FROM golang
WORKDIR /
COPY --from=build-env /go/bin/uploadchart /
CMD ["/uploadchart"]

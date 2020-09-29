FROM golang:1.15.2-alpine3.12 AS build-env
WORKDIR /go/src/github.com/koron/httpreqinfo
COPY . .
RUN go install

FROM alpine:3.12.0
COPY --from=build-env /go/bin/httpreqinfo /usr/local/bin/
CMD ["/usr/local/bin/httpreqinfo"]

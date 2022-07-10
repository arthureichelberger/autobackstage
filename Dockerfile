FROM golang:1.18 as builder

ADD . /build
WORKDIR /build

RUN go get -v ./... && make build
RUN ls

FROM scratch

COPY --from=builder /build/deployment/bin/autobackstage /opt/bin/autobackstage
COPY --from=builder /etc/ssl/certs/ca-certificates.crt /etc/ssl/certs/

CMD ["/opt/bin/autobackstage"]

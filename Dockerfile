FROM golang:1.19.1-buster as builder

ADD . /certigo
WORKDIR /certigo
RUN go build

FROM golang:1.19.1-buster
COPY --from=builder /certigo/certigo /
COPY --from=builder /certigo/test-certs/*.crt /testsuite/

ENTRYPOINT []
CMD ["/certigo", "dump", "@@"]


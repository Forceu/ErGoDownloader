FROM golang:1.16.6

## To compile:
## cd Ergodownloader/build/
## docker build . --tag ergo-builder
## docker run --rm -it -v ../:/usr/src/myapp -w /usr/src/myapp ergo-builder

RUN \
  apt-get update && \
  apt-get install -y ca-certificates openssl zip && \
  update-ca-certificates && \
  rm -rf /var/lib/apt

COPY go.mod /tmp/tmp/go.mod

RUN cd /tmp/tmp/ && go mod download && rm -r /tmp/tmp

COPY entrypoint.sh /entrypoint.sh

ENTRYPOINT ["/entrypoint.sh"]

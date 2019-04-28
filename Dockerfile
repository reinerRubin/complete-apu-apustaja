FROM golang:1.12 as builder

WORKDIR /src

# cache dependencies
COPY go.mod go.su* /src/
RUN go mod download

# build the app
# TODO: make more catchable and copy only go code
ADD . /src
RUN make build

# run tests
FROM builder as test
WORKDIR /src
RUN go test -v -count 1 ./...

# do a release container
FROM alpine:latest as release
RUN apk update && apk add ca-certificates && rm -rf /var/cache/apk/*
WORKDIR /root
COPY --from=builder /src/bin/completer .

CMD ["/root/completer"]

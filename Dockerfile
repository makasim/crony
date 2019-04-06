FROM golang AS build-env


RUN apt-get update && \
    apt-get -y --no-install-recommends --no-install-suggests install git && \
    rm -rf /var/lib/apt/lists/*

ADD . /go/src/github.com/formapro/crony
WORKDIR /go/src/github.com/formapro/crony

RUN go get -d -v
RUN CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -o crony -v

FROM alpine
WORKDIR /app

COPY --from=build-env /go/src/github.com/formapro/crony/crony /app
RUN chmod u+x /app/crony

ENTRYPOINT ["./crony"]

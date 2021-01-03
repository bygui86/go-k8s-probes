
FROM golang:1.15-buster AS gobuilder

WORKDIR /go/src/github.com/bygui86/go-k8s-probes
COPY . .

ENV GO111MODULE=on
ENV CGO_ENABLED=0
ENV GOOS=linux

RUN go build -a -installsuffix cgo -o /bin/app .

# ---

FROM alpine

RUN apk update --no-cache
RUN apk add --no-cache bash
RUN apk add --no-cache curl

WORKDIR /usr/bin/
COPY --from=gobuilder /bin/app .

EXPOSE 8080
EXPOSE 9090
EXPOSE 9091

USER 1001

ENTRYPOINT "/usr/bin/app"

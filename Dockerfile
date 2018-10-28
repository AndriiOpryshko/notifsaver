FROM golang:1.11.1-alpine3.8 as build-env
RUN mkdir /notifsaver
WORKDIR /notifsaver
COPY go.mod .
COPY go.sum .


RUN apk add --update --no-cache ca-certificates git
RUN go mod download
COPY . .

RUN CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -a -installsuffix cgo -o /go/bin/notifsaver

FROM scratch
COPY --from=build-env /go/bin/notifsaver /go/bin/notifsaver
COPY config.yml .
ENTRYPOINT ["/go/bin/notifsaver"]
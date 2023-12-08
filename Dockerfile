# Stage 1
FROM golang:1.17-alpine3.14 as builder

# Add git
RUN apk update && \
    apk add git && \
    apk add gcc && \
    apk add libc-dev

RUN mkdir /ncronus

ADD . /ncronus

WORKDIR /ncronus/cmd

RUN env GOOS=linux GOARCH=amd64 go build .

# Stage 2
FROM alpine:3.14

RUN apk update && apk add ca-certificates
RUN update-ca-certificates

COPY --from=builder /ncronus/cmd /

EXPOSE 3001

CMD ["./cmd"]

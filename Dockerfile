FROM golang:latest AS build

COPY migrations /migrations

COPY bravo-svc /bravo-svc

EXPOSE 8888

WORKDIR /

CMD ["./bravo-svc"]

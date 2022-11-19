FROM golang:latest AS build

COPY bravo /bravo

EXPOSE 8888

WORKDIR /

CMD ["./bravo"]

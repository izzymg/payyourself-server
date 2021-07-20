FROM golang:1.16
WORKDIR /src
COPY ./server .
RUN go get -d -v ./...
RUN CGO_ENABLED=0 GOOS=linux go build -o /bin/py-server .

FROM alpine:latest
COPY --from=0 /bin/py-server /bin/py-server
EXPOSE 6002
CMD ["/bin/py-server"]
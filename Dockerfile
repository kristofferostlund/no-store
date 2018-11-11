FROM golang:latest AS build

# RUN mkdir /usr/app
ADD . /usr/app

WORKDIR /usr/app

RUN go get -v ./...

RUN CGO_ENABLED=0 GOOS=linux go build -a -installsuffix cgo -v -o nostore-server server/main.go

FROM scratch

WORKDIR /usr/app

COPY --from=build /usr/app/nostore-server /nostore-server

EXPOSE 80

ENTRYPOINT [ "/nostore-server" ]
CMD ["-address=0.0.0.0", "-port=80"]

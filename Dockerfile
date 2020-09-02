FROM golang:alpine

WORKDIR /go/review_app

COPY ./ReviewApp/go.mod .
COPY ./ReviewApp/go.sum .

RUN go mod download

COPY ./ReviewApp/ .

RUN go install

RUN GOOS=linux GOARCH=amd64 go build

EXPOSE 5550

## Add the wait script to the image
ADD https://github.com/ufoscout/docker-compose-wait/releases/download/2.7.3/wait /wait
RUN chmod +x /wait

# CMD [ "/wait && ./ReviewApp" ]


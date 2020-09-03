FROM golang:alpine

WORKDIR /go/sertis_app

COPY ./sertis_backend/go.mod .
COPY ./sertis_backend/go.sum .

RUN go mod download

COPY ./sertis_backend/ .

RUN go install

RUN GOOS=linux GOARCH=amd64 go build

EXPOSE 8880

## Add the wait script to the image
ADD https://github.com/ufoscout/docker-compose-wait/releases/download/2.7.3/wait /wait
RUN chmod +x /wait

# CMD [ "/wait && ./sertis_app" ]

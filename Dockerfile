FROM golang:1.21

WORKDIR /usr/local/app

COPY . .

RUN go build
EXPOSE 8080
CMD [ "/usr/local/app/api" ]
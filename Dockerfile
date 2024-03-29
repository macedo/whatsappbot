FROM golang:alpine

WORKDIR /app

COPY . /app/

RUN go build -o whatsappbot

ENTRYPOINT [ "./whatsappbot" ]

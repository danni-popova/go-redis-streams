WORKDIR /app

COPY go.mod go.sum ./
COPY *.go ./

RUN go mod download

RUN go build -o /go/bin/app cmd/consumer/main.go

#EXPOSE 4080

CMD [ "/cmd/main" ]
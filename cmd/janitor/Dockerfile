WORKDIR /app

COPY go.mod go.sum ./
COPY *.go ./

RUN go mod download

RUN go build -o /go/bin/app cmd/janitor/main.go

EXPOSE 4085

CMD [ "/cmd/main" ]
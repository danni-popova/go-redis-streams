FROM golang:1.21 as builder

WORKDIR /app

COPY go.mod go.sum ./
COPY cmd .

RUN go mod download

RUN go build -o /go/bin/app cmd/producer/main.go

# We use a separate locked down runner container, so copy the binaries in
FROM gcr.io/distroless/static-debian12:nonroot as runner
WORKDIR /app
COPY --from=builder /go/bin/app /app/main

CMD [ "/app/main" ]
FROM golang:1.21 as builder

WORKDIR /app

COPY go.mod go.sum ./
COPY stream/ ./stream
COPY cmd/ ./cmd

RUN go mod download && go mod verify
ENV CGO_ENABLED=0
ENV GOOS=linux
ENV GOARCH=arm64
RUN go build -o /go/bin/app cmd/producer/main.go

# We use a separate locked down runner container, so copy the binaries in
FROM gcr.io/distroless/static-debian12:nonroot as runner

WORKDIR /app

COPY --from=builder /go/bin/app /app/main

CMD [ "/app/main" ]
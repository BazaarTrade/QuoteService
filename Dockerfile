FROM golang:1.23.1 AS builder

WORKDIR /app
COPY . .

RUN go mod tidy

WORKDIR /app/cmd
RUN go build -o quote .

FROM gcr.io/distroless/base

COPY --from=builder /app/cmd/quote /quote

EXPOSE 50052
CMD ["/quote"]
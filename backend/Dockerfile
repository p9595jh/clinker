FROM golang:latest AS builder
WORKDIR /backend
COPY . .
RUN go mod download
RUN go build -o backend ./cmd/backend

FROM golang:latest
WORKDIR /backend
COPY --from=builder /backend /backend
CMD ["/backend"]

EXPOSE 3001

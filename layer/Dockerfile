FROM golang:latest AS builder
WORKDIR /layer
COPY . .
RUN go mod download
RUN go build -o layer ./cmd/layer

FROM golang:latest
WORKDIR /layer
COPY --from=builder /layer /layer
CMD ["/layer"]

EXPOSE 3002

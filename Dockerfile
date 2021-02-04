FROM golang:1.15-alpine3.13 as builder

WORKDIR /app

COPY go.mod go.sum ./
RUN go mod download

COPY . .

RUN CGO_ENABLED=0 GOOS=linux go build ./cmd/foundry

FROM archlinux:base-devel-20210124.0.14185

WORKDIR /app

COPY --from=builder /app/foundry .

EXPOSE 10201
EXPOSE 10200

ENTRYPOINT ["./foundry"]


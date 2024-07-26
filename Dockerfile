FROM golang:1.22.5 AS builder
WORKDIR /app

COPY src ./

RUN go mod download
RUN CGO_ENABLED=0 GOOS=linux go build -o /auth

FROM golang:1.22.5-alpine3.20
COPY --from=builder /auth /auth

EXPOSE 5000

ENTRYPOINT [ "/auth" ]
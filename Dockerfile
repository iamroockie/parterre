FROM golang:1.27-trixie AS builder

WORKDIR /app

COPY go.mod go.sum ./
RUN go mod download

COPY . ./

RUN CGO_ENABLED=0 GOOS=linux \
    go build -trimpath -ldflags="-s -w" -o /parterre-api ./cmd/api

###

FROM alpine:3.24 AS runtime

RUN apk add --no-cache tzdata
RUN adduser -D -u 10001 app

COPY --from=builder /parterre-api /parterre-api

USER app
EXPOSE 8080

CMD [ "/parterre-api" ]

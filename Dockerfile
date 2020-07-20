FROM golang:1.14-alpine AS dev

RUN apk add --no-cache git build-base

WORKDIR /code
RUN go get -u github.com/cosmtrek/air

COPY go.mod go.sum ./
RUN go mod download
COPY . .
RUN go build -o ./tmp/goat .
CMD ["air"]


FROM alpine:latest AS prod

ARG ENV=local
ARG VERSION

COPY --from=dev /code/tmp/goat /goat
RUN apk add --no-cache curl

ENV ENV=$ENV
ENV VERSION=$VERSION

RUN chmod +x /goat
CMD ["/goat"]

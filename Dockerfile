FROM golang:1.14-alpine AS dev

RUN apk add --no-cache git build-base

WORKDIR /code
RUN go get -u github.com/cosmtrek/air

COPY go.mod go.sum ./
RUN go mod download
COPY . .
RUN go build -o ./tmp/goat . && go build -o ./tmp ./cmd/...
CMD ["air"]


FROM alpine:latest AS prod

ARG ENV=local
ARG VERSION

COPY --from=dev /code/tmp/ /
RUN apk add --no-cache curl

ENV ENV=$ENV
ENV VERSION=$VERSION

RUN chmod +x /goat
CMD ["/goat"]

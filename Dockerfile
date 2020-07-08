FROM golang:1.14-alpine AS dev

RUN apk add --no-cache git build-base
WORKDIR /code
COPY go.mod go.sum ./
RUN go mod download
RUN go get -u github.com/cosmtrek/air
COPY . .
RUN go build -o ./tmp/goat .
CMD ["air"]


FROM alpine:latest AS prod
COPY --from=dev /code/tmp/goat /goat
RUN chmod +x /goat
CMD ["/goat"]

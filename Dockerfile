FROM golang:1.17.0-alpine3.14 as build
ENV CGO_ENABLED 0
COPY . /service

WORKDIR /service/app/web
RUN go build

FROM alpine:3.14
COPY --from=build /service/app/web/web /service/web
WORKDIR /service
CMD ["./web"]

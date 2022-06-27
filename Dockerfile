FROM golang:alpine3.16 as build
WORKDIR /app

COPY . .
RUN go mod download

RUN go build -o /app/main .

FROM alpine:3.16 as final

COPY --from=build /app/main /
CMD ["/main"]
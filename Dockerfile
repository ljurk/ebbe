FROM golang:alpine as build
WORKDIR /
COPY go.mod go.sum ./
RUN go mod download
COPY . .
RUN go build -o app .

FROM alpine:latest
WORKDIR /
COPY --from=build /app /app
ENV EBBE_REDIS_URL="" \
    EBBE_PIXELFLUT_URL=""
CMD ["sh", "-c", "./app redis worker --url $EBBE_REDIS_URL --host $EBBE_PIXELFLUT_URL"]

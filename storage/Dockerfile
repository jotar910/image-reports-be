# Dockerfile.production
FROM golang:1.19-alpine as builder

ENV APP_HOME /go/src/irbe-storage

WORKDIR "$APP_HOME"
COPY . .

RUN go mod download
RUN go mod verify
RUN go build -o irbe-storage ./storage/.

FROM golang:1.19-alpine

ENV APP_HOME /go/src/irbe-storage
RUN mkdir -p "$APP_HOME"
WORKDIR "$APP_HOME"

COPY --from=builder "$APP_HOME" $APP_HOME

EXPOSE 8080
ENTRYPOINT ["./irbe-storage"]
CMD [ "-m", "prod" ]
# Dockerfile.production
FROM golang:1.19-alpine as builder

ENV APP_HOME /go/src/irbe-realtime

WORKDIR "$APP_HOME"
COPY . .

RUN go mod download
RUN go mod verify
RUN go build -o irbe-realtime ./realtime/.

FROM golang:1.19-alpine

ENV APP_HOME /go/src/irbe-realtime
RUN mkdir -p "$APP_HOME"
WORKDIR "$APP_HOME"

COPY --from=builder "$APP_HOME" $APP_HOME

EXPOSE 8080
ENTRYPOINT ["./irbe-realtime"]
CMD [ "-m", "prod" ]
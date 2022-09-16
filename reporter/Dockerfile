# Dockerfile.production
FROM golang:1.19-alpine as builder

ENV APP_HOME /go/src/irbe-reporter

WORKDIR "$APP_HOME"
COPY . .

RUN go mod download
RUN go mod verify
RUN go build -o irbe-reporter ./reporter/.

FROM golang:1.19-alpine

ENV APP_HOME /go/src/irbe-reporter
RUN mkdir -p "$APP_HOME"
WORKDIR "$APP_HOME"

COPY --from=builder "$APP_HOME" $APP_HOME

EXPOSE 8080
ENTRYPOINT ["./irbe-reporter"]
CMD [ "-m", "prod" ]
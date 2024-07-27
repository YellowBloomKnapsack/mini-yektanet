FROM golang:1.22.4-bullseye as builder
# FROM registry.docker.ir/golang:1.22.4-alpine as builder

ARG FOLDER_NAME
# ENV FOLDER_NAME=$FOLDER_NAME

ENV APP_HOME /app/

WORKDIR "$APP_HOME"

COPY go.mod ./
RUN go mod download


COPY ${FOLDER_NAME} .
COPY common .
COPY adserver/.env adserver
COPY eventserver/.env eventserver
COPY panel/.env panel
COPY publisherwebsite/.env publisherwebsite

RUN go mod verify
WORKDIR ${FOLDER_NAME}
RUN go build -o ${FOLDER_NAME}.out

# FROM registry.docker.ir/golang:1.22.4-alpine
FROM golang:1.22.4-bullseye

ENV APP_HOME /app/${FOLDER_NAME}
RUN mkdir -p "$APP_HOME"
WORKDIR "$APP_HOME"

# COPY src/conf/ conf/
# COPY src/views/ views/
COPY --from=builder "$APP_HOME"/${FOLDER_NAME}.out $APP_HOME
COPY --from=builder "$APP_HOME"/adserver $APP_HOME
COPY --from=builder "$APP_HOME"/eventserver $APP_HOME
COPY --from=builder "$APP_HOME"/panel $APP_HOME
COPY --from=builder "$APP_HOME"/publisherwebsite $APP_HOME

EXPOSE 8083
CMD ["./${FOLDER_NAME}.out"]

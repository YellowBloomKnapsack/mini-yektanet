FROM golang:1.22.4-bullseye as builder
# FROM registry.docker.ir/golang:1.22.4-alpine as builder

ARG FOLDER_NAME
ARG PORT_NUMBER
# ENV FOLDER_NAME=$FOLDER_NAME
ENV BUILDER_HOME /app/

WORKDIR "$BUILDER_HOME"

COPY go.mod ./
RUN go mod download


COPY ${FOLDER_NAME} ./${FOLDER_NAME}
COPY common ./common
COPY adserver/.env adserver/.env
COPY eventserver/.env eventserver/.env
COPY panel/.env panel/.env
COPY publisherwebsite/.env publisherwebsite/.env

COPY go.sum ./
RUN go mod verify
WORKDIR ${FOLDER_NAME}
RUN go build -o output.out

# FROM registry.docker.ir/golang:1.22.4-alpine
FROM golang:1.22.4-bullseye
ARG FOLDER_NAME
ENV envFOLDER_NAME=$FOLDER_NAME
RUN mkdir -p /app
ENV APP_HOME /app/${FOLDER_NAME}
ENV BUILDER_HOME /app/
RUN mkdir -p "$APP_HOME"
WORKDIR "$APP_HOME"

# COPY src/conf/ conf/
# COPY src/views/ views/
COPY --from=builder "$BUILDER_HOME"/"${FOLDER_NAME}"/output.out $APP_HOME
COPY --from=builder "$BUILDER_HOME"/"${FOLDER_NAME}"/.env $APP_HOME
COPY --from=builder "$BUILDER_HOME"/adserver /app/adserver
COPY --from=builder "$BUILDER_HOME"/eventserver /app/eventserver
COPY --from=builder "$BUILDER_HOME"/panel /app/panel
COPY --from=builder "$BUILDER_HOME"/publisherwebsite /app/publisherwebsite
COPY --from=builder "$BUILDER_HOME"/common /app/common

EXPOSE $PORT_NUMBER

# CMD ["/bin/bash"]
CMD ["./output.out"]
# CMD ["./${FOLDER_NAME}.out"]
# CMD ./${envFOLDER_NAME}.out
# CMD ["sh", "-c", "${envFOLDER_NAME}.out"]


FROM golang:1.22.4-bullseye as builder

ARG FOLDER_NAME
ARG PORT_NUMBER
ENV BUILDER_HOME /app/

WORKDIR "$BUILDER_HOME"

COPY go.mod ./
COPY go.sum ./
RUN go mod download
RUN go mod verify

COPY ${FOLDER_NAME} ./${FOLDER_NAME}
COPY common ./common
COPY adserver/.env adserver/.env
COPY eventserver/.env eventserver/.env
COPY panel/.env panel/.env
COPY publisherwebsite/.env publisherwebsite/.env

WORKDIR ${FOLDER_NAME}
RUN go build -o output.out

FROM golang:1.22.4-bullseye
ARG FOLDER_NAME
ENV envFOLDER_NAME=$FOLDER_NAME
RUN mkdir -p /app
ENV APP_HOME /app/${FOLDER_NAME}
ENV BUILDER_HOME /app/
RUN mkdir -p "$APP_HOME"
WORKDIR "$APP_HOME"

COPY --from=builder "$BUILDER_HOME"/"${FOLDER_NAME}"/output.out $APP_HOME
COPY --from=builder "$BUILDER_HOME"/"${FOLDER_NAME}"/.env $APP_HOME
COPY --from=builder "$BUILDER_HOME"/adserver /app/adserver
COPY --from=builder "$BUILDER_HOME"/eventserver /app/eventserver
COPY --from=builder "$BUILDER_HOME"/panel /app/panel
COPY --from=builder "$BUILDER_HOME"/publisherwebsite /app/publisherwebsite
COPY --from=builder "$BUILDER_HOME"/common /app/common

EXPOSE $PORT_NUMBER

CMD ["./output.out"]


ARG GOLANG_IMAGE
ARG ALPINE_IMAGE

## GOLANG BUILD STAGE##
FROM $GOLANG_IMAGE as build-env

ARG APP_ENV
ARG APP_NAME
ARG APP_VERSION

ENV APP_ENV=$APP_ENV
ENV APP_NAME=$APP_NAME
ENV APP_VERSION=$APP_VERSION

WORKDIR /src

COPY ["go.mod", "go.sum", "./"]

RUN go mod download
# add flag `-x` for more verbose

COPY Makefile .
COPY configs  ./configs
COPY cmd  ./cmd
COPY internal  ./internal
COPY protocol  ./protocol
COPY api  ./api
COPY tools  ./tools
COPY migrations  ./migrations
COPY docs  ./docs

# go build project and store artifact to Nexus artifactory
RUN : \
    && mkdir -p /dist /nexus /nexus/configs                              \
#    && make swagger_docs                                                 \
    && make build                                                        \
    && cp ./bin/${APP_NAME} /nexus                                       \
    && cp ./configs/default.yml ./configs/${APP_ENV}.yml /nexus/configs  \
    && mv -v ./bin/${APP_NAME} /dist                                     \
    && tar -czvf ${APP_NAME}_${APP_VERSION}.tar.gz /nexus                \
    && mv -v ${APP_NAME}_${APP_VERSION}.tar.gz /dist

## APP STAGE ##
FROM $ALPINE_IMAGE as app

ARG APP_NAME
ENV APP_NAME=$APP_NAME

RUN : mkdir -p /opt/${APP_NAME}/                  \
               /opt/${APP_NAME}/configs

COPY --from=build-env /dist/${APP_NAME}    /opt/${APP_NAME}/${APP_NAME}
COPY --from=build-env /src/configs         /opt/${APP_NAME}/configs

# fix x509 error  https://stackoverflow.com/questions/52969195/docker-container-running-golang-http-client-getting-error-certificate-signed-by
COPY --from=build-env /etc/ssl/certs/ca-certificates.crt /etc/ssl/certs/

WORKDIR /opt/${APP_NAME}
ENTRYPOINT ["./${APP_NAME}"]


## STORE TO NEXUS ARTIFACTORY
FROM build-env as nexus

# Artifact's storagePostgres
ARG NEXUS_HOST
ARG NEXUS_PORT
ARG NEXUS_USER
ARG NEXUS_PSWD

ENV NEXUS_HOST=$NEXUS_HOST
ENV NEXUS_PORT=$NEXUS_PORT
ENV NEXUS_USER=$NEXUS_USER
ENV NEXUS_PSWD=$NEXUS_PSWD

WORKDIR /dist
RUN export \
    && ls -la \
    && curl -v -u ${NEXUS_USER}:${NEXUS_PSWD} --upload-file ${APP_NAME}_${APP_VERSION}.tar.gz http://${NEXUS_HOST}:${NEXUS_PORT}/repository/theseus-${APP_NAME}/${APP_NAME}_${APP_VERSION}.tar.gz


## DLV DEBUG STAGE ##
FROM app as app-debug

# Allow delve to run on Alpine based containers.
RUN apk add --no-cache libc6-compat

# Copy dlv binary file from build-env stage
COPY --from=build-env /go/bin/dlv /opt/${APP_NAME}/dlv

# IF Need dlv to debug use this
ENTRYPOINT ["sleep", "10000000"]

FROM docker.io/golang:1.26-alpine AS builder

WORKDIR /app
ENV CGO_ENABLED=0

COPY . .

RUN for p in cmd/*; do go build -buildvcs=false -o "${p/cmd\//}" "./$p"; done

FROM gcr.io/distroless/static-debian11

ARG BUILD_DATE
ARG SOURCE_COMMIT

LABEL org.opencontainers.image.authors="blankdots"
LABEL org.opencontainers.image.created=$BUILD_DATE
LABEL org.opencontainers.image.source="https://github.com/blankdots/minimal-kube-app"
LABEL org.opencontainers.image.licenses="Apache-2.0"
LABEL org.opencontainers.image.title="minimal-kube-app"

COPY --from=builder /app/api /usr/bin/api
COPY --from=builder /app/cronjob /usr/bin/cronjob

USER 65534
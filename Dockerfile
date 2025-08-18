FROM golang:latest AS builder
WORKDIR /app
COPY . .
RUN CGO_ENABLED=0 go build -v -trimpath -ldflags '-d -w -s'

FROM scratch
ARG REVISION
LABEL org.opencontainers.image.title="Fast vanity age X25519 identity generator"
LABEL org.opencontainers.image.description="This tool generates age X25519 identity with a recipient that has a specified prefix"
LABEL org.opencontainers.image.authors="Alexander Yastrebov <yastrebov.alex@gmail.com>"
LABEL org.opencontainers.image.url="https://github.com/AlexanderYastrebov/age-vanity-keygen"
LABEL org.opencontainers.image.licenses="BSD-3-Clause"
LABEL org.opencontainers.image.revision="${REVISION}"

COPY --from=builder /app/age-vanity-keygen /age-vanity-keygen

ENTRYPOINT ["/age-vanity-keygen"]

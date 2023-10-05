FROM alpine:3.18

ARG BINARY_NAME=server

EXPOSE 8080

RUN apk update && \
    apk add --no-cache ca-certificates && \
    update-ca-certificates

COPY --chmod=0755 $BINARY_NAME /server

CMD ["/server"]

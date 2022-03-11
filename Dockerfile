FROM alpine:3.13
RUN apk add --no-cache ca-certificates
COPY check-christmas-lottery-numbers /usr/bin/check-christmas-lottery-numbers
ENTRYPOINT ["/usr/bin/check-christmas-lottery-numbers"]
FROM golang:1.14-alpine3.13 as builder

# Add Maintainer Info
LABEL maintainer="Xoan Mallon <xoanmallon@gmail.com>"

RUN apk --no-cache add git
WORKDIR /go/src/build
COPY . .
RUN export CGO_ENABLED=0 \
    && mkdir -p dist \
    && go mod vendor \
    && go build -o dist/check-christmas-lottery-numbers ./cmd

FROM alpine:3.13
COPY --from=builder /go/src/build/dist/ .

ENTRYPOINT ["./check-christmas-lottery-numbers"]
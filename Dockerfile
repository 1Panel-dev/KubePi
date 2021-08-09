FROM node:14-alpine as stage-web-build

LABEL stage=stage-web-build

RUN apk add make

WORKDIR /build/ekko/web

COPY . .

RUN make build_web

RUN rm -fr web


FROM golang:1.16 as stage-bin-build

ENV GOPROXY="https://goproxy.cn,direct"

LABEL stage=stage-bin-build

WORKDIR /build/ekko/bin

COPY --from=stage-web-build /build/ekko/web .

ENV GO111MODULE=on

RUN go mod download

RUN make build_bin

FROM alpine:3.14

WORKDIR /

COPY --from=stage-bin-build /build/ekko/bin/dist/usr /usr

EXPOSE 2019

USER root

CMD ["ekko-server"]
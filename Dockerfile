FROM golang:1.13.4 as build

ENV GOPROXY https://goproxy.cn
ENV GO111MODULE on
ENV GIN_MODE relaesa

WORKDIR /go/cache

ADD go.mod .
RUN go mod download

WORKDIR /go/release

ADD . .

RUN GOOS=linux CGO_ENABLED=0 go build -ldflags="-s -w" -installsuffix cgo -o app

FROM scratch as prod

COPY --from=build /usr/share/zoneinfo/Asia/Shanghai /etc/localtime
COPY --from=build /etc/ssl/certs/ca-certificates.crt /etc/ssl/certs/ca-certificates.crt
COPY --from=build /go/release/app /
COPY --from=build /go/release/conf /conf

CMD ["/app"]


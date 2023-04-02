FROM golang:1.18.4-bullseye
# author xuri.li

RUN mkdir /usr/local/goc
RUN mkdir /usr/local/app

ADD . /usr/local/goc

RUN go env -w GO111MODULE=on
RUN go env -w GOPROXY=https://goproxy.cn

WORKDIR /usr/local/goc

RUN go build ./runner/main.go

RUN mv main /usr/local/app/authz
# RUN mv config /usr/local/app
# users should pass the config file

WORKDIR /usr/local/app

RUN rm -rf /usr/local/goc

ENTRYPOINT ["./authz"]
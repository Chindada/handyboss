FROM golang:1.13.15

USER root
ENV symbol=\!
ENV GO111MODULE=on
ENV TZ=Asia/Taipei
RUN ln -snf /usr/share/zoneinfo/$TZ /etc/localtime && echo $TZ > /etc/timezone
WORKDIR /go/src/app
RUN git clone http://timhsu:${symbol}Minnotec2025@172.20.10.50/timhsu/handyboss.git
WORKDIR handyboss
RUN go build -o handyboss
ENTRYPOINT ./handyboss
EXPOSE 27001

FROM paulhindemith/golang-dev:1.13

WORKDIR /root

RUN wget -q https://github.com/prometheus/prometheus/releases/download/v2.15.2/prometheus-2.15.2.linux-amd64.tar.gz

RUN tar -zxf prometheus-2.15.2.linux-amd64.tar.gz && mv prometheus-2.15.2.linux-amd64/prometheus /usr/local/bin
FROM paulhindemith/golang-dev:1.13

WORKDIR /root

RUN apt-get update

RUN apt-get install -y apt-transport-https \
&& apt-get install -y software-properties-common wget \
&& wget -q -O - https://packages.grafana.com/gpg.key | apt-key add -  \
&& add-apt-repository "deb https://packages.grafana.com/oss/deb stable main" \
&& apt-get update \
&& apt-get install  -y grafana=6.6.0

WORKDIR /go/src/github.com/paulhindemith

CMD ["/bin/bash"]
FROM golang:1.4.1
ENV DEBIAN_FRONTEND noninteractive
RUN apt-get update
RUN apt-get install tor supervisor -y
ADD deployment/.env /.env
ADD deployment/start-tor.sh /start-tor.sh
ADD deployment/start-goffee.sh /start-goffee.sh
ADD deployment/run.sh /run.sh
RUN chmod 755 /*.sh
ADD deployment/supervisord-tor.conf /etc/supervisor/conf.d/supervisord-tor.conf
ADD deployment/supervisord-goffee.conf /etc/supervisor/conf.d/supervisord-goffee.conf
ADD . /go/src/github.com/goffee/goffee
WORKDIR /go/src/github.com/goffee/goffee
RUN go install github.com/goffee/goffee
CMD ["/run.sh"]
EXPOSE 80

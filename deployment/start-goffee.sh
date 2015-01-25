#!/bin/bash
source /.env
exec /go/bin/goffee -clientid="$CLIENTID" -secret="$SECRET" -bind :80 -mandrill="$MANDRILLKEY" -redisaddress="$REDIS_MASTER_SERVICE_HOST:$REDIS_MASTER_SERVICE_PORT" -mysql="$MYSQL" -sessionsecret="$SESSIONSECRET"

FROM java:8-alpine
MAINTAINER vishal

ENV KAFKA_BROKER_HOST=localhost \
    KAFKA_BROKER_PORT=9092 \
    appDir=/usr/src/app

RUN mkdir -p $appDir
WORKDIR $appDir

ADD build/libs/svc-payment-io.jar .
ADD data ./data

CMD ["java", \
            "-Djava.security.egd=file:/dev/./urandom", \
            "-Dspring.kafka.bootstrap-servers=${KAFKA_BROKER_HOST}:${KAFKA_BROKER_PORT}", \
            "-jar", \
            "svc-payment-io.jar"]

EXPOSE 8080
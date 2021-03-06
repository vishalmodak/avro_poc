FROM java:8-alpine
MAINTAINER vishal

ENV KAFKA_BROKER_URL=localhost:9092 \
    SCHEMA_REGISTRY_URL=http://localhost:8081 \
    LOAN_TOPIC_CONSUME=loanIn \
    LOAN_TOPIC_PRODUCE=loanOut \
    PAYMENT_TOPIC_CONSUME=paymentIn \
    STATSD_SERVER_HOST=localhost \
    PAYMENT_SERVICE_HOST=localhost \
    appDir=/usr/src/app

RUN mkdir -p $appDir
WORKDIR $appDir

ADD build/libs/svc-loan-io.jar .
ADD data ./data

CMD ["java", \
            "-Djava.security.egd=file:/dev/./urandom", \
            "-Dspring.kafka.bootstrap-servers=${KAFKA_BROKER_URL}", \
            "-Dspring.kafka.producer.properties.schema.registry.url=${SCHEMA_REGISTRY_URL}", \
            "-Dspring.kafka.consumer.properties.schema.registry.url=${SCHEMA_REGISTRY_URL}", \
            "-Dloan.topic.consume=${LOAN_TOPIC_CONSUME}", \
            "-Dloan.topic.produce=${LOAN_TOPIC_PRODUCE}", \
            "-Dpayment.topic.produce=${PAYMENT_TOPIC_CONSUME}", \
            "-Dservice.payments.uri=${PAYMENT_SERVICE_HOST}", \
            "-Dspring.cloud.stream.schemaRegistryClient.endpoint=${SCHEMA_REGISTRY_URL}", \
            "-Dmanagement.metrics.export.statsd.host=${STATSD_SERVER_HOST}", \
            "-jar", \
            "svc-loan-io.jar"]

EXPOSE 8080

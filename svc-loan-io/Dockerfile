FROM java:8-alpine
MAINTAINER vishal

ENV KAFKA_BROKER_URL=localhost:9092 \
    SCHEMA_REGISTRY_URL=http://localhost:8081 \
    appDir=/usr/src/app

RUN mkdir -p $appDir
WORKDIR $appDir

ADD build/libs/svc-loan-io.jar .
ADD data ./data

CMD ["java", \
            "-Djava.security.egd=file:/dev/./urandom", \
            "-Dspring.kafka.bootstrap-servers=${KAFKA_BROKER_URL}", \
            "-Dspring.kafka.producer.properties.schema.registry.url=${SCHEMA_REGISTRY_URL}", \
            "-Dspring.cloud.stream.schemaRegistryClient.endpoint=${SCHEMA_REGISTRY_URL}", \
            "-jar", \
            "svc-loan-io.jar"]

EXPOSE 8080
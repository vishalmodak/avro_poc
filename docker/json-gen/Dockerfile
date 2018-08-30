FROM airdock/oracle-jdk:1.8

WORKDIR /root

RUN curl -k -o json-data-generator-bin.tar -L https://github.com/acesinc/json-data-generator/releases/download/json-data-generator-1.3.0/json-data-generator-1.3.0-bin.tar; tar -xf json-data-generator-bin.tar -C /; mv /json-data-generator-1.3.0 /json-data-generator

WORKDIR /json-data-generator

RUN cp conf/exampleSimConfig.json conf/SimConfig.json; mkdir logs

VOLUME ["/json-data-generator/logs", "/json-data-generator/conf"]

ENTRYPOINT ["java", "-jar", "json-data-generator-1.3.0.jar", "SimConfig.json" ]

package com.loan.io.kafka;

import java.io.IOException;
import java.nio.file.Files;
import java.nio.file.Paths;

import org.slf4j.Logger;
import org.slf4j.LoggerFactory;
import org.springframework.beans.factory.annotation.Autowired;
import org.springframework.beans.factory.annotation.Value;
import org.springframework.kafka.core.KafkaTemplate;
import org.springframework.stereotype.Service;

import com.lss.models.Loan;

import tech.allegro.schema.json2avro.converter.JsonAvroConverter;

@Service
public class LoanPublisher {
    
    private static final Logger log = LoggerFactory.getLogger(LoanPublisher.class);

    @Autowired
    private KafkaTemplate<String, Loan> kafkaTemplate;
    
    @Value("${kafka.topic.produce}")
    private String producerTopic;

    public void send(String payload) {
        log.info("sending payload='{}' to topic='{}'", payload, producerTopic);
        Loan loan = convertJsonToAvro(payload);
        kafkaTemplate.send(producerTopic, loan);
    }

    public Loan convertJsonToAvro(String json) {
        JsonAvroConverter converter = new JsonAvroConverter();
        try {
            byte[] schemaBytes = Files.readAllBytes(Paths.get("src/main/avro/loan.avsc"));
            String schema = new String(schemaBytes);
            return converter.convertToSpecificRecord(json.getBytes(), Loan.class, schema);
        } catch (IOException e) {
            e.printStackTrace();
        }
        return null;
    }
}

package com.loan.io.kafka;

import java.io.File;
import java.io.IOException;

import org.apache.commons.io.FileUtils;
import org.slf4j.Logger;
import org.slf4j.LoggerFactory;
import org.springframework.beans.factory.annotation.Autowired;
import org.springframework.beans.factory.annotation.Value;
import org.springframework.core.io.ClassPathResource;
import org.springframework.core.io.Resource;
import org.springframework.kafka.core.KafkaTemplate;
import org.springframework.stereotype.Service;
import org.springframework.util.FileCopyUtils;
import org.springframework.util.ResourceUtils;

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
        Loan loan = convertJsonToAvro(payload);
        log.info("sending loan='{}' to topic='{}'", loan, producerTopic);
        kafkaTemplate.send(producerTopic, loan);
    }

    public Loan convertJsonToAvro(String json) {
        JsonAvroConverter converter = new JsonAvroConverter();
        try {
            ClassPathResource loanSchema = new ClassPathResource("avro/loan.avsc");
            byte[] bdata = FileCopyUtils.copyToByteArray(loanSchema.getInputStream());
            String schema = new String(bdata);
            return converter.convertToSpecificRecord(json.getBytes(), Loan.class, schema);
        } catch (IOException e) {
            e.printStackTrace();
        }
        return null;
    }
}

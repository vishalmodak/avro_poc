package com.payment.io.kafka;

import org.slf4j.Logger;
import org.slf4j.LoggerFactory;
import org.springframework.beans.factory.annotation.Autowired;
import org.springframework.beans.factory.annotation.Value;
import org.springframework.kafka.core.KafkaTemplate;
import org.springframework.stereotype.Service;

@Service
public class PaymentPublisher {
    
    private static final Logger log = LoggerFactory.getLogger(PaymentPublisher.class);

    @Autowired
    private KafkaTemplate<String, String> kafkaTemplate;
    
    @Value("${kafka.topic.produce}")
    private String producerTopic;

    public void send(String payload) {
        log.info("sending payload='{}' to topic='{}'", payload, producerTopic);
        kafkaTemplate.send(producerTopic, payload);
    }
}

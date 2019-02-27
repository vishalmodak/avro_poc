package com.loan.io.kafka;

import org.slf4j.Logger;
import org.slf4j.LoggerFactory;
import org.springframework.beans.factory.annotation.Autowired;
import org.springframework.beans.factory.annotation.Value;
import org.springframework.kafka.core.KafkaTemplate;
import org.springframework.stereotype.Service;

import com.lending.models.Loan;

@Service
public class LoanPublisher {

    private static final Logger log = LoggerFactory.getLogger(LoanPublisher.class);

    @Autowired
    private KafkaTemplate<String, Loan> kafkaTemplate;
    
    @Autowired
    private LoanRepository loanRepository;

    @Value("${loan.topic.produce}")
    private String producerTopic;

    public void send(String payload) {
        Loan loan = loanRepository.convertJsonToAvro(payload);
        log.info("sending loan='{}' to topic='{}'", loan, producerTopic);
        kafkaTemplate.send(producerTopic, loan);
    }

}

package com.loan.io.kafka;

import java.io.IOException;
import java.nio.file.Files;
import java.nio.file.Paths;
import java.util.HashMap;
import java.util.Map;

import javax.annotation.PostConstruct;

import org.apache.avro.Schema;
import org.apache.avro.generic.GenericData;
import org.slf4j.Logger;
import org.slf4j.LoggerFactory;
import org.springframework.beans.factory.annotation.Autowired;
import org.springframework.kafka.annotation.KafkaListener;
import org.springframework.stereotype.Service;

import com.jayway.jsonpath.JsonPath;

@Service
public class LoanListener {
    private static final Logger log = LoggerFactory.getLogger(LoanListener.class);
    
    @Autowired
    private LoanPublisher loanPublisher;
    
    Map<String, String> loanMap = new HashMap<>();
    
    @PostConstruct
    public void readMockData() {
        try {
            Files.walk(Paths.get("data"))
            .filter(Files::isRegularFile)
            .forEach(file -> { 
                String json;
                try {
                    json = new String(Files.readAllBytes(file));
                    String loanNumber = JsonPath.read(json, "$.loan.loanNumber");
                    loanMap.put(loanNumber, json);
                } catch (IOException e) {
                    e.printStackTrace();
                }
            });
        } catch (IOException io) {
            io.printStackTrace();
        }
    }

    @KafkaListener(topics="${kafka.topic.consume}")
    public void processMessage(String message) {
        log.info("LoanNumber: " + message);
        String loanNumber = JsonPath.read(message, "$.loan.loanNumber");
        String loan = loanMap.get(loanNumber);
        log.info(loan);
        loanPublisher.send(loan);
    }
}

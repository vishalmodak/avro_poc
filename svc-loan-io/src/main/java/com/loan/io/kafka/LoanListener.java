package com.loan.io.kafka;

import java.io.IOException;
import java.nio.file.Files;
import java.nio.file.Paths;
import java.util.HashMap;
import java.util.Map;

import javax.annotation.PostConstruct;

import org.apache.kafka.clients.consumer.ConsumerConfig;
import org.apache.kafka.common.serialization.StringDeserializer;
import org.slf4j.Logger;
import org.slf4j.LoggerFactory;
import org.springframework.beans.factory.annotation.Autowired;
import org.springframework.boot.autoconfigure.condition.ConditionalOnMissingBean;
import org.springframework.context.annotation.Bean;
import org.springframework.kafka.annotation.KafkaListener;
import org.springframework.kafka.config.ConcurrentKafkaListenerContainerFactory;
import org.springframework.kafka.core.ConsumerFactory;
import org.springframework.kafka.core.DefaultKafkaConsumerFactory;
import org.springframework.kafka.core.KafkaTemplate;
import org.springframework.stereotype.Service;

import com.fasterxml.jackson.databind.ObjectMapper;
import com.jayway.jsonpath.JsonPath;

import io.confluent.kafka.serializers.KafkaAvroDeserializer;
import io.confluent.kafka.serializers.KafkaAvroDeserializerConfig;

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
                    String loanJson = new ObjectMapper().writeValueAsString(JsonPath.read(json, "$.loan"));
                    loanMap.put(loanNumber, loanJson);  
                } catch (IOException e) {
                    e.printStackTrace();
                }
            });
        } catch (IOException io) {
            io.printStackTrace();
        }
    }

    @KafkaListener(id="loanlistener", topics="${kafka.topic.consume}", groupId = "loan", containerFactory = "loanListenerFactory")
    public void processMessage(String message) {
        log.info("LoanNumber: " + message);
        String loanNumber = JsonPath.read(message, "$.loan.loanNumber");
        String loan = loanMap.get(loanNumber);
        log.info(loan);
        loanPublisher.send(loan);
    }
    
    @Bean
    @ConditionalOnMissingBean(name="kafkaConsumerFactory")
    public ConsumerFactory<String, String> loanConsumerFactory() {
       Map<String, Object> props = new HashMap<>();
       props.put(ConsumerConfig.BOOTSTRAP_SERVERS_CONFIG, "localhost:9092");
       props.put(ConsumerConfig.KEY_DESERIALIZER_CLASS_CONFIG, StringDeserializer.class);
       props.put(ConsumerConfig.VALUE_DESERIALIZER_CLASS_CONFIG, StringDeserializer.class);
       props.put(ConsumerConfig.GROUP_ID_CONFIG, "loan");
       return new DefaultKafkaConsumerFactory<String, String>(props);
    }
    
    @Bean(name = "loanListenerFactory")
    public ConcurrentKafkaListenerContainerFactory<String, String> loanListenerFactory() {
       ConcurrentKafkaListenerContainerFactory<String, String> factory =
           new ConcurrentKafkaListenerContainerFactory<>();
       factory.setConsumerFactory(loanConsumerFactory());
//       ContainerProperties containerProperties = factory.getContainerProperties();
//       containerProperties.setPollTimeout(...);
//       containerProperties.setAckMode(AckMode...);
       return factory;
    }
}

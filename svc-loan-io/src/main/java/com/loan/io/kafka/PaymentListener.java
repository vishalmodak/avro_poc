package com.loan.io.kafka;

import java.util.HashMap;
import java.util.Map;

import org.apache.kafka.clients.consumer.ConsumerConfig;
import org.apache.kafka.clients.consumer.ConsumerRecord;
import org.apache.kafka.common.serialization.StringDeserializer;
import org.apache.kafka.common.serialization.StringSerializer;
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
import org.springframework.kafka.listener.config.ContainerProperties;
import org.springframework.stereotype.Service;

import com.lss.models.Payment;

import io.confluent.kafka.serializers.KafkaAvroDeserializer;
import io.confluent.kafka.serializers.KafkaAvroDeserializerConfig;

@Service
public class PaymentListener {
    
    private static final Logger log = LoggerFactory.getLogger(LoanListener.class);
    
    @Bean
    @ConditionalOnMissingBean(name="kafkaConsumerFactory")
    public ConsumerFactory<Object, Object> paymentConsumerFactory() {
       Map<String, Object> props = new HashMap<>();
       props.put(ConsumerConfig.BOOTSTRAP_SERVERS_CONFIG, "localhost:9092");
//       props.put(ConsumerConfig.KEY_DESERIALIZER_CLASS_CONFIG, KafkaAvroDeserializer.class);
//       props.put(ConsumerConfig.VALUE_DESERIALIZER_CLASS_CONFIG, KafkaAvroDeserializer.class);
       props.put(ConsumerConfig.KEY_DESERIALIZER_CLASS_CONFIG, StringDeserializer.class);
       props.put(ConsumerConfig.VALUE_DESERIALIZER_CLASS_CONFIG, StringDeserializer.class);
       props.put(KafkaAvroDeserializerConfig.SPECIFIC_AVRO_READER_CONFIG, true);
       props.put(KafkaAvroDeserializerConfig.SCHEMA_REGISTRY_URL_CONFIG, "http://localhost:8081");
       props.put(ConsumerConfig.GROUP_ID_CONFIG, "payment");
       return new DefaultKafkaConsumerFactory<Object, Object>(props);
    }


    @Bean(name = "paymentListenerFactory")
    public ConcurrentKafkaListenerContainerFactory<Object, Object> paymentListenerFactory() {
       ConcurrentKafkaListenerContainerFactory<Object, Object> factory =
           new ConcurrentKafkaListenerContainerFactory<>();
       factory.setConsumerFactory(paymentConsumerFactory());
//       ContainerProperties containerProperties = factory.getContainerProperties();
//       containerProperties.setPollTimeout(...);
//       containerProperties.setAckMode(AckMode...);
       return factory;
    }

    @KafkaListener(id="paymentlistener", topics="payment", groupId = "payment", containerFactory = "paymentListenerFactory")
    public void processMessage(ConsumerRecord<String, Payment> message) {
        log.info("Payment: " + message.value());
    }
    
}

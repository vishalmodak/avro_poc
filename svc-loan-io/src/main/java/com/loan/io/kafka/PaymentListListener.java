package com.loan.io.kafka;

import java.util.HashMap;
import java.util.Map;

import org.apache.kafka.clients.consumer.ConsumerConfig;
import org.apache.kafka.clients.consumer.ConsumerRecord;
import org.apache.kafka.common.serialization.StringDeserializer;
import org.slf4j.Logger;
import org.slf4j.LoggerFactory;
import org.springframework.beans.factory.annotation.Value;
import org.springframework.boot.autoconfigure.condition.ConditionalOnMissingBean;
import org.springframework.context.annotation.Bean;
import org.springframework.kafka.annotation.KafkaListener;
import org.springframework.kafka.config.ConcurrentKafkaListenerContainerFactory;
import org.springframework.kafka.core.ConsumerFactory;
import org.springframework.kafka.core.DefaultKafkaConsumerFactory;
import org.springframework.stereotype.Service;

import com.lss.models.PaymentList;

import io.confluent.kafka.serializers.KafkaAvroDeserializerConfig;

@Service
public class PaymentListListener {

    private static final Logger log = LoggerFactory.getLogger(PaymentListener.class);

    @Value("${spring.kafka.bootstrap-servers}")
    private String brokerAddress;
    @Value("${spring.kafka.consumer.properties.schema.registry.url}")
    private String schemaRegistryURL;

    @Bean
    @ConditionalOnMissingBean(name="paymentListConsumerFactory")
    public ConsumerFactory<Object, Object> paymentConsumerFactory() {
       Map<String, Object> props = new HashMap<>();
       props.put(ConsumerConfig.BOOTSTRAP_SERVERS_CONFIG, brokerAddress);
//       props.put(ConsumerConfig.KEY_DESERIALIZER_CLASS_CONFIG, KafkaAvroDeserializer.class);
//       props.put(ConsumerConfig.VALUE_DESERIALIZER_CLASS_CONFIG, KafkaAvroDeserializer.class);
       props.put(ConsumerConfig.KEY_DESERIALIZER_CLASS_CONFIG, StringDeserializer.class);
       props.put(ConsumerConfig.VALUE_DESERIALIZER_CLASS_CONFIG, StringDeserializer.class);
       props.put(KafkaAvroDeserializerConfig.SPECIFIC_AVRO_READER_CONFIG, true);
       props.put(KafkaAvroDeserializerConfig.SCHEMA_REGISTRY_URL_CONFIG, schemaRegistryURL);
       props.put(ConsumerConfig.GROUP_ID_CONFIG, "payment");
       return new DefaultKafkaConsumerFactory<>(props);
    }


    @Bean(name = "paymentListListenerFactory")
    public ConcurrentKafkaListenerContainerFactory<Object, Object> paymentListenerFactory() {
       ConcurrentKafkaListenerContainerFactory<Object, Object> factory =
           new ConcurrentKafkaListenerContainerFactory<>();
       factory.setConsumerFactory(paymentConsumerFactory());
       return factory;
    }

    @KafkaListener(id="paymentListlistener", topics="${payment.topic.consume}", groupId = "payment", containerFactory = "paymentListListenerFactory")
    public void processMessage(ConsumerRecord<String, PaymentList> message) {
        log.info("PaymentList: " + message.value());
    }

}

package com.loan.io;

import org.springframework.beans.factory.annotation.Value;
import org.springframework.boot.SpringApplication;
import org.springframework.boot.autoconfigure.SpringBootApplication;
import org.springframework.cloud.stream.schema.client.ConfluentSchemaRegistryClient;
import org.springframework.cloud.stream.schema.client.EnableSchemaRegistryClient;
import org.springframework.cloud.stream.schema.client.SchemaRegistryClient;
import org.springframework.context.annotation.Bean;
import org.springframework.context.annotation.ComponentScan;
import org.springframework.context.annotation.Configuration;
import org.springframework.context.annotation.FilterType;
import org.springframework.kafka.annotation.EnableKafka;

import com.avro.converter.AvroHttpMessageConverter;

@EnableSchemaRegistryClient
@SpringBootApplication
@EnableKafka
public class SvcLoanIoApplication {

	public static void main(String[] args) {
		SpringApplication.run(SvcLoanIoApplication.class, args);
	}
	
    @Configuration
    static class ConfluentSchemaRegistryConfiguration {
        @Bean
        public SchemaRegistryClient schemaRegistryClient(@Value("${spring.cloud.stream.schemaRegistryClient.endpoint}") String endpoint){
            ConfluentSchemaRegistryClient client = new ConfluentSchemaRegistryClient();
            client.setEndpoint(endpoint);
            return client;
        }
    }
}

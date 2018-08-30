package com.loan.io;

import org.springframework.boot.SpringApplication;
import org.springframework.boot.autoconfigure.SpringBootApplication;
import org.springframework.kafka.annotation.EnableKafka;

@SpringBootApplication
@EnableKafka
public class SvcLoanIoApplication {

	public static void main(String[] args) {
		SpringApplication.run(SvcLoanIoApplication.class, args);
	}
	

}

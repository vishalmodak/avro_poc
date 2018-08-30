package com.loan.io;

import java.io.IOException;
import java.nio.file.Files;
import java.nio.file.Paths;

import org.springframework.boot.SpringApplication;
import org.springframework.boot.autoconfigure.SpringBootApplication;
import org.springframework.kafka.annotation.EnableKafka;

import com.jayway.jsonpath.JsonPath;

@SpringBootApplication
@EnableKafka
public class SvcLoanIoApplication {

	public static void main(String[] args) {
		SpringApplication.run(SvcLoanIoApplication.class, args);
	}
	

}

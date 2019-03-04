package com.loan.io.controller;

import java.util.Arrays;
import java.util.HashMap;
import java.util.Map;

import org.slf4j.Logger;
import org.slf4j.LoggerFactory;
import org.springframework.beans.factory.annotation.Autowired;
import org.springframework.beans.factory.annotation.Value;
import org.springframework.context.annotation.Bean;
import org.springframework.http.HttpEntity;
import org.springframework.http.HttpHeaders;
import org.springframework.http.HttpMethod;
import org.springframework.http.MediaType;
import org.springframework.http.ResponseEntity;
import org.springframework.web.bind.annotation.PathVariable;
import org.springframework.web.bind.annotation.RequestHeader;
import org.springframework.web.bind.annotation.RequestMapping;
import org.springframework.web.bind.annotation.RequestParam;
import org.springframework.web.bind.annotation.ResponseBody;
import org.springframework.web.bind.annotation.RestController;
import org.springframework.web.client.RestTemplate;

import com.avro.converter.AvroHttpMessageConverter;
import com.lending.models.Loan;
import com.lending.models.PaymentList;
import com.loan.io.kafka.LoanRepository;

import io.micrometer.core.annotation.Timed;


@RestController
public class LoanController {
    private static final Logger log = LoggerFactory.getLogger(LoanController.class);
    
    @Bean
    AvroHttpMessageConverter avroHttpMessageConverter() {
        AvroHttpMessageConverter avroConverter = new AvroHttpMessageConverter<>();
        return avroConverter;
    }

    @Autowired
    private LoanRepository loanRepository;
    
    private static MediaType avroMediaType = MediaType.valueOf("avro/binary");
    private static MediaType jsonMediaType = MediaType.valueOf("application/json");
    
    @Value("${service.payments.uri}")
    private String paymentServiceURI;

    
    @RequestMapping(value="/loan/{loanNumber}", produces= {"avro/binary", "application/json"})
    @ResponseBody
    public Loan lookupLoan(@PathVariable String loanNumber) {
        Loan loan = LoanRepository.convertJsonToAvro(loanRepository.lookup(loanNumber));
        return loan;
    }
    
    @Timed(percentiles = {0.90, 0.95, 0.99}, histogram = true)
    @RequestMapping(value="/loan/payments/{loanNumber}", produces= {"application/json"})
    @ResponseBody
    public String lookuPaymentJSON(@PathVariable String loanNumber, @RequestParam(name = "dup", required = false) Integer dupCount) {
        String paymentURI = paymentServiceURI + "/v1/payments/{loanNumber}";
        if (dupCount != null && dupCount > 1) {
            paymentURI += "?dup=" + dupCount;
        }
        
        RestTemplate restTemplate = new RestTemplate();
        
        HttpHeaders outheaders = new HttpHeaders();
        outheaders.setAccept(Arrays.asList(jsonMediaType));

        HttpEntity<String> entity = new HttpEntity<String>(outheaders);
        
        Map<String, String> params = new HashMap<>();
        params.put("loanNumber", loanNumber);
        
        ResponseEntity<String> response = restTemplate.exchange(paymentURI, HttpMethod.GET, entity, String.class, params);
        log.info("Payment: " + response);
        
        return response.getBody();

    }
    
    @Timed(percentiles = {0.90, 0.95, 0.99}, histogram = true)
    @RequestMapping(value="/payments/{loanNumber}", produces={"avro/binary"})
    @ResponseBody
    public PaymentList lookuPayment(@PathVariable String loanNumber, @RequestParam(name = "dup", required = false) Integer dupCount) {
        String paymentURI = paymentServiceURI + "/v1/payments/{loanNumber}";
        if (dupCount != null && dupCount > 1) {
            paymentURI += "?dup=" + dupCount;
        }
        
        RestTemplate restTemplate = new RestTemplate();
        restTemplate.getMessageConverters().add(new AvroHttpMessageConverter<PaymentList>());
        
        HttpHeaders outheaders = new HttpHeaders();
        outheaders.setAccept(Arrays.asList(avroMediaType));

        HttpEntity<String> entity = new HttpEntity<String>(outheaders);
        
        Map<String, String> params = new HashMap<>();
        params.put("loanNumber", loanNumber);
        
        ResponseEntity<PaymentList> response = restTemplate.exchange(paymentURI, HttpMethod.GET, entity, PaymentList.class, params);
        log.info("Payment: " + response);
        
        return response.getBody();

    }
}

package com.loan.io.controller;

import java.util.Arrays;
import java.util.HashMap;
import java.util.Map;

import org.slf4j.Logger;
import org.slf4j.LoggerFactory;
import org.springframework.beans.factory.annotation.Autowired;
import org.springframework.context.annotation.Bean;
import org.springframework.http.HttpEntity;
import org.springframework.http.HttpHeaders;
import org.springframework.http.HttpMethod;
import org.springframework.http.MediaType;
import org.springframework.http.ResponseEntity;
import org.springframework.web.bind.annotation.PathVariable;
import org.springframework.web.bind.annotation.RequestMapping;
import org.springframework.web.bind.annotation.ResponseBody;
import org.springframework.web.bind.annotation.RestController;
import org.springframework.web.client.RestTemplate;

import com.loan.io.kafka.LoanRepository;
import com.lss.models.Loan;
import com.lss.models.PaymentList;


@RestController
public class LoanController {
    private static final Logger log = LoggerFactory.getLogger(LoanController.class);
    
    @Bean
    AvroHttpMessageConverter<Loan> avroHttpMessageConverter() {
        AvroHttpMessageConverter<Loan> avroConverter = new AvroHttpMessageConverter<>();
        return avroConverter;
    }

    @Autowired
    private LoanRepository loanRepository;

    
    @RequestMapping(value="/loan/{loanNumber}", produces= {"avro/binary", "application/json"})
    @ResponseBody
    public Loan lookupLoan(@PathVariable String loanNumber) {
        Loan loan = LoanRepository.convertJsonToAvro(loanRepository.lookup(loanNumber));
        return loan;
    }
    
    @RequestMapping(value="/payments/{loanNumber}", produces= {"avro/binary", "application/json"})
    @ResponseBody
    public PaymentList lookuPayment(@PathVariable String loanNumber) {
        String paymentURI = "http://localhost:3000/v1/payments/{loanNumber}";
        
        RestTemplate restTemplate = new RestTemplate();
        restTemplate.getMessageConverters().add(new AvroHttpMessageConverter<PaymentList>());
        
        HttpHeaders headers = new HttpHeaders();
        headers.setAccept(Arrays.asList((MediaType.valueOf("avro/binary"))));
        HttpEntity<String> entity = new HttpEntity<String>("parameters", headers);
        
        Map<String, String> params = new HashMap<>();
        params.put("loanNumber", loanNumber);
        
        ResponseEntity<PaymentList> response = restTemplate.exchange(paymentURI, HttpMethod.GET, entity, PaymentList.class, params);
        log.info("Payment: " + response);
        
        return response.getBody();
    }
}

package com.loan.io.kafka;

import java.io.IOException;
import java.nio.file.Files;
import java.nio.file.Paths;
import java.util.HashMap;
import java.util.Map;

import javax.annotation.PostConstruct;

import com.fasterxml.jackson.databind.ObjectMapper;
import com.jayway.jsonpath.JsonPath;
import com.lss.models.Loan;

import tech.allegro.schema.json2avro.converter.JsonAvroConverter;

import org.springframework.core.io.ClassPathResource;
import org.springframework.stereotype.Service;
import org.springframework.util.FileCopyUtils;

@Service
public class LoanRepository {
    
    public static ClassPathResource loanSchema = new ClassPathResource("avro/loan.avsc");

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
    
    public String lookup(String loanNumber) {
        return loanMap.get(loanNumber);
    }
    
    public static Loan convertJsonToAvro(String json) {
        JsonAvroConverter converter = new JsonAvroConverter();
        try {
            byte[] bdata = FileCopyUtils.copyToByteArray(loanSchema.getInputStream());
            String schema = new String(bdata);
            return converter.convertToSpecificRecord(json.getBytes(), Loan.class, schema);
        } catch (IOException e) {
            e.printStackTrace();
        }
        return null;
    }
}

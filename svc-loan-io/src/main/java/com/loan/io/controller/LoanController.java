package com.loan.io.controller;

import org.springframework.web.bind.annotation.PathVariable;
import org.springframework.web.bind.annotation.RequestMapping;
import org.springframework.web.bind.annotation.RequestMethod;
import org.springframework.web.bind.annotation.RestController;

@RestController
public class LoanController {
	
    @RequestMapping("/")
    public String index() {
        return "Loan IO Service";
    }
    
    @RequestMapping(value="/loan/lookup/{loanNumber}", method=RequestMethod.GET, produces={"application/json","application/x-protobuf"})
    public String lookupLoan(@PathVariable String loanNumber) {
        return "Loan IO Service";
    }
}

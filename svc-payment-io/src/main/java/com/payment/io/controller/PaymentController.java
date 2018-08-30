package com.payment.io.controller;

import org.springframework.web.bind.annotation.PathVariable;
import org.springframework.web.bind.annotation.RequestMapping;
import org.springframework.web.bind.annotation.RequestMethod;
import org.springframework.web.bind.annotation.RestController;

@RestController
public class PaymentController {
	
    @RequestMapping("/")
    public String index() {
        return "Payment IO Service";
    }
    
    @RequestMapping(value="/payment/lookup/{loanNumber}", method=RequestMethod.GET, produces={"application/json","application/x-protobuf"})
    public String lookupLoan(@PathVariable String loanNumber) {
        return "Payment IO Service";
    }
}

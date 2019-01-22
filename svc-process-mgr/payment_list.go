// Code generated by github.com/actgardner/gogen-avro. DO NOT EDIT.
/*
 * SOURCES:
 *     payment.avsc
 *     payment_list.avsc
 *     loan.avsc
 */

package main

import (
	"io"
)

type Payment_list struct {
	Payments []*Payment
}

func DeserializePayment_list(r io.Reader) (*Payment_list, error) {
	return readPayment_list(r)
}

func NewPayment_list() *Payment_list {
	v := &Payment_list{
		Payments: make([]*Payment, 0),
	}
	v.Payments = make([]*Payment, 0)

	return v
}

func (r *Payment_list) Schema() string {
	return "{\"fields\":[{\"default\":[],\"name\":\"payments\",\"type\":{\"items\":{\"fields\":[{\"name\":\"paid\",\"type\":\"boolean\"},{\"name\":\"datePaid\",\"type\":\"string\"},{\"name\":\"loanNumber\",\"type\":\"string\"},{\"name\":\"amountInCents\",\"type\":\"int\"},{\"name\":\"sourceAccountNumber\",\"type\":\"string\"},{\"name\":\"sourcePaymentNumber\",\"type\":\"string\"},{\"name\":\"sourceObligationNumber\",\"type\":\"string\"}],\"name\":\"payment\",\"type\":\"record\"},\"type\":\"array\"}}],\"name\":\"payment_list\",\"type\":\"record\"}"
}

func (r *Payment_list) Serialize(w io.Writer) error {
	return writePayment_list(r, w)
}
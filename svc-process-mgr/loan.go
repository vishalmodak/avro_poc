// Code generated by github.com/actgardner/gogen-avro. DO NOT EDIT.
/*
 * SOURCES:
 *     payment_list.avsc
 *     loan.avsc
 */

package main

import (
	"io"
)

type Loan struct {

	// loan status
	Status       string
	CalledDue    bool
	DateClosed   string
	DateOpened   string
	LoanNum      string
	DaysPastDue  int32
	CurrentOwner string
	TransferFrom string
}

func DeserializeLoan(r io.Reader) (*Loan, error) {
	return readLoan(r)
}

func NewLoan() *Loan {
	v := &Loan{}
	v.Status = ""
	v.CalledDue = true
	v.DateClosed = ""
	v.DateOpened = ""
	v.LoanNum = ""
	v.DaysPastDue = 0
	v.CurrentOwner = ""
	v.TransferFrom = ""

	return v
}

func (r *Loan) Schema() string {
	return "{\"fields\":[{\"default\":\"\",\"doc\":\"loan status\",\"name\":\"status\",\"type\":\"string\"},{\"default\":true,\"doc\":\"\",\"name\":\"calledDue\",\"type\":\"boolean\"},{\"default\":\"\",\"doc\":\"\",\"name\":\"dateClosed\",\"type\":\"string\"},{\"default\":\"\",\"doc\":\"\",\"name\":\"dateOpened\",\"type\":\"string\"},{\"default\":\"\",\"doc\":\"\",\"name\":\"loanNum\",\"type\":\"string\"},{\"default\":0,\"doc\":\"\",\"name\":\"daysPastDue\",\"type\":\"int\"},{\"default\":\"\",\"doc\":\"\",\"name\":\"currentOwner\",\"type\":\"string\"},{\"default\":\"\",\"doc\":\"\",\"name\":\"transferFrom\",\"type\":\"string\"}],\"name\":\"Loan\",\"namespace\":\"com.lss.models\",\"type\":\"record\",\"version\":\"1\"}"
}

func (r *Loan) Serialize(w io.Writer) error {
	return writeLoan(r, w)
}

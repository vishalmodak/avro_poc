// Code generated by github.com/actgardner/gogen-avro. DO NOT EDIT.
/*
 * SOURCE:
 *     loan.avsc
 */

package main

import (
	"fmt"
	"io"
	"math"
)

type ByteReader interface {
	ReadByte() (byte, error)
}

type ByteWriter interface {
	Grow(int)
	WriteByte(byte) error
}

type StringWriter interface {
	WriteString(string) (int, error)
}

func encodeInt(w io.Writer, byteCount int, encoded uint64) error {
	var err error
	var bb []byte
	bw, ok := w.(ByteWriter)
	// To avoid reallocations, grow capacity to the largest possible size
	// for this integer
	if ok {
		bw.Grow(byteCount)
	} else {
		bb = make([]byte, 0, byteCount)
	}

	if encoded == 0 {
		if bw != nil {
			err = bw.WriteByte(0)
			if err != nil {
				return err
			}
		} else {
			bb = append(bb, byte(0))
		}
	} else {
		for encoded > 0 {
			b := byte(encoded & 127)
			encoded = encoded >> 7
			if !(encoded == 0) {
				b |= 128
			}
			if bw != nil {
				err = bw.WriteByte(b)
				if err != nil {
					return err
				}
			} else {
				bb = append(bb, b)
			}
		}
	}
	if bw == nil {
		_, err := w.Write(bb)
		return err
	}
	return nil

}

func readBool(r io.Reader) (bool, error) {
	var b byte
	var err error
	if br, ok := r.(ByteReader); ok {
		b, err = br.ReadByte()
	} else {
		bs := make([]byte, 1)
		_, err = io.ReadFull(r, bs)
		if err != nil {
			return false, err
		}
		b = bs[0]
	}
	return b == 1, nil
}

func readInt(r io.Reader) (int32, error) {
	var v int
	buf := make([]byte, 1)
	for shift := uint(0); ; shift += 7 {
		if _, err := io.ReadFull(r, buf); err != nil {
			return 0, err
		}
		b := buf[0]
		v |= int(b&127) << shift
		if b&128 == 0 {
			break
		}
	}
	datum := (int32(v>>1) ^ -int32(v&1))
	return datum, nil
}

func readLoan(r io.Reader) (*Loan, error) {
	var str = &Loan{}
	var err error
	str.Status, err = readString(r)
	if err != nil {
		return nil, err
	}
	str.CalledDue, err = readBool(r)
	if err != nil {
		return nil, err
	}
	str.DateClosed, err = readString(r)
	if err != nil {
		return nil, err
	}
	str.DateOpened, err = readString(r)
	if err != nil {
		return nil, err
	}
	str.LoanNumber, err = readString(r)
	if err != nil {
		return nil, err
	}
	str.DaysPastDue, err = readInt(r)
	if err != nil {
		return nil, err
	}
	str.CurrentOwner, err = readString(r)
	if err != nil {
		return nil, err
	}
	str.TransferFrom, err = readString(r)
	if err != nil {
		return nil, err
	}
	str.SourceAccountNumber, err = readString(r)
	if err != nil {
		return nil, err
	}
	str.AmountPastDueInCents, err = readInt(r)
	if err != nil {
		return nil, err
	}
	str.CurrentBalanceInCents, err = readInt(r)
	if err != nil {
		return nil, err
	}
	str.LastPaymentContingent, err = readBool(r)
	if err != nil {
		return nil, err
	}
	str.TransferEffectiveDate, err = readString(r)
	if err != nil {
		return nil, err
	}
	str.CalledDueAmountInCents, err = readInt(r)
	if err != nil {
		return nil, err
	}
	str.OustandingBalanceCents, err = readInt(r)
	if err != nil {
		return nil, err
	}
	str.InCollectionBeforeClosed, err = readBool(r)
	if err != nil {
		return nil, err
	}
	str.OriginalLoanAmountInCents, err = readInt(r)
	if err != nil {
		return nil, err
	}

	return str, nil
}

func readLong(r io.Reader) (int64, error) {
	var v uint64
	buf := make([]byte, 1)
	for shift := uint(0); ; shift += 7 {
		if _, err := io.ReadFull(r, buf); err != nil {
			return 0, err
		}
		b := buf[0]
		v |= uint64(b&127) << shift
		if b&128 == 0 {
			break
		}
	}
	datum := (int64(v>>1) ^ -int64(v&1))
	return datum, nil
}

func readString(r io.Reader) (string, error) {
	len, err := readLong(r)
	if err != nil {
		return "", err
	}

	// makeslice can fail depending on available memory.
	// We arbitrarily limit string size to sane default (~2.2GB).
	if len < 0 || len > math.MaxInt32 {
		return "", fmt.Errorf("string length out of range: %d", len)
	}

	if len == 0 {
		return "", nil
	}

	bb := make([]byte, len)
	_, err = io.ReadFull(r, bb)
	if err != nil {
		return "", err
	}
	return string(bb), nil
}

func writeBool(r bool, w io.Writer) error {
	var b byte
	if r {
		b = byte(1)
	}

	var err error
	if bw, ok := w.(ByteWriter); ok {
		err = bw.WriteByte(b)
	} else {
		bb := make([]byte, 1)
		bb[0] = b
		_, err = w.Write(bb)
	}
	if err != nil {
		return err
	}
	return nil
}

func writeInt(r int32, w io.Writer) error {
	downShift := uint32(31)
	encoded := uint64((uint32(r) << 1) ^ uint32(r>>downShift))
	const maxByteSize = 5
	return encodeInt(w, maxByteSize, encoded)
}

func writeLoan(r *Loan, w io.Writer) error {
	var err error
	err = writeString(r.Status, w)
	if err != nil {
		return err
	}
	err = writeBool(r.CalledDue, w)
	if err != nil {
		return err
	}
	err = writeString(r.DateClosed, w)
	if err != nil {
		return err
	}
	err = writeString(r.DateOpened, w)
	if err != nil {
		return err
	}
	err = writeString(r.LoanNumber, w)
	if err != nil {
		return err
	}
	err = writeInt(r.DaysPastDue, w)
	if err != nil {
		return err
	}
	err = writeString(r.CurrentOwner, w)
	if err != nil {
		return err
	}
	err = writeString(r.TransferFrom, w)
	if err != nil {
		return err
	}
	err = writeString(r.SourceAccountNumber, w)
	if err != nil {
		return err
	}
	err = writeInt(r.AmountPastDueInCents, w)
	if err != nil {
		return err
	}
	err = writeInt(r.CurrentBalanceInCents, w)
	if err != nil {
		return err
	}
	err = writeBool(r.LastPaymentContingent, w)
	if err != nil {
		return err
	}
	err = writeString(r.TransferEffectiveDate, w)
	if err != nil {
		return err
	}
	err = writeInt(r.CalledDueAmountInCents, w)
	if err != nil {
		return err
	}
	err = writeInt(r.OustandingBalanceCents, w)
	if err != nil {
		return err
	}
	err = writeBool(r.InCollectionBeforeClosed, w)
	if err != nil {
		return err
	}
	err = writeInt(r.OriginalLoanAmountInCents, w)
	if err != nil {
		return err
	}

	return nil
}

func writeLong(r int64, w io.Writer) error {
	downShift := uint64(63)
	encoded := uint64((r << 1) ^ (r >> downShift))
	const maxByteSize = 10
	return encodeInt(w, maxByteSize, encoded)
}

func writeString(r string, w io.Writer) error {
	err := writeLong(int64(len(r)), w)
	if err != nil {
		return err
	}
	if sw, ok := w.(StringWriter); ok {
		_, err = sw.WriteString(r)
	} else {
		_, err = w.Write([]byte(r))
	}
	return err
}

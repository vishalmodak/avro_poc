record :payment do
  required :paid, :boolean
  required :datePaid, :string
  required :loanNumber, :string
  required :amountInCents, :int
  required :sourceAccountNumber, :string
  required :sourcePaymentNumber, :string
  required :sourceObligationNumber, :string
end

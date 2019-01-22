class PaymentListProducerWorker < ApplicationWorker
  # include Sidekiq::Worker
  sidekiq_options retry: false

  def perform

    payments = Payments_by_loanNumber["2015CA169772974"]
    data = JSON.parse(payments)
    paymentList = Avro::PaymentList.new
    data['payments'].map do |payment|
      # puts payment['loanNumber'] + "/" + payment['sourcePaymentNumber']

      paymentList.payments << Avro::Payment.new(
          paid: payment['paid'],
          datePaid: !payment['datePaid'].nil? ? payment['datePaid'] : "",
          amountInCents: payment['amountInCents'].to_i,
          loanNumber: payment['loanNumber'],
          sourceAccountNumber: payment['sourceAccountNumber'],
          sourcePaymentNumber: payment['sourcePaymentNumber'],
          sourceObligationNumber: payment['sourceObligationNumber']
      )
    end

    Payment::PaymentResponder.call(paymentList)
  end
end

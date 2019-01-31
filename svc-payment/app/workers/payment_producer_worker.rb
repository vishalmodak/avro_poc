class PaymentProducerWorker < ApplicationWorker
  # include Sidekiq::Worker
  sidekiq_options retry: false

  def perform

    payment = Avro::Payment.new(loanNumber: "2015CA169772974", paid: true, datePaid: "2015-12-23", amountInCents: 12151,
      sourceAccountNumber: "8601860", sourceObligationNumber: "11544267", sourcePaymentNumber: "2012323")
    Payment::PaymentResponder.call(payment.avro_message_value)

  end
end

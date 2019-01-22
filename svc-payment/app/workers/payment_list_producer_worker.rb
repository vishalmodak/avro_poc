class PaymentProducerWorker < ApplicationWorker
  # include Sidekiq::Worker
  sidekiq_options retry: false

  def perform

    payment = Avro::Payment.new(loanNumber: "2015CA169772974", paid: true, datePaid: "2015-12-23", amountInCents: 12151,
      sourceAccountNumber: "8601860", sourceObligationNumber: "11544267", sourcePaymentNumber: "2012323")
    Payment::PaymentResponder.call(payment)

    # payments = Payments_by_loanNumber["2015CA169772974"]
    # data = JSON.parse(payments)
    # paymentList = Avro::PaymentList.new
    # data['payments'].map do |payment|
    #   # puts payment['loanNumber'] + "/" + payment['sourcePaymentNumber']
    #
    #   paymentList.payments << Avro::Payment.new(
    #       paid: payment['paid'],
    #       datePaid: !payment['datePaid'].nil? ? payment['datePaid'] : "",
    #       amountInCents: payment['amountInCents'].to_i,
    #       loanNumber: payment['loanNumber'],
    #       sourceAccountNumber: payment['sourceAccountNumber'],
    #       sourcePaymentNumber: payment['sourcePaymentNumber'],
    #       sourceObligationNumber: payment['sourceObligationNumber']
    #   )
    # end
    #
    # Payment::PaymentResponder.call(paymentList)
  end
end

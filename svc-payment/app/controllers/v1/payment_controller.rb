module V1
  class PaymentController < ApplicationController

    # GET /payments/:payment_number
    def index

    end

    # GET /payments/:loan_number
    def lookup_payments_loan_number
      payments = Payments_by_loanNumber[params[:loan_number]]
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
      puts paymentList
      respond_to do |format|
        format.json { render json: paymentList }
        format.avro { render plain: paymentList.avro_message_value, content_type: 'avro/binary' }
      end
    end

  end
end

module V1
  class PaymentController < ApplicationController

    # GET /payments/:payment_number
    def index

    end

    # GET /payments/:loan_number
    def lookup_payments_loan_number

      metricName = "lookup_payments_loan_number"
      if request.headers["Accept"] == "avro/binary"
        metricName += ".avro"
      elsif request.headers["Accept"] == "application/json"
        metricName += ".json"
      end

      puts metricName
      StatsD.measure(metricName) do
        payments = Payments_by_loanNumber[params[:loan_number]]
        data = JSON.parse(payments)
        paymentList = Avro::PaymentList.new

        paymentsArray =[]

        dup = params[:dup]? params[:dup].to_i : 1

        dup.times do
          data['payments'].map do |payment|
            # puts payment['loanNumber'] + "/" + payment['sourcePaymentNumber']

            paymentList.payments << Avro::Payment.new(
                # paymentsArray << Avro::Payment.new(
                paid: payment['paid'],
                datePaid: !payment['datePaid'].nil? ? payment['datePaid'] : "",
                amountInCents: payment['amountInCents'].to_i,
                loanNumber: payment['loanNumber'],
                sourceAccountNumber: payment['sourceAccountNumber'],
                sourcePaymentNumber: payment['sourcePaymentNumber'],
                sourceObligationNumber: payment['sourceObligationNumber']
            )
          end
        end

        respond_to do |format|
          format.json { render json: paymentList }
          format.avro { render plain: paymentList.avro_message_value, content_type: 'avro/binary' }
        end
      end
    end
  end

end

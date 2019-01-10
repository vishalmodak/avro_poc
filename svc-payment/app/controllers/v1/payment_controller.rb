module V1
  class PaymentController < ApplicationController

    # GET /payments/:payment_number
    def index

    end

    # GET /payments/:loan_number
    def lookup_payments_loan_number
      payments = Payments_by_loanNumber[params[:loan_number]]
      render json: payments
    end

  end
end

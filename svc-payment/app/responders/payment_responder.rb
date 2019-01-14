# frozen_string_literal: true

module Payment
  class PaymentResponder < ApplicationResponder
    topic :payment

    # @param data [Hash] any data that we want to send
    # @note Will be serialized to json
    def respond(data)
      respond_to :payment, data
    end
  end
end

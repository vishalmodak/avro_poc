# frozen_string_literal: true

class LoanConsumer < ApplicationConsumer

  def consume
    Rails.logger.info params
  end

end

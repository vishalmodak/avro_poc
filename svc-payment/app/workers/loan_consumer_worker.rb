class LoanConsumerWorker < ApplicationWorker
  # include Sidekiq::Worker
  # sidekiq_options retry: false

  def perform
    Rails.logger.info "LoanConsumerWorker...."
    LoanConsumer.new.consume
  end
end

# frozen_string_literal: true

ENV['RAILS_ENV'] ||= 'development'
ENV['KARAFKA_ENV'] = ENV['RAILS_ENV']

require ::File.expand_path('../config/environment', __FILE__)
Rails.application.eager_load!

require 'active_support/core_ext/hash'

# This lines will make Karafka print to stdout like puma or unicorn
if Rails.env.development?
  Rails.logger.extend(
    ActiveSupport::Logger.broadcast(
      ActiveSupport::Logger.new($stdout)
    )
  )
end

class App < Karafka::App
  setup do |config|
    config.kafka.seed_brokers = [ENV['KAFKA_HOST'] || 'kafka://localhost:9092']
    config.client_id = 'svc-payment'
    config.logger = Rails.logger
    config.backend = :inline
    config.batch_consuming = false
    # config.parser = AvroParser
  end

  after_init do
    WaterDrop.setup { |config| config.deliver = !Karafka.env.test? }
  end

  after_init do |_config|
    Sidekiq::Logging.logger = Karafka::App.logger
  end
end

# Karafka.monitor.subscribe(Karafka::Instrumentation::Listener)

# Per topic
App.consumer_groups.draw do
  consumer_group :loan do
    topic :intake do
      consumer LoanConsumer
      parser AvroParser
      # worker LoanConsumerWorker
      # backend :sidekiq
    end
  end
end

App.boot!

# frozen_string_literal: true
require 'statsd-instrument'

StatsD.prefix = 'app.svc-payment'

case Rails.env
when 'test'
  StatsD.backend = StatsD::Instrument::Backends::NullBackend.new
when 'development', 'staging'
  StatsD.backend = StatsD::Instrument::Backends::UDPBackend.new(ENV['STATSD_SERVER'])
when 'production'
  StatsD.backend = StatsD::Instrument::Backends::UDPBackend.new('localhost:8125')
end
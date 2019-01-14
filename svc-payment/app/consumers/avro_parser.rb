# frozen_string_literal: true

class AvroParser
  def self.parse(message)
    Rails.logger.info message

    # loan = Avro::Loan.avro_message_decode(message)

    # avro = AvroTurf::Messaging.new(registry_url: "http://localhost:8081/")  
    loan = avro.decode(message)

  end
end

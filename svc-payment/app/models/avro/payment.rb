module Avro
  class Payment
    include Avromatic::Model.build(schema_name: 'payment')
  end

  Payment.register_schemas!
end

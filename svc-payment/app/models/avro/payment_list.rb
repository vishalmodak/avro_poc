module Avro
  class PaymentList
    include Avromatic::Model.build(schema_name: 'payment_list')
  end


  PaymentList.register_schemas!
end

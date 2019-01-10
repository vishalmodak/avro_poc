module Avro
  class User
    include Avromatic::Model.build(schema_name: 'user')
  end
end

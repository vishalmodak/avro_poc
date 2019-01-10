Avromatic.configure do |avro|
  avro.schema_store = Avro::Builder::SchemaStore.new(path: 'lib/avro/dsl')
  avro.registry_url = 'http://localhost:8081'
  avro.build_messaging!
end

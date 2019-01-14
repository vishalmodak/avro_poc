Avromatic.configure do |avro|
  avro.schema_store = AvroTurf::SchemaStore.new(path: 'avro/schema')
  avro.schema_registry = AvroTurf::SchemaRegistry.new('http://localhost:8081')
  avro.build_messaging!
end

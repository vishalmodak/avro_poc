module Avro
  class Loan
    # schema = Avro::Builder.build_schema do
    #   record :Loan, namespace: 'com.lending.models' do
    #     required :loanNumber, :string
    #     required :sourceAccountNumber, :string
    #   end
    # end
    include Avromatic::Model.build(schema_name: 'com.lending.models.Loan')
  end
end

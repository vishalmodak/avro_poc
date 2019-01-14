module Avro
  class Loan
    # schema = Avro::Builder.build_schema do
    #   record :Loan, namespace: 'com.lss.models' do
    #     required :loanNumber, :string
    #     required :sourceAccountNumber, :string
    #   end
    # end
    include Avromatic::Model.build(schema_name: 'com.lss.models.Loan')
  end
end

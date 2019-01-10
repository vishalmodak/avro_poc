require 'jsonpath'

Payments_by_loanNumber = Hash.new

Dir.glob("#{Rails.root}/data/**/*").each do |file|
  # Rails.logger.info file
  # file = File.join(Rails.root, 'data', '2015CA169772974.json')
  data_hash = JSON.parse(File.read(file))
  # Rails.logger.info data_hash
  loan_num = JsonPath.on(data_hash, '$.payments[0].loanNumber')
  # byebug
  Payments_by_loanNumber[loan_num.first] = data_hash
end

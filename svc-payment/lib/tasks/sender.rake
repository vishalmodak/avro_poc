require 'json'

namespace :waterdrop do
  desc 'Generate messages to kafka server'
  file = File.join(Rails.root, 'data', '2015CA169772974.json')
  data_hash = JSON.parse(File.read(file))
  puts data_hash
  task :send do
    3.times do
      message = data_hash.to_json
      WaterDrop::AsyncProducer.call(message, topic:'payment')
    end
  end
end

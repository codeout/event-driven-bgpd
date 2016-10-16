require 'kafka'
require_relative 'event_handler'

kafka = Kafka.new(seed_brokers: 'localhost:9092')
handler = EventHandler.new('10.0.0.0/8')

kafka.each_message(topic: 'edb', start_from_beginning: false) do |message|
  handler.process message.value
end

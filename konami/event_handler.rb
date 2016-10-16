require 'json'
require_relative 'gobgp_client'

class EventHandler
  KONAMI = [38, 38, 40, 40, 37, 39, 37, 39, 66, 65]

  def initialize(carrier)
    @prefix = carrier
    @gobgp = GobgpClient.new
  end

  def reset
    $stderr.puts "Reset state"
    @sequence = []
  end

  def configure_peer
    if @sequence==KONAMI
      @gobgp.update_policy default_accept: true
    else
      @gobgp.update_policy default_accept: false
    end
  end

  def process(event)
    event = JSON.parse(event)

    case event['type']
    when 'peer_state'
      process_peer_state event['value']
    when 'best_path'
      process_best_path event['value'] unless event['value']['withdrawal']
    end
  end


  private

  def process_peer_state(body)
    if body['State'] == 5  # Established
      reset
    else
      configure_peer
    end
  end

  def process_best_path(body)
    return unless prefix(body)==@prefix

    community = community(body).first
    $stderr.puts %(Received community #{community})
    @sequence << community
  end

  def prefix(body)
    body['nlri']['prefix']
  end

  def community(body)
    communities = body['attrs'].select {|a| a['type']==8 }.first&.fetch('communities') || []
    communities.map {|c| c & 0xffff }
  end
end

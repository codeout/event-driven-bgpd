$: << File.expand_path('..', __FILE__)

require 'gobgp_pb'
require 'gobgp_services_pb'

class GobgpClient
  def initialize
    @stub = Gobgpapi::GobgpApi::Stub.new('localhost:50051', :this_channel_is_insecure)
  end

  def update_policy(default_accept: true)
    default = default_accept ? :ACCEPT : :REJECT
    policy = current_policy
    policy.type = :EXPORT
    policy.default = default

    arg = Gobgpapi::ReplacePolicyAssignmentRequest.new(assignment: policy)
    @stub.replace_policy_assignment(arg)
    $stderr.puts "Changed policy default to #{default}"
  end


  private

  def current_policy
    arg = Gobgpapi::GetPolicyAssignmentRequest.new(assignment: Gobgpapi::PolicyAssignment.new(type: :EXPORT))
    @stub.get_policy_assignment(arg).assignment
  end
end

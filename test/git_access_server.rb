require "rack"
require "puma"
require "json"

class GitAccessServer
  attr_reader :url

  def initialize(responses)
    @server = Puma::Server.new(
      ->(env) do
        params   = Rack::Utils.parse_nested_query(env["QUERY_STRING"])
        user_id  = params["user"]
        response = responses.find { |r| r[:user] == user_id.to_i }

        [ response[:response], { "Content-Type" => "text/plain" }, [ response[:body] || "" ] ]
      end,
    )

    socket = @server.add_tcp_listener("127.0.0.1", 0)
    addr   = socket.addr(:numeric)
    @url   = "http://#{addr[3]}:#{addr[1]}/"

    start
  end

  def start
    @server.run(true)
  end

  def shutdown
    @server.halt
  end
end

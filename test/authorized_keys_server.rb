require "rack"
require "puma"
require "json"

class AuthorizedKeysServer
  attr_reader :url

  def initialize
    @keys = [
      { user_id: 1, keys: ["ssh-rsa AAA1...== something@example"] },
      { user_id: 2, keys: ["ssh-dsa ABC2...==", "ssh-rsa AAA3...== me@host"] },
    ]

    @server = Puma::Server.new(
      ->(_env) do
        [ 200, { "Content-Type" => "application/json" }, [ @keys.to_json ] ]
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

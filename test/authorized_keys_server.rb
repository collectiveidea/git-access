require 'rack'
require 'thin'
require 'json'

class AuthorizedKeysServer
  attr_reader :port
  attr_reader :keys

  def initialize
    @port = (rand * 10000 + 1000).to_i

    @keys = [
      {user: "1", keys: ["ssh-rsa AAA1...== something@example"]},
      {user: "2", keys: ["ssh-dsa ABC2...==", "ssh-rsa AAA3...== me@host"]}
    ]

    start
  end

  def start
    Thin::Logging.silent = true
    Thread.abort_on_exception = true

    @server_thread = Thread.new do
      Rack::Handler::Thin.run RackHandler.new(@keys), :Port => @port
    end
    sleep 1
  end

  def shutdown
    @server_thread.kill
  end

  class RackHandler
    def initialize(keys)
      @keys = keys
    end

    def call(env)
      [ 200, { "Content-Type" => "application/json" }, [ @keys.to_json ] ]
    end
  end
end

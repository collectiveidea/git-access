require 'rack'
require 'thin'
require 'json'

class GitAccessServer
  attr_reader :port
  attr_reader :keys

  def initialize(response: 200)
    @port = (rand * 10000 + 1000).to_i
    @response = response
    start
  end

  def start
    Thin::Logging.silent = true
    Thread.abort_on_exception = true

    @server_thread = Thread.new do
      Rack::Handler::Thin.run RackHandler.new(@response), :Port => @port
    end
    sleep 1
  end

  def shutdown
    @server_thread.kill
  end

  class RackHandler
    def initialize(response)
      @response_code = response
    end

    def call(env)
      [ @response_code, { "Content-Type" => "text/plain" }, [ ] ]
    end
  end
end

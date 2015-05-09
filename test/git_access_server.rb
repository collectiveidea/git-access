require 'rack'
require 'thin'
require 'json'

class GitAccessServer
  attr_reader :port
  attr_reader :keys

  def initialize(responses)
    @port = (rand * 10000 + 1000).to_i
    @responses = responses
    start
  end

  def start
    Thin::Logging.silent = true
    Thread.abort_on_exception = true

    @server_thread = Thread.new do
      Rack::Handler::Thin.run RackHandler.new(@responses), :Port => @port
    end
    sleep 1
  end

  def shutdown
    @server_thread.kill
  end

  class RackHandler
    def initialize(responses)
      @responses = responses
    end

    def call(env)
      params = Rack::Utils.parse_nested_query(env["QUERY_STRING"])
      userId = params["user"]
      response = @responses.find {|r| r[:user] == userId.to_i }

      [ response[:response] , { "Content-Type" => "text/plain" }, [ response[:body] || "" ] ]
    end
  end
end

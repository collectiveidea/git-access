require "minitest/autorun"
require "open3"

class Minitest::Test

  class Output < Struct.new(:stdout, :stderr, :status)
    def output
      stdout + stderr
    end
  end

  def call_with_opts(*opts)
    binary = File.expand_path("../../bin/git_access", __FILE__)
    Output.new(
      *Open3.capture3([binary, *opts].join(" "))
    )
  end

end

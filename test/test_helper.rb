require "minitest/autorun"
require "open3"

class Minitest::Test

  class Output < Struct.new(:stdout, :stderr, :status)
    def output
      stdout + stderr
    end
  end

  def git_access(command, params = nil)
    git_access = File.expand_path("../../bin/git_access", __FILE__)
    Output.new(
      *Open3.capture3({"SSH_ORIGINAL_COMMAND" => command}, "#{git_access} #{params}")
    )
  end

end

require "minitest/autorun"
require "open3"

class Minitest::Test
  Output = Struct.new(:stdout, :stderr, :status) do
    def output
      stdout + stderr
    end
  end

  def git_access(command, params = nil)
    git_access = File.expand_path("../git-access", __dir__)
    Output.new(
      *Open3.capture3({ "SSH_ORIGINAL_COMMAND" => command }, "#{git_access} #{params}"),
    )
  end
end

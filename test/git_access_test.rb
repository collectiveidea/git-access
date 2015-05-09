require "test_helper"
require "git_access_server"

class GitAccessTest < Minitest::Test

  def test_forwards_only_a_specified_set_of_git_commands
    server = GitAccessServer.new([{user: 4, response: 200}])
    args = "--user 4 --permission-check-url=http://localhost:#{server.port}"

    %w(git-receive-pack git-upload-pack git-upload-archive).each do |command|
      assert_match(
        /a git repository/,
        git_access("#{command} repo.git", args).output,
        "Wrong output for #{command}"
      )

      result = git_access("cat /etc/passwd", args)
      assert_equal "", result.output
      assert !result.status.success?
    end
  ensure
    server.shutdown
  end

  def test_queries_for_user_permissions_to_the_requested_repository
    server = GitAccessServer.new([{user: 2, response: 402}, {user: 4, response: 200}])
      # User 2 Does not have access, nothing runs
      result = git_access(
        "git-upload-pack 'test_repo.git'",
        "--user 2 --permission-check-url=http://localhost:#{server.port}"
      )
      assert_equal("", result.output)

      # User 4 Has access, git command is executed
      result = git_access(
        "git-upload-pack 'test_repo.git'",
        "--user 4 --permission-check-url=http://localhost:#{server.port}"
      )
      assert_match(/'test_repo.git' does not appear to be a git repository/, result.output)
  ensure
    server.shutdown
  end

end

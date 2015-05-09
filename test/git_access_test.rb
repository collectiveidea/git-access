require "test_helper"
require "git_access_server"

class GitAccessTest < Minitest::Test

  def test_forwards_only_a_specified_set_of_git_commands
    with_git_access_server do |server|
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
    end
  end

  def test_queries_for_user_permissions_to_the_requested_repository
    with_git_access_server(response: 200) do |server|
      result = git_access(
        "git-upload-pack 'test_repo.git'",
        "--user 4 --permission-check-url=http://localhost:#{server.port}"
      )

      assert_match(/'test_repo.git' does not appear to be a git repository/, result.output)
    end
  end

  def test_closes_connection_with_error_if_permission_denied
    with_git_access_server(response: 402) do |server|
      result = git_access(
        "git-upload-pack 'test_repo.git'",
        "--user 4 --permission-check-url=http://localhost:#{server.port}"
      )

      assert_equal("", result.output)
    end
  end

  private

  def with_git_access_server(server_args = {})
    server = GitAccessServer.new(**server_args)
    yield server
  ensure
    server.shutdown
  end

end

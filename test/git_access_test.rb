require "test_helper"
require "git_access_server"

class GitAccessTest < Minitest::Test
  def test_forwards_only_a_specified_set_of_git_commands
    server = GitAccessServer.new([{ user: 4, response: 200 }])
    args = "--user 4 --permission-check-url=#{server.url}"

    %w[git-receive-pack git-upload-pack].each do |command|
      assert_match(
        /a git repository/,
        git_access("#{command} repo.git", args).output,
        "Wrong output for #{command}",
      )
    end

    assert_match(
      /No such file or directory/,
      git_access("git-upload-archive repo.git", args).output,
      "Wrong output for git-upload-archive",
    )

    result = git_access("cat /etc/passwd", args)
    assert_match(/Permission denied/, result.output)
    assert !result.status.success?
  ensure
    server.shutdown
  end

  def test_queries_for_user_permissions_to_the_requested_repository
    server = GitAccessServer.new([
      { user: 2, response: 402 },
      { user: 4, response: 200, body: "path/to/real/repo" },
    ])

    # User 2 Does not have access
    result = git_access(
      "git-upload-pack 'test_repo.git'",
      "--user 2 --permission-check-url=#{server.url}",
    )
    assert_match(/Permission denied/, result.output)

    # User 4 Has access, git command is executed
    result = git_access(
      "git-upload-pack 'test_repo.git'",
      "--user 4 --permission-check-url=#{server.url}",
    )
    assert_match(%r{'path/to/real/repo' does not appear to be a git repository}, result.output)
  ensure
    server.shutdown
  end
end

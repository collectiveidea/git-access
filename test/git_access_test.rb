require "test_helper"

class GitAccessTest < Minitest::Test

  def test_forwards_git_receive_pack
    assert_match(/usage: git receive-pack/, call_with_opts("git-receive-pack").output)
  end

  def test_forwards_git_upload_pack
    assert_match(/usage: git upload-pack/, call_with_opts("git-upload-pack").output)
  end

  def test_forwards_git_upload_archive
    assert_match(/usage: git upload-archive/, call_with_opts("git-upload-archive").output)
  end

  def test_ignores_other_command_requests
    result = call_with_opts("cat", "/etc/passwd")

    assert_equal "", result.output
    assert !result.status.success?
  end

end

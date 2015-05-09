require "test_helper"

class GitAccessTest < Minitest::Test

  def test_forwards_git_receive_pack
    assert_match(/usage: git receive-pack/, git_access("git-receive-pack").output)
  end

  def test_forwards_git_upload_pack
    assert_match(/usage: git upload-pack/, git_access("git-upload-pack").output)
  end

  def test_forwards_git_upload_archive
    assert_match(/usage: git upload-archive/, git_access("git-upload-archive").output)
  end

  def test_ignores_other_command_requests
    result = git_access("cat /etc/passwd")

    assert_equal "", result.output
    assert !result.status.success?
  end

end

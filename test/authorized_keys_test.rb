require "test_helper"
require "authorized_keys_server"

class AuthorizedKeysTest < Minitest::Test

  def test_requests_authorized_keys_from_configured_url
    server = AuthorizedKeysServer.new

    result = call_with_opts("-A", "--authorized-keys-url=http://localhost:#{server.port}")

    keys = result.output.split("\n")
    assert_equal 3, keys.size
  ensure
    server.shutdown
  end

  def test_adds_authorized_keys_command_info_to_returned_keys
    server = AuthorizedKeysServer.new

    result = call_with_opts("-A", "--authorized-keys-url=http://localhost:#{server.port}")
    keys = result.output.split("\n")

    auth_keys_options = "no-user-rc,no-X11-forwarding,no-agent-forwarding,no-pty"

    keys.each_with_index do |key, i|
      assert_equal(
        %|command="git-access --user=#{i + 1}",#{auth_keys_options} #{server.keys[i]}|,
        key,
        "Invalid command line for key #{i}"
      )
    end

    assert_equal 3, keys.size
  ensure
    server.shutdown
  end

  def test_errors_if_no_url_given
    result = call_with_opts("-A")

    assert_match(/--authorized-keys-url is required/, result.output)
  end

end

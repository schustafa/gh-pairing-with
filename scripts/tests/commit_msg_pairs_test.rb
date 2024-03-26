require "minitest/autorun"
require "../commit_msg_pairs"

class PairingTest < MiniTest::Test
  def test_nobody_pairing
    message = "This is a regular commit message"
    assert_empty parse_pairing_handles(message)
  end

  def test_one_pair
    message = <<-MSG
    Does some stuff.

    Pairing with @eeyore.
    MSG

    assert_equal 1, parse_pairing_handles(message).length
    assert_equal ["@eeyore"], parse_pairing_handles(message)
  end

  def test_two_pairs_on_two_lines
    message = <<-MSG
    Exciting new functionality.

    Pairing with @pooh.
    Pairing with @tigger.
    MSG
    assert_equal 2, parse_pairing_handles(message).length
    assert_equal ["@pooh", "@tigger"], parse_pairing_handles(message)
  end

  def test_many_pairs_on_one_line
    message = <<-MSG
    We did it! Pairing with @pooh, @tigger, and @piglet.

    @eeyore didn't help at all.
    MSG

    assert_equal 3, parse_pairing_handles(message).length
    assert_equal ["@pooh", "@tigger", "@piglet"], parse_pairing_handles(message)
  end
end

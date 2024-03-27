require "minitest/autorun"
require_relative "../commit_msg_pairs"

class PairingTest < Minitest::Test
  def test_nobody_pairing
    message = "This is a regular commit message"
    assert_pairs [], message
  end

  def test_one_pair
    message = <<-MSG
    Does some stuff.

    Pairing with @eeyore.
    MSG

    assert_pairs ["eeyore"], message
  end

  def test_two_pairs_on_two_lines
    message = <<-MSG
    Exciting new functionality.

    Pairing with @pooh.
    Pairing with @tigger.
    MSG

    assert_pairs ["pooh", "tigger"], message
  end

  def test_many_pairs_on_one_line
    message = <<-MSG
    We did it! Pairing with @pooh, @tigger, and @piglet.

    @eeyore didn't help at all.
    MSG

    assert_pairs ["pooh", "tigger", "piglet"], message
  end

  def test_when_you_are_really_excited
    message = <<-MSG
    It finally worked! Pairing with @christoph3rr0bin!
    MSG

    assert_pairs ["christoph3rr0bin"], message
  end

  def test_when_you_are_not_sure_what_just_happened
    message = <<-MSG
    Is this it? Pairing with @roo?
    MSG

    assert_pairs ["roo"], message
  end

  def test_no_punctuation
    message = <<-MSG
    fixed it. pairing with @owl
    MSG

    assert_pairs ["owl"], message
  end

  def test_with_a_colon
    message = <<-MSG
    fixed it. pairing with: @owl, @roo, @eeyore
    MSG

    assert_pairs ["owl", "roo", "eeyore"], message
  end

  def test_other_wordings
    wordings = ["collaborating with", "working with"]
    wordings.each do |wording|
      message = <<-MSG
      fixed! #{wording} @owl.
      MSG

      assert_pairs ["owl"], message
    end
  end

  def assert_pairs(expected_pairs, message)
    result = parse_pairing_handles(message)
    assert_equal expected_pairs.length, result.length
    assert_equal expected_pairs, result
  end
end

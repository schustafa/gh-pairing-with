def parse_pairing_handles(commit_msg)
  # Valid pairing formats
  # * pairing with @username
  # * Pairing with @username, @username2, and @username3
  # * pairing with @username and @username2

  regex = /pairing with (@[^.\r\n]*)/i
  match = commit_msg.scan(regex)

  return [] unless match

  pairs = []
  match.flatten.each do |substring|
    substring.split do |word|
      next unless word.start_with?("@")
      word.gsub!(/^@/, "")
      word.gsub!(/[,;.]$/i, "")
      pairs << word unless pairs.include?(word)
    end
  end

  pairs
end

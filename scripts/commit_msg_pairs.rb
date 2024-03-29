def parse_pairing_handles(commit_msg)
  descriptors = [
    "pairing with",
    "collaborating with",
    "working with"
  ]

  regex = /(?:#{descriptors.join("|")}):? (@[^.\r\n]*)/i
  match = commit_msg.scan(regex)

  return [] unless match

  pairs = []
  match.flatten.each do |substring|
    substring.split do |word|
      next unless word.start_with?("@")
      word.gsub!(/^@/, "")
      word.gsub!(/[,;.!?]$/i, "")
      pairs << word unless pairs.include?(word)
    end
  end

  pairs
end

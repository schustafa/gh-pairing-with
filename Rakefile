require "rake/testtask"

task default: "test"

Rake::TestTask.new do |task|
  task.libs = ["scripts/tests"]
  task.test_files = FileList["scripts/tests/*_test.rb"]
  task.options = "--pride"
end

file "githooks/commit-msg" do
  sh "erb githooks/commit-msg.TEMPLATE.erb > githooks/commit-msg && chmod +x githooks/commit-msg"
  ruby "-c githooks/commit-msg"
end

desc "Generate the commit-msg hook and verify its syntax"
task generate_hook: %w[githooks/commit-msg]

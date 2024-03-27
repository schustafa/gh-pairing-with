require "rake/testtask"

task default: "test"

Rake::TestTask.new do |task|
  task.libs = ["scripts/tests"]
  task.test_files = FileList["scripts/tests/*_test.rb"]
  task.options = "--pride"
end

desc "Generate the commit-msg hook and verify its syntax"
file "generate-hook" do
  sh "erb githooks/commit-msg.TEMPLATE.erb > githooks/commit-msg"
  ruby "-c githooks/commit-msg"
end

require "rake/testtask"

task default: "test"

Rake::TestTask.new do |task|
  task.libs = ["scripts/tests"]
  task.test_files = FileList["scripts/tests/*_test.rb"]
  task.options = "--pride"
end

require 'rake/testtask'

task :build do
  rm_f "bin/git-access"
  sh "gb build"
end

Rake::TestTask.new do |t|
  t.libs = %w(test)
  t.pattern = "test/*_test.rb"
end

task default: [:build, :test]

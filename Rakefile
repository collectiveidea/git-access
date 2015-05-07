require 'rake/testtask'

task :build do
  sh "gb build"
end

Rake::TestTask.new do |t|
  t.pattern = "test/*_test.rb"
end

task default: [:build, :test]

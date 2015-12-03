require 'rake/testtask'

task :build do
  rm_f "bin/git-access"
  sh "gb build"
end

desc "Build cross-compiled binaries"
task :release do
  sh "GOOS=linux GOARCH=amd64 gb build"
end

Rake::TestTask.new do |t|
  t.libs = %w(test)
  t.pattern = "test/*_test.rb"
end

task default: [:build, :test]

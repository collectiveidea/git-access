require 'rake/testtask'

task :build do
  sh "go build ./src/git-access"
end

desc "Build cross-compiled binaries"
task :release do
  sh "GOOS=linux GOARCH=amd64 go build ./src/git-access"
end

task :vet do
  sh "go vet ./..."
end

Rake::TestTask.new do |t|
  t.libs = %w(test)
  t.pattern = "test/*_test.rb"
end

task default: [:build, :test, :vet]

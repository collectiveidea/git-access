require "rake/testtask"

task :build do
  rm "./git-access"
  sh "go build ./cmd/git-access"
end

desc "Build cross-compiled binaries"
task :release do
  sh "GOOS=linux GOARCH=amd64 go build ./cmd/git-access"
end

task :vet do
  sh "go vet ./..."
end

Rake::TestTask.new do |t|
  t.libs = %w[test]
  t.pattern = "test/*_test.rb"
end

task default: %i[build test vet]

task :default => :launch_server

task :install do
  #Should probably use some form of package maangement
  `go get github.com/revel/revel`
  `go get github.com/revel/cmd/revel`
  `go get github.com/boltdb/bolt/...`
  `go build tools/initdb.go`
end

task :init_db do
  `./initdb` if !File.exist?('db/HumanPredictions.db')
end

task :launch_server, [:threads] => :init_db do |t, args|
  max_procs = ENV['GOMAXPROCS']
  args.with_defaults(:threads => max_procs || 4)
  ENV['GOMAXPROCS'] = args[:threads]

  `revel run github.com/nickjanus/ProteinGraphQuery`
end

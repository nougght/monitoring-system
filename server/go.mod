module github.com/nougght/monitoring-system/server

go 1.26.0

require (
	github.com/google/uuid v1.6.0
	github.com/nougght/monitoring-system/shared/go v0.0.0
	google.golang.org/grpc v1.82.0
)

require (
	golang.org/x/net v0.53.0 // indirect
	golang.org/x/sys v0.43.0 // indirect
	golang.org/x/text v0.36.0 // indirect
	google.golang.org/genproto/googleapis/rpc v0.0.0-20260414002931-afd174a4e478 // indirect
	google.golang.org/protobuf v1.36.11 // indirect
)

replace github.com/nougght/monitoring-system/shared/go => ../shared/go

package proto

//go:generate protoc -I. -I../../../third_party --gofast_out=. --gofast_opt=paths=source_relative --zmicro-gin_out=. --zmicro-gin_opt=paths=source_relative --zmicro-gin_opt allow_empty_patch_body=true hello.proto

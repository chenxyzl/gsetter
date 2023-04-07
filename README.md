# gsetter

# build
go build -o ~/go/bin/protoc-gen-go-dirty main.go

# gen
protoc --go-dirty_out=. --go-dirty_opt=paths=source_relative example.proto
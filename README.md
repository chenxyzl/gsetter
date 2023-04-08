# gsetter

# build
go build -o ~/go/bin/protoc-gen-go-dirty main.go

# gen pb
protoc --go_out=. --proto_path=./ example.proto

# gen dirty
protoc --go-dirty_out=. --go-dirty_opt=paths=source_relative example.proto
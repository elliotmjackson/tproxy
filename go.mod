module github.com/elliotmjackson/tproxy

go 1.18

require (
	buf.build/gen/go/bufbuild/eliza/bufbuild/connect-go v1.4.0-20221108170037-30afbf7c670d.1
	buf.build/gen/go/bufbuild/eliza/protocolbuffers/go v1.28.1-20221108170037-30afbf7c670d.4
	github.com/bufbuild/connect-go v1.4.0
	github.com/fatih/color v1.13.0
	golang.org/x/net v0.4.0
	google.golang.org/protobuf v1.28.2-0.20220831092852-f930b1dc76e8
)

require (
	github.com/mattn/go-colorable v0.1.9 // indirect
	github.com/mattn/go-isatty v0.0.14 // indirect
	golang.org/x/sys v0.3.0 // indirect
	golang.org/x/text v0.5.0 // indirect
)

replace github.com/bufbuild/buf v1.10.0 => ../buf

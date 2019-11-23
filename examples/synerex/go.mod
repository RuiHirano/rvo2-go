module github.com/synerex/synerex_alpha/provider/simulation/rvo2-go/examples/synerex

require (
	github.com/RuiHirano/rvo2-go v0.0.0-20191123120741-33bd53a271f0 // indirect
	github.com/synerex/synerex_alpha/api v0.0.0
	github.com/synerex/synerex_alpha/api/fleet v0.0.0
	github.com/synerex/synerex_alpha/api/simulation/area v0.0.0
	github.com/synerex/synerex_alpha/api/simulation/clock v0.0.0
	github.com/synerex/synerex_alpha/provider/simulation/rvo2-go/src/rvosimulator v0.0.0-00010101000000-000000000000 // indirect
	github.com/synerex/synerex_alpha/provider/simulation/simutil v0.0.0
	github.com/synerex/synerex_alpha/sxutil v0.0.0
	google.golang.org/grpc v1.22.1
)

replace (
	github.com/synerex/synerex_alpha/api => ./../../../../../api
	github.com/synerex/synerex_alpha/api/adservice => ../../../../../api/adservice
	github.com/synerex/synerex_alpha/api/common => ./../../../../../api/common
	github.com/synerex/synerex_alpha/api/fleet => ../../../../../api/fleet
	github.com/synerex/synerex_alpha/api/library => ../../../../../api/library
	github.com/synerex/synerex_alpha/api/ptransit => ../../../../../api/ptransit
	github.com/synerex/synerex_alpha/api/rideshare => ../../../../../api/rideshare
	github.com/synerex/synerex_alpha/api/routing => ../../../../../api/routing
	github.com/synerex/synerex_alpha/api/simulation/agent => ./../../../../../api/simulation/agent
	github.com/synerex/synerex_alpha/api/simulation/area => ./../../../../../api/simulation/area
	github.com/synerex/synerex_alpha/api/simulation/clock => ./../../../../../api/simulation/clock
	github.com/synerex/synerex_alpha/api/simulation/participant => ./../../../../../api/simulation/participant
	github.com/synerex/synerex_alpha/api/simulation/route => ./../../../../../api/simulation/route
	github.com/synerex/synerex_alpha/nodeapi => ./../../../../../nodeapi
	github.com/synerex/synerex_alpha/provider/simulation/rvo2-go/src/rvosimulator => ../../src/rvosimulator
	github.com/synerex/synerex_alpha/provider/simulation/simutil => ../../../simutil
	github.com/synerex/synerex_alpha/sxutil => ./../../../../../sxutil
)

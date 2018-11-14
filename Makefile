oopt:
	docker build -t oopt .
	docker run -v `pwd`:/data --rm oopt cp /go/bin/oopt /data/

model:
	generator -compress_paths -generate_fakeroot -package_name model_gnmi -exclude_modules openconfig-lldp,ietf-interfaces -path ./submodules/public/release,./submodules/pyang/modules/ietf,./submodules/pyang/modules/iana  -output_file ./pkg/model_gnmi/packet-transport.go ./submodules/public/release/models/optical-transport/openconfig-terminal-device.yang ./submodules/public/release/models/platform/openconfig-platform.yang
	gofmt -w ./pkg/model_gnmi/packet-transport.go
	generator -compress_paths -generate_fakeroot -package_name model -exclude_modules openconfig-platform,openconfig-terminal-device,openconfig-interfaces,ietf-interfaces -path ./submodules/public/release,./submodules/pyang/modules/ietf,./submodules/pyang/modules/iana  -output_file ./pkg/model/packet-transport.go ./yang/packet-transport.yang
	gofmt -w ./pkg/model/packet-transport.go

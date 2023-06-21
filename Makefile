all: bin/terraform-provider-ccx

bin/terraform-provider-ccx:
	 go build -o ./bin/terraform-provider-ccx ./cmd/terraform-provider-ccx

clean:
	rm -rf ./bin/terraform-provider-ccx

install: bin/terraform-provider-ccx
	mkdir -p ~/.terraform.d/plugins/registry.terraform.io/severalnines/ccx/1.5.0/linux_amd64
	cp ./bin/terraform-provider-ccx ~/.terraform.d/plugins/registry.terraform.io/severalnines/ccx/1.5.0/linux_amd64/

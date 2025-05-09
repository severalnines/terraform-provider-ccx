.PHONY: docs

VERSION=0.4.6

all: bin/terraform-provider-ccx

bin/terraform-provider-ccx:
	 go build -o ./bin/terraform-provider-ccx .

clean:
	rm -rf ./bin/terraform-provider-ccx

install: bin/terraform-provider-ccx
	mkdir -p ~/.terraform.d/plugins/registry.terraform.io/severalnines/ccx/${VERSION}/linux_amd64
	cp ./bin/terraform-provider-ccx ~/.terraform.d/plugins/registry.terraform.io/severalnines/ccx/${VERSION}/linux_amd64/

docs:
	go run github.com/hashicorp/terraform-plugin-docs/cmd/tfplugindocs


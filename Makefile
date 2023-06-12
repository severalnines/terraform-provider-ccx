all: bin/terraform-provider-ccx

bin/terraform-provider-ccx:
	 go build -o ./bin/terraform-provider-ccx ./cmd/terraform-provider-ccx

clean:
	rm -rf ./bin/terraform-provider-ccx
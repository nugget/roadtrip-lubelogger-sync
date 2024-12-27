BINARYNAME=rt2ll

fetchdata:
	curl --output ./data/dropbox.zip --location "https://www.dropbox.com/scl/fo/bhu66tpp1f1rbth9ry855/AGfHN9BGzzypwMcDtlySXQc?rlkey=66a2wjdthk9v9gvouw1kmb14f&st=17g3t3vx&dl=0"
	cd data && unzip -o dropbox.zip

localdev:
	go mod edit -replace=github.com/nugget/roadtrip="/Users/nugget/src/Vehicle Fleet/roadtrip-api-go"
	go mod tidy

productiondev:
	go mod edit -dropreplace=github.com/nugget/roadtrip
	go mod tidy

testrun: localdev
	clearbuffer && go mod tidy && go build -o dist/$(BINARYNAME)
	./dist/$(BINARYNAME)

linux: productiondev
	env GOOS=linux GOARCH=amd64 go build -o dist/$(BINARYNAME)-linux-amd64

prod: productiondev linux
	scp -rp dist/$(BINARYNAME)-linux-amd64 roadtrip-sync:.local/bin/


prodrun: prod
	ssh roadtrip-sync ".local/bin/rt2ll-linux-amd64 -csvpath \"/home/nugget/Dropbox/Road Trip Data/CSV\""
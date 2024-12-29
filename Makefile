BINARYNAME=rt2ll

DATADIR=./testdata

moduleupdate:
	 go get -u github.com/nugget/roadtrip-go/roadtrip
	 go mod tidy
	 git commit go.mod go.sum -m "make moduleupdate"

fetchdata:
	curl --output $(DATADIR)/dropbox.zip --location "https://www.dropbox.com/scl/fo/bhu66tpp1f1rbth9ry855/AGfHN9BGzzypwMcDtlySXQc?rlkey=66a2wjdthk9v9gvouw1kmb14f&st=17g3t3vx&dl=0"
	cd $(DATADIR) && unzip -o dropbox.zip

localdev:
	go mod edit -replace=github.com/nugget/roadtrip-go/roadtrip="/Users/nugget/src/Vehicle Fleet/roadtrip-go/roadtrip"
	go mod tidy

productiondev:
	go mod edit -dropreplace=github.com/nugget/roadtrip-go/roadtrip
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

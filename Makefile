BINARYNAME=rt2ll

DATADIR=./testdata
PRODUCTION_HOST=roadtrip-sync

mod:
	 go get -u github.com/nugget/roadtrip-go/roadtrip
	 go get -u
	 go mod tidy
	 git commit go.mod go.sum -m "make mod"

fetchdata:
	@echo "Fetching current Dropbox files from $(PRODUCTION_HOST)"
	rsync -auz "$(PRODUCTION_HOST):Dropbox/Road\ Trip\ Data/*" "$(DATADIR)"

localdev:
	go mod edit -replace=github.com/nugget/roadtrip-go/roadtrip="/Users/nugget/src/Vehicle Fleet/roadtrip-go/roadtrip"
	go mod tidy

productiondev:
	go mod edit -dropreplace=github.com/nugget/roadtrip-go/roadtrip
	go mod tidy

testrun: 
	clearbuffer && go mod tidy && go build -o dist/$(BINARYNAME)
	./dist/$(BINARYNAME) -v

linux: productiondev
	env GOOS=linux GOARCH=amd64 go build -o dist/$(BINARYNAME)-linux-amd64

prod: productiondev linux
	scp -rp dist/$(BINARYNAME)-linux-amd64 $(PRODUCTION_HOST):.local/bin/


prodrun: prod
	ssh roadtrip-sync ".local/bin/rt2ll-linux-amd64 -csvpath \"/home/nugget/Dropbox/Road Trip Data/CSV\""

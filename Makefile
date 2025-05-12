build: clean
	mkdir out
	cp -r web out/web
	go build -o out/consumables ./cmd/web

clean:
	rm -rf out

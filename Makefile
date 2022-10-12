test: 
	go test ./...

clean:
	go clean -testcache

db:
# Make sure you have psql installed and a postgres instance is running on port 5432.
	cd repository && ./db.sh

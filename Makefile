run:
	go build ./
	./tsukemen

reset-db:
	sh script/reset_db.sh

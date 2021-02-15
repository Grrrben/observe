select-observations:
	docker-compose exec db psql -U postgres observe -c "SELECT * FROM observation;"
select-projects:
	docker-compose exec db psql -U postgres observe -c "SELECT * FROM project;"

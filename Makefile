select-observations:
	docker-compose exec db psql -U postgres observe -c "SELECT * FROM observation;"
select-projects:
	docker-compose exec db psql -U postgres observe -c "SELECT * FROM project;"

remove-images:
	sudo rm -rf ./static/images/observation/thumb/*.png
	sudo rm -rf ./static/images/observation/small/*.png
	sudo rm -rf ./static/images/observation/medium/*.png
	sudo rm -rf ./static/images/observation/large/*.png
	sudo rm -rf ./static/images/observation/raw/*.png
	@echo ----------------------------
	@echo Don\'t forget to flush the db
select-observations:
	docker-compose exec db psql -U postgres observe -c "SELECT * FROM observation;"
select-projects:
	docker-compose exec db psql -U postgres observe -c "SELECT * FROM project;"

remove-images:
	sudo rm ./static/images/observation/thumb/*.png
	sudo rm ./static/images/observation/small/*.png
	sudo rm ./static/images/observation/medium/*.png
	sudo rm ./static/images/observation/large/*.png
	sudo rm ./static/images/observation/raw/*.png
	@echo ----------------------------
	@echo Don\'t forget to flush the db
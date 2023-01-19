bash:
	docker-compose run --rm app bash

clean:
	docker-compose kill
	docker-compose rm -f

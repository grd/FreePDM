CONTAINER_NAME := freepdm
CONTAINER_DATA_FOLDER := /home/user/docker/freepdm

full: stop rm start

start restart docker: stop rm
################	docker network create FreePDM
# rmlog
#	mkdir -p ${CONTAINER_DATA_FOLDER}
#	chown -Rv 1000:1000 ${CONTAINER_DATA_FOLDER}
	docker-compose -p ${CONTAINER_NAME} up -d --build
        #$(call showlog, "It works, yey!")

stop:
	docker stop ${CONTAINER_NAME} || true

rm: stop
	docker rm ${CONTAINER_NAME} || true

#rmdata: stop
#	if [ -d "${CONTAINER_DATA_FOLDER}" ]; then rm -R -I "${CONTAINER_DATA_FOLDER}"; fi || true

#rmlog:
# https://stackoverflow.com/questions/42510002/docker-how-to-clear-the-logs-properly-for-a-docker-container
#   truncate -s 0 $(docker inspect --format='{{.LogPath}}' ${CONTAINER_NAME})


shell:
	docker exec -it ${CONTAINER_NAME} /bin/sh


logs:
	docker logs FreePDM



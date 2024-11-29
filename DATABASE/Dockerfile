FROM mysql:latest
ENV MYSQL_ROOT_PASSWORD=root
# ENV MYSQL_ROOT_PASSWORD=root \
#     DB_HOST=mysql_container:3306 \
#     DB_USER=root \
#     DB_PASSWORD=A7med2452003 \
#     DB_NAME=package_tracking_system
EXPOSE 3306

COPY ./Query.sql /docker-entrypoint-initdb.d/

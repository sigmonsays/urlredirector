
drop dead simple url redirector service which uses redis as a backend

create a new redirect with a post request

     curl -d '{ "id": "/bbq", "url":"http://grepped.org/" }' localhost:8080/api/create

get said redirect

     curl localhost:8080/bbq

build docker container

     docker build -t urlredirectord .


run container

     docker run --name redis1 -d redis
     docker run --name urlredir1 -d --link redis1 urlredirectord


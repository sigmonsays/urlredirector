
drop dead simple url redirector service

create a new redirect with a post request

     curl -d '{ "id": "/bbq", "url":"http://grepped.org/" }' localhost:8080/api/create

get said redirect

     curl localhost:8080/bbq

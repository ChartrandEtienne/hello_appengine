#!/bin/sh

rm cookie_file

echo "\n\ndeleting contents of database\n\n"

curl http://localhost:8080/dev/clearDatabase

sleep 1

echo "\n\ntrying out /signup\n\n"

curl -H "Content-Type: application/json" -X POST -d '{"name":"Doob","password":"billionsofstars"}' http://localhost:8080/signup

sleep 1

echo "\n\ntrying out the session stuff\n\n"

# curl writes to file with option --cookie-jar and reads from it with option --cookie.
# browsers read/write/create/whatever
# all automagically
# apparently
curl --cookie-jar cookie_file -H "Content-Type: application/json" -X POST -d '{"name":"Doob","password":"billionsofstars"}' http://localhost:8080/login

echo "\n\nam I logged in? \n\n"

curl --cookie cookie_file -H "Content-Type: application/json" http://localhost:8080/

echo "\n\nlogout\n\n"

curl --cookie cookie_file -H "Content-Type: application/json" http://localhost:8080/logout

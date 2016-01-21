#!/bin/sh

rm cookie_file

echo "\n\ndeleting contents of database\n\n"

curl http://localhost:8080/dev/clearDatabase

sleep 1

echo "testing /\n"
echo "should get an error about not being logged\n\n"

curl http://localhost:8080/

# TODO: add the situation where an erroneous cookie exists 

echo "\n\ntesting /signup\n"
echo "should get an error about: not a POST\n\n"

curl http://localhost:8080/signup

echo "\n\nshould get an error about name and password being mandatory\n\n"

curl -H "Content-Type: application/json" -X POST -d '{"name":"Doob"}' http://localhost:8080/signup

echo "\n\nshould get an error about name and password being mandatory again\n\n"

curl -H "Content-Type: application/json" -X POST -d '{"password":"billionsofstars"}' http://localhost:8080/signup

echo "\n\nshould return status okay\n\n"

curl -H "Content-Type: application/json" -X POST -d '{"name":"Doob","password":"billionsofstars"}' http://localhost:8080/signup

echo "\n\nshould return error about existing user"

curl -H "Content-Type: application/json" -X POST -d '{"name":"Doob","password":"billionsofstars"}' http://localhost:8080/signup

echo "\n\ntesting login\n"

echo "should get an error about name and password being mandatory\n\n"

curl --cookie-jar cookie_file --cookie cookie_file -H "Content-Type: application/json" -X POST -d '{"name":"Doob"}' http://localhost:8080/login

echo "\n\nshould get an error about name and password being mandatory again\n\n"

curl -H "Content-Type: application/json" -X POST -d '{"password":"billionsofstars"}' http://localhost:8080/login

echo "\n\nshould not have logged in the user\n\n"

curl --cookie-jar cookie_file --cookie cookie_file http://localhost:8080/

echo "\n\nshould return status okay\n\n"

curl --cookie-jar cookie_file --cookie cookie_file -H "Content-Type: application/json" -X POST -d '{"name":"Doob","password":"billionsofstars"}' http://localhost:8080/login


sleep 1

echo "\n\ntrying out the session stuff\n\n"

# curl writes to file with option --cookie-jar and reads from it with option --cookie.
# browsers read/write/create/whatever
# all automagically
# apparently
curl --cookie-jar cookie_file -H "Content-Type: application/json" -X POST -d '{"name":"Doob","password":"billionsofstars"}' http://localhost:8080/login

echo "\n\nam I logged in? \n\n"

curl --cookie cookie_file --cookie-jar cookie_file -H "Content-Type: application/json" http://localhost:8080/

echo "\n\nlogout\n\n"

curl --cookie-jar cookie_file -H "Content-Type: application/json" http://localhost:8080/logout

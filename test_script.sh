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
echo "should return status okay\n\n"

curl -H "Content-Type: application/json" -X POST -d '{"name":"DeGrasse","password":"billionsofstars"}' http://localhost:8080/signup

sleep 1

echo "\n\nshould return error about existing user\n\n"

curl -H "Content-Type: application/json" -X POST -d '{"name":"DeGrasse","password":"billionsofstars"}' http://localhost:8080/signup

echo "\n\ntesting login\n"

echo "should return status okay\n\n"

curl --cookie-jar cookie_file -H "Content-Type: application/json" -X POST -d '{"name":"DeGrasse","password":"billionsofstars"}' http://localhost:8080/login


sleep 1

echo "\n\ntrying out the session stuff\n\n"

# curl writes to file with option --cookie-jar and reads from it with option --cookie.
# browsers read/write/create/whatever
# all automagically
# apparently

echo "\n\nshould return the username of the current logged in user\n\n"

curl --cookie cookie_file --cookie-jar cookie_file -H "Content-Type: application/json" http://localhost:8080/

echo "\n\nlogout\n\n"

curl --cookie-jar cookie_file -H "Content-Type: application/json" http://localhost:8080/logout

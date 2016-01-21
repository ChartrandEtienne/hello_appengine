A very simple chat application. I think this would cover most of the basic things we do in Vidao. This project has two parts. We can go one by one. 

1) Simple authentication. You need to create 4 API.
    
    /              -  show current user's login info
    /signup   -  take username (unique) and password and create an account for the user. encrypt the password while storing (md5 checksum is enough)
    /login      - log in the user with provided username and password
    /logout    - clear the current session of the user

Try to use "datastore" as database. Use cookies to maintain session. you can keep as simple as possible. 

package hello

import (
  "fmt"
  "net/http"
  "appengine"
  "appengine/datastore"
  "io/ioutil"
  "encoding/json"
  "errors"
)

type UserCredentials struct {
  // FieldA int    `json:"field_a"
  // http://stackoverflow.com/questions/11693865/lower-case-key-names-with-json-marshal-in-go
  // that's pretty wacky
  Name string       `json:"name"`
  Password string   `json:"password"`
}

func init() {
  // apparently this handler gets called as a default.
  // it kept getting called for favicon.ico.
  http.HandleFunc("/", rootHandler)
  http.HandleFunc("/favicon.ico", func (w http.ResponseWriter, r *http.Request) {
    http.ServeFile(w, r, "favicon.ico")
  })
  http.HandleFunc("/signup", signupHandler)
  // goapp states that this handler is now default
  // for some reason
  // TODO: how does this lib picks the default handler?
  http.HandleFunc("/login", loginHandler)
  http.HandleFunc("/logout", logoutHandler)

  http.HandleFunc("/dev/clearDatabase", clearDatabaseHandler)
}

func rootHandler(writer http.ResponseWriter, request *http.Request) {
  context := appengine.NewContext(request)
  user_credentials, err := getUserFromCookie(request, context)
  if err != nil {
    context.Debugf("logged in err ", err)
    fmt.Fprint(writer, `{ "error": "you are not currently logged in"}`)
    return
  }
  fmt.Fprint(writer, `{ "username": "` + user_credentials.Name + `"}`)
}

func loginHandler(writer http.ResponseWriter, request *http.Request) {
  context := appengine.NewContext(request)
  user, err := readUserCredentialsBody(writer, request)
  if err != nil {
    return
  }
  query := datastore.NewQuery("UserCredentials").
    Filter("Name=", user.Name).
    Ancestor(userKey(context))
    // Filter("Password=", user.Password)
  var found_users []UserCredentials
  keys, err := query.GetAll(context, &found_users)
  if err != nil {
    // Handle error.
    return
  }
  if len(found_users) == 0 {
    fmt.Fprint(writer, `{ "error": "this username/password pair does not match an existing user"}`)
  }
  if len(found_users) == 1 {
    cookie := http.Cookie{Name: "my_auth_cookie", Value: keys[0].Encode(), MaxAge: 0, Secure: false, HttpOnly: false}
    http.SetCookie(writer, &cookie)
    fmt.Fprint(writer, `{ "status": "okay"}`)
  }
  // if len(found_users) > 1 // can't happen
}

func logoutHandler(writer http.ResponseWriter, request *http.Request) {
  cookie := http.Cookie{Name: "my_auth_cookie", Value: "loggedout", MaxAge: 0, Secure: false, HttpOnly: false}
  http.SetCookie(writer, &cookie)
  fmt.Fprint(writer, ` { "status": "okay" } `)
}

func signupHandler(writer http.ResponseWriter, request *http.Request) {
  context := appengine.NewContext(request)
  user, err := readUserCredentialsBody(writer, request)
  if err != nil {
    return
  }
  query := datastore.NewQuery("UserCredentials").
    Filter("Name=", user.Name)
    // Filter("Password=", user.Password)
  var found_users []UserCredentials
  if _, err = query.GetAll(context, &found_users); err != nil {
    // Handle error.
    return
  }
  if len(found_users) > 0 {
    fmt.Fprint(writer, `{ "error": "username [username] already exists"}`)
    return
  }
  key := datastore.NewIncompleteKey(context, "UserCredentials", userKey(context))
  key, err = datastore.Put(context, key, user)
  if err != nil {
    fmt.Fprint(writer, `{ "error": Put returned err for reason" } `, err)
  }
  context.Debugf("inserted key: ", key)
  fmt.Fprint(writer, `{ "status": "okay"}`)
}

func readUserCredentialsBody(writer http.ResponseWriter, request *http.Request) (*UserCredentials, error) {
  if nil == request.Body {
    fmt.Fprint(writer, `{ "error": "must be POST request"}`)
    return nil, errors.New("must be a POST request")
  }
  body, err := ioutil.ReadAll(request.Body)
  if nil != err {
    fmt.Fprint(writer, `{ "error": "not too sure what happened"}`)
    return nil, errors.New("not too sure")
  }
  var user UserCredentials
  err = json.Unmarshal(body, &user)
  if nil != err {
    fmt.Fprint(writer, `{ "error": "keys username: string and password: string must be set"}`)
    return nil, errors.New("keys username: string and password: string must be set")
  }
  return &user, nil
}


func getUserFromCookie(request *http.Request, context appengine.Context) (*UserCredentials, error) {
  cookie, err := request.Cookie("my_auth_cookie")
  if err != nil {
    return nil, err
  }
  var user_credentials UserCredentials

  context.Debugf("key from cookie ", cookie.Value)
  key, err := datastore.DecodeKey(cookie.Value)
  if err != nil {
    return nil, err
  }
  err = datastore.Get(context, key, &user_credentials)
  if err != nil {
    return nil, err
  }
  return &user_credentials, nil
}

// key for all User entries
func userKey(c appengine.Context) *datastore.Key {
  // The string "default_guestbook" here could be varied to have multiple guestbooks.
  // nice thx
  return datastore.NewKey(c, "UserCredentials", "default_guestbook", 0, nil)
}



func clearDatabaseHandler(writer http.ResponseWriter, request *http.Request) {
  context := appengine.NewContext(request)
  // query := datastore.NewQuery("UserCredentials").Ancestor(userKey(context))
  query := datastore.NewQuery("UserCredentials")
  for iterator := query.Run(context); ; {
    var user UserCredentials
    key, err := iterator.Next(&user)
    if err == datastore.Done {
      break
    }
    context.Debugf("deleting ", user, "\n")
    err = datastore.Delete(context, key)
  }
}

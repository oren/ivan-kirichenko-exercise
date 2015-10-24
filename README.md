# ivan-kirichenko-exercise
short coding exercise for Ivan Kirechenko

Additional question may be provided via issues to this repo. Good luck and have fun! :)

###Create HTTP Rest API:
1. Use echo for web handler 
2. Implement login endpoint with JWT token, and simple middleware that checks header for 'Authorization: Bearer %jwt_token%' in each request. Otherwise return 403 and json struct with error 
3. Implement endpoint that will use OAuth2 authorization for FB, to login and issue access_token
3. Log each request including status code (logrus)
4. Implement persistence with MySQL and Gorm (https://github.com/jinzhu/gorm) 
5. Use tool of your choice for DB migrations 
6. Implement save endpoint for Task object 
7. Implement update endpoint for Task object 
8. Implement get endpoint for Task object 
9. Implement delete endpoint for Task object (just update IsDeleted field)  
10. Use CORS (reply with header Access-Control-Allow-Origin: *) 
11. Add support for OPTION HTTP method for each endpoints  
12. Configure daemon over simple YAML config. Specify path as process flag for daemon. Required params: ListenAddress, DatabaseUri. 
13. Implement 3-rd party libs vendoring with tool of your choice. (godeps)
14. Put in comments below description of taken architectural decisions


###Task:
```
type Task struct {
    Id          int64
    Title       string
    Description string
    Priority    int
    CreatedAt   *time.Time
    UpdatedAt   *time.Time
    CompletedAt *time.Time
    IsDeleted   bool
    IsCompleted bool
}
```

# Design decisions

- error messages for simplicity were not declared as constants. This approach allows to quick find errors in code as they appear in logs.
- there are separate packages for models and handlers, but currently CRUD logic is located in handlers code for simplicity. Ideally we should have input and output structures in handlers package and Model structures in model package.
- authorization logic uses session storage in order to check CSRF tokens. Currently simple in-memory storage is used. But storage is passed as interface, so we can quickly substitute it with any other kind of storage (memcache, Aerospike, mysql, etc) we want.
- for simplicity SQLite datastorage is being used. Hopefully, golang database logic allows to change datastorage quickly. We can switch it with MySQL, for instance.
- logger is created in `main.go` in order to log messages that can appear outside of the application to the same logging channel.
- 
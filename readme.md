# users

- login (post)
- logout (post)

# reports

- reports (post)
- reports/:id (get)
- reports/:id (put)

- sample body for login
```
{
    "username": "admin",
    "password": "admin"
}
```

- sample body for logout
```
{
    "username": "admin"
}
```
- sample body for reports
```
{
    "title": "chilli",
    "type": "cooking",
    "description": "chilli is hot",
    "img": "www.storagepath.com/image.png"
}
```

- sample response API
{
    "status": "200",
    "message": "messages"
    "errors": [],
    "data": []
}

# Migration
- api/v1/migrate/on
- api/v1/migrate/off



# Run
- make setup-dev
- /bin/air
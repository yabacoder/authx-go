``About AuthX
This program allows user to register, login and search over a RESTful API. The search page is secured for ONLY logged in users to access. Each searches are also recorded. As per of setting up a profile, users also can upload their profile pictures to Amazon S3 platform

Strcuture
main.go
server
models
controllers
middleware

As you know, the main.go is the entrance into the app while the rest do the following jobs

.env // app.env holds the environment variables for the db and AWS services

Server
- server.go
This file has the startup, setup and also routing for each path.

Model
This folders has the table structures and also contains the database connection.

middleware
- veriftyToken: This handles the JWT verifcation to secure the pages.
- GetLoggedInUSer: This is a function used to retrieve logged in user data so that it can be used across the project.


controllers
- user.go : this handles the users operations from signing up, login, photo upload and logout.
- search.go : Handle searches and also records searches made by users. Only logged in users can search.



**Todo
- View own search history
- Delete Images
- Test
- API Documentation
- Delete and update photo
- Change password


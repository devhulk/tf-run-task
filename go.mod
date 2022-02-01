module github.com/devhulk/test-task

go 1.16

replace github.com/devhulk/runtask => ./runtask

replace github.com/devhulk/jwt => ./jwt

require github.com/golang-jwt/jwt v3.2.2+incompatible

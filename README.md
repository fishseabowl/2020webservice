Start server:
go run main.go

1. Support basic auth with admin:123456
2. Support Postman Get method: curl localhost:8080
3. Support Postman Post method: localhost:8080?name=blk&val=book
4. Support Postman Put method:  localhost:8080?name=hello&val=book
5. Support Postman Delete method: localhost:8080?name=hello


Build docker image and run it
docker build -t 2020webservice .
docker run --publish 8080:8080 --name test --rm 2020webservice

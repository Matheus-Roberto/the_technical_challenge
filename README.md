The Technical Challenge consists of creating an API with Golang using gRPC with stream pipes that exposes an upvote service endpoints.

Technical requirements:
- Keep the code in Github

API:
- The API must guarantee the typing of user inputs. If an input is expected as a string, it can only be received as a string.
- The structs used with your mongo model should support Marshal/Unmarshal with bson, json and struct
- The API should contain unit test of methods it uses

Extra:
- Deliver the whole solution running in some free cloud service


Postman - https://www.getpostman.com/collections/0dbf002cf291d44f83c8;

clone code - git clone git@github.com:Matheus-Roberto/the_technical_challenge_klever.git  

run server - go run cmd\server\main.go

run client - go run cmd\client\main.go

run test - go run test\test.go
  

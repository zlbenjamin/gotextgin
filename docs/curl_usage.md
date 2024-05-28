# Text
in Windows Command Prompt  

## POST /api/text/add
- curl -X POST http://localhost:40000/api/text/add -H "Content-Type: application/json" -d "{\"type\":\"text\"}"  
failed.  
{"code":400,"message":"Bad request","data":null}  
- curl -X POST http://localhost:40000/api/text/add -H "Content-Type: application/json" -d "{\"content\": \"text001\", \"type\":\"text\"}"  
success.  
{"code":200,"message":"OK","data":9}  

## GET /api/text/:id

- curl http://localhost:40000/api/text/99999999999999999999999999999  
{"code":400,"message":"Bad request","data":null}

- curl http://localhost:40000/api/text/-1  
{"code":400,"message":"Bad request \u003c 1","data":null}  
\u003c is <. TODO

- curl http://localhost:40000/api/text/2  
{"code":400,"message":"No data for id","data":null}

- curl http://localhost:40000/api/text/3  
{"code":200,"message":"OK","data":{}}

## DELETE /api/text/:id
- curl -X DELETE http://localhost:40000/api/text/2  
{"code":200,"message":"Do nothing","data":true}  

- curl -X DELETE http://localhost:40000/api/text/3  
{"code":200,"message":"OK","data":true}  

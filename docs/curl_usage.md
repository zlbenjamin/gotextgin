# Text
in Windows Command Prompt  

## POST /api/text/add
- curl -X POST http://localhost:40000/api/text/add -H "Content-Type: application/json" -d "{\"type\":\"text\"}"  
failed.  
{"code":400,"message":"Bad request","data":null}  
- curl -X POST http://localhost:40000/api/text/add -H "Content-Type: application/json" -d "{\"content\": \"text001\", \"type\":\"text\"}"  
success.  
{"code":200,"message":"OK","data":9}  



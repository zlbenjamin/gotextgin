# gotextgin
A text management program based on Golang, Gin, MySQL, etc..   

## run app
go run . -c config_file  
or  
go run . --config config_file  
Note, powered by viper.  

## go swagger
output directory: prjswagger  
CMD:  
$ swag init --output prjswagger  

After run the app, access the swagger docs:  
http://localhost:40000/swagger/index.html  


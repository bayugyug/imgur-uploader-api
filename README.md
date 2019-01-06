## imgur-uploader-api



### A simple golang script that bridge the Imgur REST API 
    
	( customize-image-upload-only as a mini-test )

	
### Registration Quickstart

	 
		- [x] Register an Application by providing required information @ https://api.imgur.com/oauth2/addclient
		
				- Application name: [your-app-name-here]
				
				- Authorization type: [select with callback url and use the https://www.getpostman.com/oauth2/callback]
				
				- Email: [your-email-address-here]
				
				- Captcha: [select i am not a robot]
				
				- Click Submit
				
		- [x] Save somewhere safe the Client-ID and Client-Secret
		
		- [x] Format your config paramter in JSON format:
		        
				Example: 
				{
					"client_id": "{YOUR_CLIENT_ID_FROM_IMGUR}", 
					"client_secret": "{YOUR_CLIENT_SECRET_FROM_IMGUR}"
				}
				
		- [x] Pass as parameter in running the docker.
		
				--credentials='{"client_id": "{YOUR_CLIENT_ID_FROM_IMGUR}", "client_secret": "{YOUR_CLIENT_SECRET_FROM_IMGUR}"}'


				
### Docker Binary


```sh

    sudo  sysctl -w net.ipv4.ip_forward=1

    sudo  docker run -p 7000-8000:7000-8000 -v `pwd`:`pwd` -w `pwd` -d --name imgur-uploader-api-latest  bayugyug/imgur-uploader-api:latest --port 7777   --credentials='{"client_id": "{YOUR_CLIENT_ID_FROM_IMGUR}", "client_secret": "{YOUR_CLIENT_SECRET_FROM_IMGUR}"}'



    curl -i -v 'http://127.0.0.1:7777/v1/api/images'

```


### Compile and Run-In-Command-Line


```sh

     git clone https://github.com/bayugyug/imgur-uploader-api.git && cd imgur-uploader-api

     git pull && make
		 

	 $ ./imgur-uploader-api

	 Version: 0.1.0-0

		 -config string
			use to set the config file parameter with info/credentials on Imgur
		 -credentials string
			use to set the info/credentials on Imgur
		 -port string
			use to set HTTP @ Port No. (default "7777")


	  Example:


			$ ./imgur-uploader-api --credentials='{"client_id": "{YOUR_CLIENT_ID_FROM_IMGUR}", "client_secret": "{YOUR_CLIENT_SECRET_FROM_IMGUR}"}'

			$ ./imgur-uploader-api --config=user.json


```

				
### Mini-How-To (List of End-Points)


##### 	- After running the docker binary, we need to supply the authorization-code by approving a permission on this app 
	  
##### 	- Paste the below URL on your browser   (can try using chrome)
			
		https://api.imgur.com/oauth2/authorize?access_type=offline&client_id={YOUR_CLIENT_ID_FROM_IMGUR}&response_type=code&state=state

##### 	- This will auto-redirect to the callback URL you've setup in the registration quickstart.
	   
	    https://app.getpostman.com/oauth2/callback?state=state&code={IMGUR_AUTH_CODE_IS_HERE}
		
##### 	- Pass this auth-code one-time to the api-bridge 

```sh
	    curl -v -X GET  'http://127.0.0.1:7777/v1/api/credentials/{IMGUR_AUTH_CODE_IS_HERE}'
		
		@output:
		{
		  "code": 202,
		  "message": "Accepted"
		}
```		
	
##### - Upload image URLs to the api-bridge
	
```sh
		curl -v  POST 'http://127.0.0.1:7777/v1/api/images/upload' -d '{
				"urls": [
					"https://farm3.staticflickr.com/2879/11234651086_681b3c2c00_b_d.jpg",
					"https://farm4.staticflickr.com/3790/11244125445_3c2f32cd83_k_d.jpg"
					]
				}'
				
		@output:
		{
		  "jobId": "1cfc8710-c366-4241-833b-3c5a988700cf-20190106-091212"
		}
```
	
##### 	- Get image URLs list
			
```sh		
		curl -X GET  'http://127.0.0.1:7777/v1/api/images'
		
		@output:
		{
		  "uploaded": [
			"https://i.imgur.com/8yc2oCz.jpg",
			"https://i.imgur.com/u6GIFQA.jpg"
		  ]
		}
```	
		
##### - Get image URLs list by jobId
	
	
```sh
		curl -v  GET 'http://127.0.0.1:7777/v1/api/images/upload/1cfc8710-c366-4241-833b-3c5a988700cf-20190106-091212'  
		
		@output:
		{
		  "id": "1cfc8710-c366-4241-833b-3c5a988700cf-20190106-091212",
		  "created": "2019-01-06T09:12:12Z",
		  "finished": "2019-01-06T09:12:28Z",
		  "status": "complete",
		  "uploaded": {
			"complete": [
			  "https://i.imgur.com/8yc2oCz.jpg",
			  "https://i.imgur.com/u6GIFQA.jpg"
			],
			"pending": null,
			"failed": null
		  }
		}
```
			
##### - Get image URLs list by jobId (Invalid JobId)
	
```sh
		curl -v  GET 'http://127.0.0.1:7777/v1/api/images/upload/f16fbca4-dae2-4c73-8304-df2966fa8831-20190106-nocontent' 
		
		@output:
		{
		  "code": 204,
		  "message": "No Content"
		}
```
	
	
	
### License

[MIT](https://bayugyug.mit-license.org/)

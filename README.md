## imgur-uploader-api

### A simple golang script that bridge the Imgur REST API

### Registration Quickstart

		
		- [x] [Register your application] (https://api.imgur.com/oauth2/addclient) 

		- [x] Register an Application by providing required information.
		
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
					"client_secret": "YOUR_CLIENT_SECRET_FROM_IMGUR"
				}
				
		- [x] Pass as parameter in running the docker.
		
				--credentials='{"client_id": "{YOUR_CLIENT_ID_FROM_IMGUR}", "client_secret": "YOUR_CLIENT_SECRET_FROM_IMGUR"}'
			
## Docker Binary


```sh

    sudo  sysctl -w net.ipv4.ip_forward=1

    sudo  docker run -p 7000-8000:7000-8000 -v `pwd`:`pwd` -w `pwd` -d --name imgur-uploader-api-latest  bayugyug/imgur-uploader-api:latest --port 7777  --config=user.json

    curl -i -v 'http://127.0.0.1:7777/v1/api/images'

```


### Run-As-A-Command-Line


```sh

$ ./imgur-uploader-api

 Version: 0.1.0-0

	 -config string
		use to set the config file parameter with info/credentials on Imgur
	 -credentials string
		use to set the info/credentials on Imgur
	 -port string
		use to set HTTP @ Port No. (default "7777")


  Example:


	$ ./imgur-uploader-api --credentials='{"client_id": "{YOUR_CLIENT_ID_FROM_IMGUR}", "client_secret": "YOUR_CLIENT_SECRET_FROM_IMGUR"}'

	$ ./imgur-uploader-api --config=user.json


```



### License

[MIT](https://bayugyug.mit-license.org/)

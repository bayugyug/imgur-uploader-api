## imgur-uploader-api

### A simple golang script that bridge the Imgur REST API


### Command


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



## Docker Binary

- [x] In order to  use it via CURL/WGET or Browser


```sh

    sudo  sysctl -w net.ipv4.ip_forward=1

    sudo  docker run -p 7000-8000:7000-8000 -v `pwd`:`pwd` -w `pwd` -d --name imgur-uploader-api-alpine  bayugyug/imgur-uploader-api:alpine --http --port 7778

    curl -i -v 'http://127.0.0.1:7777/v1/api/images'

```

### License

[MIT](https://bayugyug.mit-license.org/)

# converter

first install docker, and maybe make and docker-compose too

`docker build --tag converter .`  
`docker run --publish 5445:5445 converter`

or 

`make up` (docker-compose shortcut)

---

then POST a file in 'data' field to http://localhost:5445/api/v1/{inputFormatNumber}/{outputFormatNumber}/

- format numbers range from 1 to 3, corresponding with the project document specs
- please see the openapi public.yaml included for other arguments regarding delimiters when using format1


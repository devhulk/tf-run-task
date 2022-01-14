# Test Run Task

## Pre-reqs
- Install Certbot
- Install Golang
- Install Envoy
- Configure DNS Provider of choice to point at your IP (NoIP, Route53, etc)
- If running locally configure port forwarding
- Use certbot to generate certs for tls. Rename them and move to the envoy project dir

Certbot Install : 
```
brew install certbot
```

Cert Generation :
```
sudo certbot certonly --standalone
```
In the project envoy directory execute.Names changed to match envoy config file.
```
sudo cp /etc/letsencrypt/live/<your_domain_name>/fullchain.pem cert.pem

sudo cp /etc/letsencrypt/live/<your_domain_name>/privkey.pem private.pem
```

## Components
- Golang Http Server
- Envoy Proxy (HTTPS)

## Steps to Run

Start the go app in one terminal
```
go run main.go
```

Start the envoy Proxy Server in another terminal
```
envoy -c ./envoy/http.yaml
```

Run test script to see if your hitting go server. Fake data so will end in error but good for testing connectivity. In another terminal (optional)
```
curl -X POST -H "Content-Type: application/json" -d @./testData/init.json https://<your_domain>
```

TFCB Config: 
1. Configure Run Task in TFCB UI Settings. 
2. Then configure the run task in a workspace of your choice. 
3. Trigger the workspace and your run task will show up after the plan phase. If the message gets printed out it was successful. 

You can change if it will pass or fail in the runtask/runtask.go file. The passOrFail() function takes a string that lets you decide if the test passes or fails. In VI its on line 77.  

Your welcome to test using mine but no guarentees it will be up. If you all want I can throw up an instance that won't be deleted. Just let me know. My domain is https://devhulk.ddns.net/. You would use that when setting up the task in TFCB. No HMAC yet. 



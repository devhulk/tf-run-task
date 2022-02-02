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
In the project envoy directory execute the following.
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

My domain is https://devhulk.ddns.net/. You would use that when setting up the task in TFCB, I have it up most of the time but not checking religiously. No HMAC yet. 

## JWT

Added 02/01/22. Just wanted to test adding a JWT token to the workspace variables. Through out my testing I couldn't use the token given back from the run task, I had to use my own TFE/C token with the proper permissions. This is a good thing I think, if the token given back from the task could be used for org/workspace admin that would make run-tasks a vulnerability. 



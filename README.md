# go-coffeeshop

Deployment

### Setup enviorment

```
cd terraform
terraform init
terraform apply
```

### EC2

```
ssh -i "your_key.pem" ubuntu@public_ip

cd coffeeshop

sudo docker compose up -d --build

```
After all check your public ip there should be hosted frontend
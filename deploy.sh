
ssh user@john.cse.taylor.edu  # Use CSE Password

ssh user@cos143xl.cse.taylor.edu # Use special cos143 password

sudo su

docker ps
docker stop imageID

git pull
docker build --no-cache -f Dockerfile -t go143:1.0.0 .

# network host to allow go to connect to localhost redis
# port 3000 gets proxied to 8080 internally for https
docker run -d -p 3000:3000 \
--restart on-failure \
--network="host" \
-e REDIS_PASSWORD="REDIS_PASSWORD_HERE" \
-e S3_ACCESS_KEY="S3_ACCESS_KEY" \
-e S3_SECRET_KEY="S3_SECRET_KEY" \
go143:1.0.0 --port=3000 --logLevel=info

# Running Redis
sudo docker run \
-p 6379:6379 \
-v /home/jhibschm/redisData:/data \
--name redis \
--restart on-failure \
-d redis:6.0.9-alpine redis-server --appendonly yes  --requirepass "REDIS_PASSWORD_HERE"


# SSH bastion
ssh -L 6000:cos143xl.cse.taylor.edu:22 user@john.cse.taylor.edu
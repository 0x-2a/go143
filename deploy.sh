
ssh user@john.cse.taylor.edu  # Use CSE Password

ssh user@cos143xl.cse.taylor.edu # Use special cos143 password

sudo su

docker ps
docker stop imageID

git pull
docker build -f Dockerfile -t go143:1.0.0 .
docker run -d --restart on-failure -p 3000:8080 -e REDIS_PASSWORD="REDIS_PASSWORD_HERE" go143:1.0.0 --port=8080 --logLevel=info


# Running Redis
sudo docker run \
-p 6379:6379 \
-v /home/jhibschm/redisData:/data \
--name redis \
--restart on-failure \
-d redis:6.0.9-alpine redis-server --appendonly yes  --requirepass "REDIS_PASSWORD_HERE"
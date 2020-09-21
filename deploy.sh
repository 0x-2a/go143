
ssh user@john.cse.taylor.edu  # Use CSE Password

ssh cos143xl@john.cse.taylor.edu # Use special cos143 password

sudo su

docker ps
docker stop imageID

git pull
docker build -f Dockerfile -t go143:1.0.0 .
docker run -d --restart on-failure -p 3000:8080 go143:1.0.0 --port=8080 --logLevel=info

## RABBIT 
docker build -t my-image -f dockerfile .
  cd /home/bruno/Documents/murilo-pond-go-back/rabbit
  docker build -t rabbit-app -f dockerfile .

  docker run --rm rabbit-app

# latest RabbitMQ 4.x
docker run -it --rm --name rabbitmq -p 5672:5672 -p 15672:15672 rabbitmq:4-management
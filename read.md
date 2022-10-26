melakukan docker build 

# docker build -t todo-api:v1 .

menjalankan perintah docker run 

# docker run -d -p 3000:3000 --name todo-api todo-api:v1

melihat status kontent yang berjalan

# docker ps -a | grep todo-api
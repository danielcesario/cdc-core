# Define a imagem base
FROM golang:1.20-alpine

# Define o diretório de trabalho para a aplicação
WORKDIR /go/src/app

# Copia o código fonte da aplicação para o contêiner
COPY . .

# Altera o diretório de trabalho para o diretório do arquivo main
WORKDIR /go/src/app/cmd/api

# Compila o código fonte da aplicação para um binário
RUN go build -o app

# Define o comando padrão a ser executado quando o contêiner for iniciado
CMD [ "./app" ]

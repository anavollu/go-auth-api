# Go Auth API


## Aplicação de login em Golang com autenticação pelo AWS Cognito e deploy na AWS App Runner


## <img src="https://cdn-icons-png.flaticon.com/24/2666/2666505.png" width="20" /> Checklist

- [x] Criar um servidor no AWS AppRunner
- [x] Implementar um endpoint de autenticação utilizando o Amazon Cognito
- [x] Desenvolver a aplicação em Go, preferencialmente com o framework Gin
- [x] Criar um makefile com os comandos "make build" e "make deploy"
- [x] Integrar a aplicação com o AWS ECR para envio de imagens Docker
- [x] Vincular o AppRunner com o ECR para deploy automático
- [x] Disponibilizar a aplicação online em xxxx.aws.com [Link da aplicação](https://mckia3dckg.us-east-1.awsapprunner.com)
- [x] Criar os endpoints: /login, /logout e /home
- [x] Implementar a página de login com campo de usuário/senha e botão de login
- [x] Permitir acesso à página /home somente após autenticação bem-sucedida
- [x] Desenvolver a funcionalidade de logout no endpoint /logout
- [x] Utilizar JWT para autenticação, com chamadas ao Cognito durante o login


## <img src="https://cdn-icons-png.flaticon.com/24/2276/2276313.png" width="20" /> Tecnologias

Golang | Docker | AWS App Runner | AWS Cognito | AWS ECR


## <img src="https://cdn-icons-png.flaticon.com/24/5050/5050273.png" width="20" /> Deploy na AWS

[Link da aplicação](https://mckia3dckg.us-east-1.awsapprunner.com)


## <img src="https://cdn-icons-png.flaticon.com/128/10786/10786711.png" width="20" /> Diagramas

Primeiro acesso:

<img height="200px" src="https://github.com/anavollu/go-auth-api/assets/25857063/5c5e0bbc-1ac3-404e-80c5-97779d0308be" />

Outros acessos:

<img height="200px" src="https://github.com/anavollu/go-auth-api/assets/25857063/7231a66a-deda-4303-9d01-a363fcf403ee" />

Infra:

<img height="200px" src="https://github.com/anavollu/go-auth-api/assets/25857063/5141185a-3ba8-483a-bc8b-ccb61f358fb7" />

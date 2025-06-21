# имя проекта
PROJECTNAME=events

# генерация сервера с 0 (замена всех файлов генерации)
# custom templates is private (telegram:@UniMaxPod) 
appgen:
	rm -rf ./internal/generated
	mkdir -p ./internal/generated
	goswagger generate server \
	-f ./swagger-api/swagger.yml \
	-t ./internal/generated -C ./swagger-templates/default-server.yml \
	--template-dir ./swagger-templates/templates \
	--name $(PROJECTNAME)
	go mod tidy

# перегенерация апи (заменяются только файлы DoNotEdit)
apiregen:
	rm -rf ./internal/generated/models
	rm -rf ./internal/generated/restapi/operations
	rm ./internal/generated/restapi/embedded_spec.go
	rm ./internal/generated/restapi/server.go
	goswagger generate server \
	-f ./swagger-api/swagger.yml \
	-t ./internal/generated -C ./swagger-templates/default-server.yml \
	--template-dir ./swagger-templates/templates \
	--name $(PROJECTNAME)
	go mod tidy	


# запуск сервиса	
run:
	go run ./cmd/main.go
	
# форматирование файлов проекта	по fmt
fmt:
	go fmt ./...

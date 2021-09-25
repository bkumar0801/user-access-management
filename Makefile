db-init:
	psql -c 'CREATE DATABASE "user-access-management"' -U $(user)
db-init-kube:
	psql -c 'CREATE DATABASE "user-access-management"' -h $(host) -p $(port) -U $(user)
migrationup:
	migrate -path db/migrations -database "postgres://$(user):$(password)@$(host):$(port)/user-access-management?sslmode=disable" -verbose up
migrationdown:
	migrate -path db/migrations -database "postgres://$(user):$(password)@$(host):$(port)/user-access-management?sslmode=disable" -verbose down
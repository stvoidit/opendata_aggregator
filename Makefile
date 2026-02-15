include .env
export

# сделать дамп DDL создания базы, схем и таблиц
ddl.dump:
	pg_dump --host=${PGHOST} --port=${PGPORT} --username=${PGUSER} --dbname=opendata_aggregator -x --format=plain --schema-only --no-owner --create > opendata_aggregator.sql
# создание напустой БД базы, схем и таблиц без данных
ddl.restore:
	psql --host=${PGHOST} --port=${PGPORT} --username=${PGUSER} -d postgres -f opendata_aggregator.sql
# получить текущий конфиг postgresql
ddl.postgresql_config:
	ssh opendata '/bin/cat ~/postgresql.conf' > ./scripts/postgresql.conf

# загрузка на сервер исполняемого файла
upload-to-work.bin: build.backend
	scp src/backend/opendataaggregator ${DESTADDR}:~
	rm src/backend/opendataaggregator
# загрузка на сервер файлов сертификатов и ключей для ЕГР*
upload-to-work.certs:
	scp EGRIP.crt EGRIP.key EGRUL.crt EGRUL.key ${DESTADDR}:~/egr_certs

# старт бэка в режиме разработки
run-dev.backend:
	$(MAKE) -C src/backend serve
# старт профнта в режиме разработки
run-dev.frontend:
	cd src/frontend && pnpm dev

# вызов .sh скрипта для разворачивания .p12 сертификатов в нужный вид
unwraping_certificates_egr:
	/bin/bash -c scripts/unwraping_certificates_egr.sh

# вызов сборки бэка
build.backend:
	$(MAKE) -C src/backend build
# вызов сборки фронта
build.frontend:
	cd src/frontend && pnpm build

build.image:
	docker build --no-cache --build-arg=VERSION_TAG=$$(git describe --tags) --build-arg=COMMIT_SHA=$$(git rev-parse --short HEAD) --build-arg=GIT_BRANCH=$$(git branch --show-current) --file Dockerfile .

db.table_maintenance:
	/bin/bash -c scripts/tables_maintenance.sh

# профайлинг, чтобы не забыть команды
pprof.collect.profile:
	curl http://localhost:6060/debug/pprof/profile?seconds=70 > ./profile.out
pprof.collect.trace:
	curl http://localhost:6060/debug/pprof/trace?seconds=70 > ./trace.out
pprof.collect.heap:
	curl http://localhost:6060/debug/pprof/heap?seconds=70 > ./heap.out
pprof.pgo:
	curl -o default.pgo "http://localhost:6060/debug/pprof/profile?seconds=310"
pprof.analyze.profile:
	go tool pprof -web profile.out
pprof.analyze.heap:
	go tool pprof -web heap.out
pprof.analyze.trace:
	go tool trace -http=':8081' trace.out

git.tag.latest:
	git describe --tags --abbrev=0

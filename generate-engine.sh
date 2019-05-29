#!/usr/bin/env bash

mkdir -p ../EngineJS/tmp/src/{engine,domain/game/}
cp internal/app/domain/game/{action.go,game.go} ../EngineJS/tmp/src/domain/game/
cp internal/app/infrastructure/repository/memory/game/{engine.go,transport.go} ../EngineJS/tmp/src/engine/
sed -i '.bak' '1s/package game/package main/' ../EngineJS/tmp/src/engine/engine.go
sed -i '.bak' '1s/package game/package main/' ../EngineJS/tmp/src/engine/transport.go
sed -i '.bak' 's#github.com/go-park-mail-ru/2019_1_The-Recovery-Team/internal/app#engine/tmp/src#' ../EngineJS/tmp/src/engine/engine.go
sed -i '.bak' 's#github.com/go-park-mail-ru/2019_1_The-Recovery-Team/internal/app#engine/tmp/src#' ../EngineJS/tmp/src/engine/transport.go
echo 'func main() {
	js.Module.Get("exports").Set("initEngine", InitEngineJS)
}' >> ../EngineJS/tmp/src/engine/engine.go

sed -i '.bak' '/import*/a\
"github.com/gopherjs/gopherjs/js"\
' ../EngineJS/tmp/src/engine/engine.go

cd ../EngineJS/
go mod init engine
go mod tidy

gopherjs build -m tmp/src/engine/{engine.go,transport.go}
rm -r tmp
rm go.*
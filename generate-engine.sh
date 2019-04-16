#!/usr/bin/env bash

mkdir -p tmp/src/{engine,domain/game/}
cp internal/domain/game/{action.go,game.go} tmp/src/domain/game/
cp internal/infrastructure/repository/memory/game/{engine.go,transport.go} tmp/src/engine/
sed -i '.bak' '1s/package game/package main/' tmp/src/engine/engine.go
sed -i '.bak' '1s/package game/package main/' tmp/src/engine/transport.go
sed -i '.bak' 's#internal#tmp/src#' tmp/src/engine/engine.go
sed -i '.bak' 's#internal#tmp/src#' tmp/src/engine/transport.go
echo 'func main() {
	js.Module.Get("exports").Set("initEngine", InitEngineJS)
}' >> tmp/src/engine/engine.go

sed -i '.bak' '/import*/a\
"github.com/gopherjs/gopherjs/js"\
' tmp/src/engine/engine.go

mkdir engineJS
cd engineJS

gopherjs build -m ../tmp/src/engine/{engine.go,transport.go}
rm -r ../tmp
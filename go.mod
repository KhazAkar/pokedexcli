module github.com/KhazAkar/pokedexcli

go 1.24.0

replace (
	github.com/KhazAkar/pokedexcli/internal/pokeapi => ./internal/pokeapi
	github.com/KhazAkar/pokedexcli/internal/pokecache => ./internal/pokecache
)

require (
	github.com/KhazAkar/pokedexcli/internal/pokeapi v0.0.0-00010101000000-000000000000
	github.com/KhazAkar/pokedexcli/internal/pokecache v0.0.0-00010101000000-000000000000
)

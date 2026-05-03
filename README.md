# pokedexcli

A retro-style Pokemon command-line explorer written in Go. It uses the PokeAPI to let you browse areas, discover wild Pokemon, catch them with different ball types, inspect their stats, and build a persistent Pokedex across sessions.

---

## Motivation

Imagine you want to build a small terminal game that still feels interactive and stateful. You could make a very simple CLI that only fetches and prints Pokemon data on demand:

```go
pokemon := fetchPokemon("pikachu")
fmt.Println(pokemon.Name)
fmt.Println(pokemon.Types)
```

That works, but it feels more like an API viewer than a world you can play in. There is no sense of exploration, no memory between runs, and no reason to think strategically about encounters. 

`pokedexcli` solves this by turning PokeAPI data into a lightweight gameplay loop. Instead of only looking up Pokemon directly, you can explore areas, roll random encounters, choose different ball types when catching, and save your progress to disk:

```go
areas := listLocationAreas()
explore("kanto-route-2")
wild := randomEncounter()
fmt.Printf("A wild %s appeared!\n", wild.Name)

caught := catchPokemon(wild.Name, "great-ball")
if caught {
	savePokedex()
}
```

This keeps the project fun while also demonstrating practical Go concepts like API integration, caching, modular design, persistence, testing, and command-driven program flow.

--- 

## Installation

Inside your terminal:

```bash
git clone https://github.com/Exia2075/pokedexcli
cd pokedexcli
go run .
```

---

## Usage

Available commands include:

- `help`: Show all supported commands
- `map`: Show the next page of location areas
- `mapb`: Show the previous page of location areas
- `explore <location-area-name>`: Explore an area and refresh its encounter pool
- `encounter`: Roll a random wild Pokemon from the last explored area
- `balls`: Show available ball types and their catch bonuses
- `catch <pokemon-name> [ball-name]`: Try to catch a Pokemon
- `inspect <pokemon-name>`: View details for a caught Pokemon
- `pokedex`: Show all caught Pokemon
- `exit`: Quit the program

The CLI also stores your progress locally, so your Pokedex and last encounter pool remain available between sessions.

---

## Contributing

I love help! Contribute by forking the repo and opening pull requests.

Please make sure your changes pass existing tests, and add tests for new behavior when appropriate.

All pull requests should be submitted to the `main` branch.

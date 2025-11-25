# pokedexcli

**pokedexcli** is a Go command-line program that lets you explore a retro-style Pokémon world. You can move around a map, discover wild Pokémon, attempt to catch them, inspect their stats, and view your growing Pokédex all from the terminal.

---

## Motivation

The goal of pokedexcli is to provide a fun and lightweight project for practicing Go and CLI design. This project effectively teaches how to implement and use an efficient in-memory cache to reduce repeated API calls. Along the way, it serves as an introduction to:

- Building command-line interfaces in Go

- Working with APIs and parsing JSON

- Implementing TTL - based caching

- Structs, maps, and modular Go project structure. 

- Error handling and clean code practices.

--- 

## 🚀 Quick Start

Clone the repository:

```bash
git clone https://github.com/Exia2075/pokedexcli
cd pokedexcli
```

Run pokedexcli:

```bash
go run main.go
```

---

## 🚀 Usage

- map: Display the world map and your current location

- catch: Attempt to catch a discovered Pokémon

- inspect <name>: View stats of a Pokemon you've caught

- pokedex: Show your collected Pokémon

- exit: Quit the program

---

## 🧑‍💻 Author

Exia2075

GitHub: https://github.com/Exia2075/pokedex-cli

---

## 👏 Contributing

I would love your help! Contribute by forking the repo and opening pull requests. Please ensure that your code passes the existing tests and linting, and write tests to test your changes if applicable.

All pull requests should be submitted to the `main` branch.

# rotom

A CLI tool to display Pokémon sprites in your terminal.

![demo](demo.png)

### Why?

Almost all existing solutions are

1. not cross-platform (written in bash)
2. tedious to setup (written in python)

Rotom is a single statically-linked binary.

### Installation

#### From releases (recommended)
Download the latest release for your platform from the [releases page](https://github.com/shravanasati/rotom/releases).

#### Using Go
```bash
go install github.com/shravanasati/rotom@latest
```

### Usage

##### Download sprites (first time setup)
```bash
rotom download
```

##### Show random Pokémon
```bash
rotom
```

##### Show specific Pokémon
```bash
rotom pikachu
rotom 25
```

## Requirements

- Go 1.24+
- Terminal with image support (Kitty, iTerm2, WezTerm, foot, or any with Sixel)

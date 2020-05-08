# go-snake

Toy project: Go (golang) rewrite of my old Snake implementation in C

## Build requirements

* Go 1.14
* SDL (see [Installing the SDL packages](#installing-the-sdl-packages) below)

## Compiling

* Run `make`.

### Installing the SDL packages

Ensure you have pkgconfig and all SDL components installed. See full instructions in [go-sdl2][].

* Arch: `pacman -S sdl2{,_image,_mixer,_ttf,_gfx}`
* Debian/Ubuntu: `apt install libsdl2{,-image,-mixer,-ttf,-gfx}-dev`
* macOS with Homebrew: `brew install sdl2{,_gfx,_image,_mixer,_ttf} pkg-config`
* [Nix](https://nixos.org/nix/manual/) package manager: `nix-shell -p pkgconfig SDL2 SDL2_{gfx,image,mixer,ttf}`

## Development environment

### Git hooks

Run `make init` so that Git uses the hooks from the provided `.githooks` dir.

Currently those hooks ensure that `gofmt` style has been followed prior to committing Go source files.

### Vim

Source the `env.vim` file to setup `make` with automatic `gofmt`.

```vim
:source env.vim
```

[go-sdl2]: https://github.com/veandco/go-sdl2

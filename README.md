## Requirements

* Golang 1.10 (other versions were not tested)

## Compiling

* Prepare the SDL packages as described in the section below.
* Run the `./build` script.

After configuration, you can use `make` to build.

### Preparing the SDL packages

* Ensure you have pkgconfig and all SDL components installed
* Set the GOPATH (`source goenv/activate`)
* Run `go get -d` to download the SDL packages *without compiling them*
* Patch the sources with the appropriate pkgconfig information
* Run `go install` for all SDL packages

Using the [Nix](https://nixos.org/nix/manual/) package manager:

```
$ nix-shell -p pkgconfig SDL2 SDL2_{gfx,image,mixer,ttf}
(nix)$ source goenv/activate
(nix)$ go get -d github.com/veandco/go-sdl2/{sdl,mix,img,ttf,gfx}
(nix)$ patch -p1 -d goenv/src/github.com/veandco/go-sdl2/ < go-sdl2_pkgconfig.patch
(nix)$ go install -v github.com/veandco/go-sdl2/{sdl,mix,img,ttf,gfx}
```

On macOS, you can use the SDL packages from Homebrew:

```
$ brew install sdl2 sdl2_{gfx,image,mixer,ttf}
$ source goenv/activate
$ go get -d github.com/veandco/go-sdl2/{sdl,mix,img,ttf,gfx}
$ patch -p1 -d goenv/src/github.com/veandco/go-sdl2/ < go-sdl2_pkgconfig.patch
$ go install -v github.com/veandco/go-sdl2/{sdl,mix,img,ttf,gfx}
```

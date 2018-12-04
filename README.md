## Requirements

* Golang 1.10 (other versions were not tested)

## Compiling

* Prepare the SDL packages as described in the section below.
* Run the `./build` script.

### Preparing the SDL packages

Using [Nix](https://nixos.org/nix/manual/):

* Enter a Nix shell containing pkgconfig and all SDL components
* Set the GOPATH (`source goenv/activate`)
* Run `go get -d` to download the SDL packages *without compiling them*
* Patch the sources with the appropriate pkgconfig information
* Run `go install` for all SDL packages

```
$ nix-shell -p pkgconfig SDL2 SDL2_{gfx,image,mixer,ttf}
(nix)$ source goenv/activate
(nix)$ go get -d github.com/veandco/go-sdl2/{sdl,mix,img,ttf,gfx}
(nix)$ patch -p1 -d goenv/src/github.com/veandco/go-sdl2/ < go-sdl2_pkgconfig.patch
(nix)$ go install -v github.com/veandco/go-sdl2/{sdl,mix,img,ttf,gfx}
```

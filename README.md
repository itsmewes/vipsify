# Vipsify
Go Image manipulation service built on top of:
- [libvips](https://github.com/libvips/libvips)
- [Bimg](https://github.com/h2non/bimg)
- [Got](https://github.com/melbahja/got)

Vipsify is a hobby project intended to be able to manipulate images on the fly for frontend developement. There are so many different size devices out there, laptops to Monitors to every size phone out there that images can either be too large for the specified device or simply the image needs cropping. There are a few tools to manipulate images outside of just resizing and cropping, view the breakdown of the options below.

## Get Started
Install the above mentioned dependencies.

### libvips
According to their [Github Readme](https://github.com/libvips/libvips)
> libvips is a demand-driven, horizontally threaded image processing library. Compared to similar libraries, libvips runs quickly and uses little memory.

[Bimg](https://github.com/h2non/bimg) uses libvips because libvips
>requires a low memory footprint and it's typically 4x faster than using the quickest ImageMagick and GraphicsMagick settings or Go native image package, and in some cases it's even 8x faster processing JPEG images

Follow the instructions [here](https://libvips.github.io/libvips/install.html) to get set up.

### Bimg
[Bimg](https://github.com/h2non/bimg) is a:
>Small Go package for fast high-level image processing using libvips via C bindings, providing a simple programmatic API.

To install run:
`go get -u github.com/h2non/bimg`

### Got
[Got](https://github.com/melbahja/got) is used to be able to download images concurrently allowing for faster downloads.

To install run:
`go get github.com/melbahja/got/cmd/got`

Update the port in the `main.go` file to the one you would like to use before building (`go build .`)

## Usage
Vipsify uses query parameters to build up the instructions needed to processes the image. Once an image has been created a cached version of that image is saved for any subsequent visits.

## Options
`src` string The url to the image\
`w` int Width\
`h` int Height\
`t` string Type, eg png/jpg/webp (if the browser supports webp and a type is not set a webp image will be created)\
`c` int Compression, 0-8\
`o` string Options (comma seperated list, eg. `&o=flip,flop,crop`)
- `flip` string Flip the image horizontally
- `flop` string Flip the image vertically
- `crop` string Crop the image by the specified width and height (default is to not crop)
- `smart` string Crop to the interesting part of the image (libvips 8.5+)
- `bw` string Make the image black and white
- `fresh` string Regenerate the image each time. Don't use the cache.

## Example
https://domain.com?src=https://domain.com/image.jpg&w=100 \
https://domain.com?src=https://domain.com/image.jpg&w=100&h=100 \
https://domain.com?src=https://domain.com/image.jpg&w=100&h=100&o=crop,smart \
https://domain.com?src=https://domain.com/image.jpg&w=100&h=100&o=crop,flip&t=png
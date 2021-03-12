# [WIP]GoSpotDl
A Golang implemetation of spotdl (Python) 

I am learning go and thought implementing spotdl would be a good learning experience.

**Have a look at the [original](https://github.com/spotdl/spotify-downloader).**

## Usage
1. Get the package
    - `$ go get -u github.com/s1as3r/gospotdl`

2. Install the package
    - `$ go install github.com/s1as3r/gospotdl`

3. Enjoy  
    - `$ gospotdl $spotifylink` 

*Alternatively*:

1. Clone this repo
    - `$ git clone https://github.com/s1as3r/gospotdl`

2. Install the Dependencies
    - `$ cd GoSpotDl`
    - `$ go get -d ./...`

3. Build GoSpotDl
    - `$ go build`

4. Enjoy
    - `$ ./gospotdl $spotifylink`

## To-Do
- [x] Change to YTMusic Api.
- [ ] Tests.
- [ ] Fix ProgressBar when downloading multiple songs.
- [x] Parallel Downloads.
- [ ] Download best audio avaiable.
- [ ] Documentation.

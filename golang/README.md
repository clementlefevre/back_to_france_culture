# Download and sync France Culture podcasts with your mp3 player
Works with mp3 player using standard File systen (e.g poor's man player)


## Usage

either copy all urls into a text file and call :
`go run cmd/main.go --file mytxtfile_with_urls_list`

or simply paste the urls in the terminal like :
`go run cmd/main.go --url https://www.franceculture.fr/emissions/une-vie-une-oeuvre/toute-une-vie-du-vendredi-15-mai-2020,https://www.franceculture.fr/emissions/benito-mussolini-un-portrait/5-benito-mussolini-luomo-vuoto-lhomme-vide`

the scripts downloads the mp3 file linked in the url, and copy them in the local `podcasts` folder.

### Sync with external device
you can also sync the downloaded files with your local devices :
`go run cmd/main.go --file mytxtfile_with_urls_list --sync`
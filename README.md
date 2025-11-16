# Grabbler

Simple frontend to download a YouTube Video from a given URL. Metadata is fetched from a given Musicbrainz Recording ID.   
Metadata is not yet added to the downloaded file.

## Usage

You can either clone the repo, run backend and frontend manually or use docker.

```docker
services:
  backend:
    container_name: grabbler-backend
    image: ghcr.io/timhi/grabbler-backend:main
    volumes:
      - PATH_TO_DOWNLOAD:/app/downloads
    ports:
      - "3333:3333"

  frontend:
    container_name: grabbler-frontend
    image: ghcr.io/timhi/grabbler-frontend:main
    ports:
      - "3000:3000"
```

Replace `PATH_TO_DOWNLOAD` with your desired target download folder.

## TODO
[ ] Enrich downloaded audio file with metadata  
[ ] Add support for time codes 
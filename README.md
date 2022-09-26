# stLib
A platform that allows the viewing and management of 3d printing related projects and assets


## Problem
Many of us hoard a tone of STL's from public sites, patreon and so forth. Keeping a tidy library easy to access and search proved difficult (at least for me).

## Approach
The main goal is to achieve a self-hostable platform that work on top of the organization you already have, allowing you to have a nice overview of what you have collected so far.

## Demo website
https://demo.knoker.eu/projects

## Screenshots

### Home page with filter by project name and tags
![Home](/docs/Home.png)
### Project Image galery
![Images](/docs/Images.jpg)
### 3DView allows you to view multiple models of the project at once with zoom and pan controls
![3DView](/docs/3DView.png)
### Slice files with print details
![SliceDetails](/docs/SliceDetails.jpg)
### Edit page
![Edit](/docs/Edit.png)

## How to use
- Download the latest release for your platform
- Edit the config.toml
    - Change the library_path to the folder where you keep your stls
- Run the binary files
- Navigate to http://localhost:8000/projects on your browser
- Since you don't have initialized projects please toggle the initialized button and the projects should start to appear.
- Enter the project
    - Go to edit and save the projects to initialize it.
- Have fun.
- When something blows up please contact me on discord :)

## configuration

```toml	
port = 8000 # port to run the server on
library_path = "/library" # path to the stl library
max_render_workers = 5 # max number of workers to render the 3d model images in parallel shouldn't exceed the number of cpu cores
file_blacklist = [".potato",".example"] # list of files to ignore when searching for stl and assets files in the library_path
model_render_color = "#167DF0" # color to render the 3d model
model_background_color =  "#FFFFFF"  # color to render the 3d model background
thingiverse_token = "your_thingiverse_token" # thingiverse token to allow the import of thingiverse projects

```

## docker-compose

```yaml
---
version: "3.6"
services:
  stlib:
    image: ghcr.io/eduardooliveira/stlib:main
    container_name: stlib
    volumes:
      - ./library:/library
      - ./config.toml:/app/config.toml
    ports:
      - 8000:8000
    environment:
      - "THINGIVERSE_TOKEN=" # needed for the thingiverse download feature
      #- "PORT=8000"
      #- "LIBRARY_PATH=./library"
      #- "MAX_RENDER_WORKERS=5"
      #- "MODEL_RENDER_COLOR=#167DF0"
      #- "MODEL_BACKGROUND_COLOR=#FFFFFF"
      #- "LOG_PATH=./log" # If you wish to log to a file
    
    restart: unless-stopped
```

## Discussion
![Discord Shield](https://discordapp.com/api/guilds/1013417395777450034/widget.png?style=shield)

Join discord if you need any support https://discord.gg/SqxKE3Ve4Z


## TODO Features

- [ ] Facelift
- [ ] Add klipper integration to send jobs to the printer
- [ ] Add support for other slice formats
- [ ] Add support for other slicers
- [ ] Add detail to the slice view
- [ ] Allow project creation and file upload
- [ ] Allow slice upload directly from the slicer
- [x] Discover other files in the filesystem
- [x] Improve the 3DView
- [x] Allow model upload
- [x] Show slice settings (print time, speed, material)
- [x] Discover slice files in the filesystem
- [x] Allow default project image definition
- [x] Discover images in the library
- [x] Find projects in the filesystem
- [x] Render a static image of the models
- [x] Allow edition of the projects
- [x] Allow search by tags
- [x] Allow search by project name
- [x] Allow 3DView of the models
- [x] Allow multiple models in the 3DView

## TODO Technical

- [x] Launch a demo instance
- [x] Build a docker image
- [x] Add a bounding box to the 3d objects to center the camera
- [x] Allow models to be set as image
- [x] Remove backend template rendering
- [x] Refactor endpoint files to match API
- [x] Cleanup static files
- [x] Host the Vue compiled files from the Echo Server
- [x] Add frontend configuration
- [x] Add a build & release system
- [x] Add backend configuration
 

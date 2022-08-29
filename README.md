# stLib
A platform that allows the viewing and managment of 3d printing related projects and assets


## Problem
Many of us hoard a tone of STL's from public sites, pathreons and so forth. Keeping a tidy library easy to access and search proved dificult (at least for me).

## Approch
The main goal is to achive a self-hostable platform that work on top of the organization you already have, allowing you to have a nice overview of what you have collected so far.

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
- Download the latest relase for your platform
- Edit the config.toml
    - Change the library_path to the folder where you keep your stls
- Run the binary files
- Navigate to http://localhost:8000/projects on your browser
- Since you don't have initialized projects please toggle the initialized button and the projects should start to appear.
- Enter the project
    - Go to edit and save the projects to initialize it.
- Have fun.
- When something blows up please contact me on discod :)

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
    ports:
      - 8000:8000
    restart: unless-stopped

```

## Discussion
![Discord Shield](https://discordapp.com/api/guilds/1013417395777450034/widget.png?style=shield)

Join discord if you have any support https://discord.gg/SqxKE3Ve4Z


## TODO Features

- [ ] Facelift
- [ ] Add klipper integration to send jobs to the printer
- [ ] Add support for other slice formats
- [ ] Add support for other slicers
- [ ] Add detail to the slice view
- [ ] Discover other files in the filesystem
- [ ] Allow project creation and file upload
- [ ] Allow slice upload directly from the slicer
- [ ] Improve the 3DView
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

## TODO Techical

- [ ] Launch a demo instance
- [ ] Build a docker image
- [ ] Add a bounding box to the 3d objets to center the camera
- [x] Allow models to be set as image
- [x] Remove backend template rendering
- [x] Refactor endpoint files to match API
- [x] Cleanup static files
- [x] Host the Vue compiled files from the Echo Server
- [x] Add frontend configuration
- [x] Add a build & release system
- [x] Add backend configuration
 

# Containers CI/CD pipeline

**Goal:**
Intuitive & easy orchestration of docker containers on single remote VPS.

In order to host a new webiste it's enough to create a module under /compose directory.
A module is composed of those necessary configuration files in order to run with docker compose (usually .env + docker-compose.yml is enough).
User should be able to easily modify modules for each website both manually, from Makefile configuration and (in the future) from web UI dashboard like a mini VPS management system.

The project is composed of a core set of *core* modules that handle orchestration of docker containers through docker *templates* (.env file + docker-compose.yml + extra sh files(optional))
that can be extended as a web UI dashboard like a mini VPS management system.

Eveything inside root dir is aligned with VPS file system through GitHub Actions workflow.

## Host new container *manually*

1. Create a new *template* under `/templates` (.env + docker-compose.yml) if doesn't already exists one.
2. Copy the files from the template to a new *module* under `/compose` (.env + docker-compose.yml) and populate the `.env` file with actual configuration values.
3. Run with docker compose under `/compose/[module]`

    ```bash
    cd /compose/[module]
    docker compose -p [module] up -d
    ```

## Host new container from Makefile

Makefile is used as glue cli interface for the Go application logic.

1. Move under root `/`
2. Run:

    ```bash
    make dock CONTAINER="webserver" TEMPLATE="traefik"
    ```

    It should run the Go app inherent to *docking* process (create configuration file and run).

## Host new container from web UI

...

## First time setup on VPS

### Allow GitHub

1. Add ssh public key

    ```bash
    ssh-keygen -t ed25519 -C "vps-github-deploy"
    ```

2. Add the public key (`~/.ssh/id_ed25519.pub`) to GitHub as a Deploy Key:
    - Go to your GitHub repo → Settings → Deploy Keys

    - Click "Add Deploy Key" and paste the contents of the public key

    - Check "Allow write access" if you plan to push too (optional)

3. Test

    ```bash
    ssh -T git@github.com
    ```

## Setup on local environment

### Clone and generate dedicated ssh key

1. Clone from GitHub

   Under `/go-docker-manager`

    ```bash
    git clone https://github.com/FrancescoCorbosiero/go-docker-manager.git
    ```

    or under current dir

    ```bash
    git clone https://github.com/FrancescoCorbosiero/go-docker-manager.git .
    ```

2. Generate SSH key

    ```bash
    ssh-keygen -t rsa -b 4096 -f deploy_github_actions
    ```

3. Add the generated public ssh key on a new line inside `~/.ssh/authorized_keys`

## Core modules

1. `/templates` - Docker compose templates

    Ready to run with `docker compose -p`

   ```docker
    templates/
    ├── webserver/
    │   ├── .env
    │   └── docker-compose.yml
    ├── website/
    │   ├── .env
    │   └── docker-compose.yml
    └── ...
    ```

2. `/operations` - Executable script

    ```bash
    chmod +x operations/[script-name].sh
    ```

3. `/compose` - contains actual .env configuration for running containers

4. Run with Make

    ```txt
    make: Shows the help message.

    make list: Shows running containers

    make logs CONTAINER=site2: Tails the logs for site2.

    make down CONTAINER=site1: Stops and removes containers for site1.

    make restart CONTAINER=site1: Stops and then starts site1.

    make dock CONTAINER=site1 TEMPLATE=template: create module (.env + compose file) under /compose if doesn't exists and run
    ```

## Utils
## Backup

### Traefik
```
cd /var/lib/docker/volumes/website_traefik-certificates/_data && zip -r /root/traefik-certificates.zip ./*
```

### Wordpress website
#### Wordpress content
```
cd /var/lib/docker/volumes/[volume_name]/_data && zip -r /root/[project_name]-data.zip ./*
```

#### Database
```
cd /var/lib/docker/volumes/[volume_name]/_data && zip -r /root/[project_name]-data.zip ./*
```

#### Permissions + cleanup of script
```
root@vmi2548180:~/docker# chmod +x ./scripts/traefik-up.sh
root@vmi2548180:~/docker# ./scripts/traefik-up.sh
-bash: ./scripts/traefik-up.sh: cannot execute: required file not found
root@vmi2548180:~/docker# sed -i 's/\r//' ./scripts/traefik-up.sh
root@vmi2548180:~/docker# ./scripts/traefik-up.sh
```
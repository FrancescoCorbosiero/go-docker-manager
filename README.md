# Tool for Docker Compose

## ⚠️: This is an experimental project in the prototyping phase — future development is highly unlikely due to poor Proof of Concept.

**Goal:**
Intuitive & easy management of `docker compose` operations.

When it comes to simple and straightforward orchestration of your docker containers, meaning:

- run, stop, delete containers
- volume backup management
- monitoring of the running containers

you don't usually need complex solutions like Kubernetes or Ansible,
`docker compose` itlself should be enough with his *declarative* and *deterministic* `docker-compose.yml` file both with the interpolation of `environment variables` defined in `.env` file.

Unfortunately, `docker compose` was not designed for scalable multi-instance reuse out-of-the-box.

While Compose supports some environment variable interpolation via .env or export, it doesn’t let you easily:

- Dynamically reuse the same file for N independent instances
- Launch N isolated containers from the same Compose file without collisions

Compose makes this possible through -p flag:

`docker compose -p [namespace] up -d`

However configuration and organization of the file system can be messy.

This application aims to let the user create and maintain his own file system organization and setup of the docker compose *easily*.

Users can create a docker compose `template` module, containing the dockerized version of the application, and run as many instances he wants without worrying about configuration files conflicts.

Another important design key of the application is to be friendly for on-point manual operations, which means that this setup is easily modifiable *manually* (directly on file system) or through CLI commands (*Makefile*).

## Host new container *manually* with Docker Compose

1. Create a new *template* module under `/templates` (.env + docker-compose.yml) if doesn't already exists one.
2. Copy the files from the template to a new *module* under `/compose` (.env + docker-compose.yml) and populate the `.env` file with actual configuration values.
3. Run with docker compose under `/compose/[module]`

    ```bash
    cd /compose/[module]
    docker compose -p [module] up -d
    ```

This manual process is perfectly fine but you would usually prefer using Makefile instead.

## Host new container with Makefile

Makefile is used as glue cli interface for the Go application logic.

1. Move under root `/`
2. Chose a template and the container name. Run:

    ```bash
    make dock CONTAINER="webserver" TEMPLATE="traefik"
    ```

3. Wait until logs are written and check for status:

    ```bash
    make list
    ```

    You should see your new container running and healthy.

## Host new container from web UI

**WIP**

## First time setup on VPS

### Allow GitHub Deploy

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
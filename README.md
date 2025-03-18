Before starting, set up the environment:
```sh
cp .envrc-template .envrc

# Edit variables in .envrc:
vi .envrc

# Source them:
source .envrc

# Alternatively, if you use direnv:
direnv allow .
```

Then, to launch the app locally:
```sh
# Start the database container and apply migrations
docker compose up -d db
scripts/migrate.sh up

# Start the server
scripts/run.sh
```

Note: [run.sh](./scripts/run.sh) requires that the `tailwindcss` binary be available in PATH (see [tailwind.sh](./scripts/tailwind.sh)).  
An alternative method that doesn't require an explicit binary download is to launch the app in Docker:
```sh
# Start the database and apply migrations as before
# Then, launch the app container
docker compose up -d hawloom-local
```

Either way, the main page should now be available at http://127.0.0.1:22440/

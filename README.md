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
# Start the database container
docker compose up -d db

# Apply migrations
scripts/migrate.sh up

# Create a key
go build ./cmd/keyrot
POSTGRES_HOST=127.0.0.1 ./keyrot -c config/dev.json

# Start the server
scripts/run.sh

# Navigate to the main page
xdg-open http://127.0.0.1:22440/
```

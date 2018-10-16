# Go DFW Meetup 
### Testing A Restful API

This codebase is an example RESTFUL API for posting cars to a postgres database.

#### PRE REQS:

  - Docker installed (https://www.docker.com/)
  - Tested on MAC OS X

#### Make Commands (ran from root of project):
 - `Make build-linux` (will build linux binary to be served in docker)
 - `Make test` (docker compose build/docker-compose up/test) <- Spins up postgres and app
 - `Make run` (docker-compose up project)
 - `Make test-func` (Runs functional tests)

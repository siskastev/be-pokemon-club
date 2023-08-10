# be-pokemon-club
Technical Test Backend Engineer about Pokemon Fight Club Rocket Team

## Project Scopes
* Get all the data and details of the Pokemon that will fight with the integration from the API [pokeapi.co/](https://pokeapi.co/)pokedex

* Create battles Pokemon for role ops and team

* Get all the data and details of the battles with criteria regular role ops and team members will be able to see data for battles they've personally created. The role bos will able to access all data.

* Get rankings of all Pokémon battles for the boss role, showcasing the highest-scoring Pokémon that deserve promotion.

## Setup
1. Clone the repository
2. Navigate to the project directory:
3. Rename or copy .env.example to .env

## Running the app with Docker
```bash
$ docker-compose up
```

## API Documentation
The API documentation for the API Pokemons Fight Club is provided in the `docs/api-doc.json` file. 
This file allows you to quickly test the API endpoints when the Docker container is running. 
You can use various API testing tools or utilities to interact with the API and view the responses for each endpoint.

## Authentication
The API uses token authentication. To authenticate requests, include the API token in the request headers:
`Authorization: Bearer {api_token}`.
1. To log in as the bos, use the following credentials:
* `email`: bos@gmail.com
* `password`: inibos
2. To log in as the ops, use the following credentials:
* `email`: ops@gmail.com
* `password`: iniops
3. To log in as the team, use the following credentials:
* `email`: hallo@gmail.com `password`: "hallo"
* `email`: care@gmail.com `password`: "care"

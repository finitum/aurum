
# Aurum Web Client

This is the webclient for Aurum. With it you can log in and change your user information. 
As an admin you can block users. 

## Running

Change the config file (`ts/Config.ts`) as you like it.

Install dependencies:
```bash
yarn install
```

Run the dev server:
```bash
yarn run
```

## Linting and Type checking
As parcel doesn't automatically run type checking you need to  do this manually,
we made a command to make this easier, just run
```bash
yarn check
```
to run `eslint` and `tsc`


## Tests
We use `ts-jest` for our tests. To run the tests make sure all dev dependencies are installed and run:
```bash
yarn test
```

## Building
To build a docker image of the client for production simply run
```
docker build . -t aurum-core
```
This will build a docker container containing everything and using `nginx` as the web server.


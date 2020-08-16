# Cloud_Asset_Uploader

# Project description:
This is a simple RESTful web service that
can be used to upload assets to AWS S3 on the behalf of the person who started it. To do that it issues presigned upload and download URLs that expire after a given period specified in seconds by the API caller.

# Architecture:
This is a simple service written in GoLang that uses MongoDB or DataStax Astra to store the applications state.
I tried to come up with a stateless architecture, but could not due to the fact that a unique asset ID has
to be generated (and in practice linked to) the asset before it is actually uploaded. The API is 100% covered by unit tests and the data layer is covered with integration tests.

# Prerequisites and Gotchas
One has to setup AWS credentials file as described in [the official documentation](https://docs.aws.amazon.com/cli/latest/userguide/cli-configure-files.html). The file has to be in "~/.aws/" directory and be named credentials, so the AWS SDK and the provided docker-compose can pick it up correctly. The provided docker-compose will mount your AWS credentials file into the container
and the application in the container will use them to make calls to the AWS API through the SDK. 
To start the service with MongoDB locally run: <br/>
`docker-compose up --build` . Tested with docker-compose version 1.25.5. (--build is not required after the first time if you have not made any changes to the code). <br/>The docker-compose file contains all the required ENV vars for the service. Those are ENV vars for configurable connection for MongoDB, the port of the server and some AWS config. The application
provides default values for the server and AWS vars, but the MongoDB ones have to be set in order to connect successfully. The provided compose file has them all set up, so it can be used as is.
To stop the service run: <br/>
`docker-compose down`
### Warning: DO NOT upload the resulting docker image to any public docker repository as you risk leaking your AWS credentials and allowing other people to use S3 on your behalf.

## Useful commands:
The service provides a Makefile that contains several useful one line commands.<br/>
To run the unit tests only:
`make unit`<br/>
To run the whole test suite (this requires an actual DB for the integration tests, so make sure the application is properly started first):<br/>
`make integration` <br/>
This will also show a test coverage report in your default browser.<br/>
`make build` will build the Docker image for the service.<br/>
`make format` will format the code.<br/>
`make linter` will run a linter for Go and show if there are any violations in the console.

# Astra instructions
In order to start the service with Astra one has to follow a few steps:
1. Provision Astra first - a 2 click process explained on the DataStax site. <br/>
2. Get your Astra instance credentials as described [in the docs](https://docs.astra.datastax.com/docs/obtaining-database-credentials) <br/>
3. Unzip them and put them in a folder, then substitute the sample value for a path in the Makefile "astra" target with the absolute path to the credentials. <br/>
4. Setup env file and provide correct path for it in the Makefile "astra" target. The following env vars should be provided: <br/>
    ASTRA_HOST=REAL_VALUE_FOR_YOUR_INSTANCE <br/>
    ASTRA_PORT=REAL_PORT_FOR_YOUR_INSTANCE <br/>
    ASTRA_USERNAME=REAL_USERNAME_FOR_THE_DB <br/>
    ASTRA_PASSWORD=REAL_PASSWORD_FOR_THE_DB <br/>
    ASTRA_KEYSPACE=REAL_DB_KEYSPACE  <br/>
    SERVER_PORT=:8090   <br/>
    AWS_REGION=eu-central-1     <br/>
    AWS_BUCKET_NAME=cloud-assets-uploader <br/>
Make sure the port is the one from the cqlshrc file. <br/>
5. Now the service can be started with `make astra`.


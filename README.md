# Voltage AutoUnlock Webhook Service

[![Deploy to DO](https://www.deploytodo.com/do-btn-white-ghost.svg)](https://cloud.digitalocean.com/apps/new?repo=https://github.com/w3irdrobot/voltageautounlock/tree/master&refcode=b61ec14cb278)
[![Deploy to Heroku](https://www.herokucdn.com/deploy/button.svg)](https://heroku.com/deploy?template=https://github.com/w3irdrobot/voltageautounlock/tree/master)

To facilitate an easy setup for [auto-unlocking on Voltage](https://docs.voltage.cloud/lightning-nodes/webhooks#example-automatic-unlock), this project allows for a quick setup of a service to run as the server receiving the webhook to then unlock the node.

## Deployment

There are a few provided deployment methods to make setting up easy. Each method requires some values to be supplied to the service so it can run correctly. These values and their locations in the Voltage database are discussed [in the Voltage documentation](https://docs.voltage.cloud/lightning-nodes/webhooks#example-automatic-unlock).

> Note: After all deployments, ensure to update the Voltage dashboard with the location of the webhook.

- [Fly.io](#flyio)
- [Railway](#railway)
- [Docker](#docker)
- [Source](#source)

### Fly.io

Deploying to [fly.io](https://fly.io/) can be easily done using the [`flyctl`](https://fly.io/docs/flyctl/installing/) CLI.

```shell
# view regions here: https://fly.io/docs/reference/regions/
flyctl launch \
    --generate-name \
    --image w3irdrobot/voltageautounlock:1.1.1 \
    --region ord \
    --no-deploy

flyctl secrets set \
    VOLTAGE_NODE_API=<insert Node API URL> \
    VOLTAGE_WEBHOOK_SECRET=<insert webhook secret> \
    VOLTAGE_WALLET_PASSWORD=<insert wallet password>

flyctl deploy

# input this hostname into voltage
echo "https://$(flyctl info --host)"
```

### Railway

Deploying to [Railway](https://railway.app) can be easily done using the [`railway`](https://docs.railway.app/develop/cli#install) CLI.

```shell
# initialize the project
railway init
# deploy the service
railway up
```

This will fail because the environment variables are missing. Go to [the Railway dashboard](https://railway.app/dashboard), navigate to the service, and add the below variables using the raw editor.

```
VOLTAGE_NODE_API=<insert Node API URL>
VOLTAGE_WEBHOOK_SECRET=<insert webhook secret>
VOLTAGE_WALLET_PASSWORD=<insert wallet password>
```

### Docker

A Docker image is already built and pushed to Docker Hub. To use this image, just run it!

```shell
docker run -it \
    --name autounlock \
    --env VOLTAGE_NODE_API=<insert Node API URL> \
    --env VOLTAGE_WEBHOOK_SECRET=<insert webhook secret> \
    --env VOLTAGE_WALLET_PASSWORD=<insert wallet password> \
    w3irdrobot/voltageautounlock
```

### Source

To deploy from source, download [the latest binary](https://github.com/w3irdrobot/voltageautounlock/releases/latest) for the system this service will be run on. Set the necessary environment variables and run the binary.

```shell
curl -Lo voltageunlock.tar.gz https://github.com/w3irdrobot/voltageautounlock/releases/download/1.1.1/voltageautounlock_1.1.1_linux_amd64.tar.gz
tar -xzf voltageunlock.tar.gz
export VOLTAGE_NODE_API=<insert Node API URL>
export VOLTAGE_WEBHOOK_SECRET=<insert webhook secret>
export VOLTAGE_WALLET_PASSWORD=<insert wallet password>
./server
```

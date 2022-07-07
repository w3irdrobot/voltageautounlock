# Voltage AutoUnlock Webhook Service

[![Deploy to DO](https://www.deploytodo.com/do-btn-blue-ghost.svg)](https://cloud.digitalocean.com/apps/new?repo=https://github.com/w3irdrobot/voltageautounlock/tree/master&refcode=0b3c9298b62d)
[![Deploy to Heroku](https://www.herokucdn.com/deploy/button.svg)](https://heroku.com/deploy?template=https://github.com/w3irdrobot/voltageautounlock/tree/master)

To facilitate an easy setup for [auto-unlocking on Voltage](https://docs.voltage.cloud/lightning-nodes/webhooks#example-automatic-unlock), this project allows for a quick setup of a service to run as the server receiving the webhook to then unlock the node.

## Deployment

There are a few provided deployment methods to make setting up easy. Each method requires some values to be supplied to the service so it can run correctly. These values and their locations in the Voltage database are discussed [in the Voltage documentation](https://docs.voltage.cloud/lightning-nodes/webhooks#example-automatic-unlock).

- [Docker](#docker)
- [DigitalOcean](#digitalocean)
- [Heroku](#heroku)
- [Source](#source)

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

### DigitalOcean

To deploy to the DigitalOcean App Platform, click the button.

[![Deploy to DO](https://www.deploytodo.com/do-btn-blue-ghost.svg)](https://cloud.digitalocean.com/apps/new?repo=https://github.com/w3irdrobot/voltageautounlock/tree/master&refcode=0b3c9298b62d)

### Heroku

To deploy to the Heroku App Platform, click the button.

[![Deploy to Heroku](https://www.herokucdn.com/deploy/button.svg)](https://heroku.com/deploy?template=https://github.com/w3irdrobot/voltageautounlock/tree/master)

### Source

To deploy from source, download [the latest binary](https://github.com/w3irdrobot/voltageautounlock/releases/latest) for the system this service will be run on. Set the necessary environment variables and run the binary.

```shell
curl -Lo voltageunlock.tar.gz https://github.com/w3irdrobot/voltageautounlock/releases/download/1.0.0/voltageautounlock_1.0.0_linux_amd64.tar.gz
tar -xzf voltageunlock.tar.gz
export VOLTAGE_NODE_API=<insert Node API URL>
export VOLTAGE_WEBHOOK_SECRET=<insert webhook secret>
export VOLTAGE_WALLET_PASSWORD=<insert wallet password>
./server
```

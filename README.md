# Motivation

To be able to post in a Discord channel based on the overnight temperature, within the constraints of all the free tiers of AWS.

Many parts of this project are for fun & education.

# Overview

## [Go](https://go.dev/)
A concise, statically typed language. If it compiles, it works?

## [Discord Webhook](https://support.discord.com/hc/en-us/articles/228383668-Intro-to-Webhooks)
As this script does not require any particularly complex Discord functionality we use a simple webhook to post to our specified Discord channel.

## Weather API ([Open-Meteo](https://open-meteo.com/en))
A fantastic and free weather API with high temporal resolution, allowing us to look up temperatures over the next day.

## [Docker multistage build](https://docs.docker.com/develop/develop-images/multistage-build/)
Utilising docker multistage builds to build our Go application but base the final image on [alpine](https://hub.docker.com/_/alpine/) to greatly reduce the size.

We do not use AWS provided Lambda images, as they are (as of writing) not up to date with the latest Go version. We instead install RIE manually on the last stage.

## [GitHub Actions](https://github.com/features/actions)
Triggered on push to master, this will build our image, push to ECR and point Lambda at our new image.

We could stop here and use GitHub Actions to run our script, as it's a once a day script, we could simply run `go run .` instead of building an image and running on Lambda. However, the taken approach gives a better idea on how to deploy serverless applications in the real world.

## [ECR](https://aws.amazon.com/ecr/)
Push our image to ECR to be used within Lambda.

We have a lifecycle policy to remove all old images to keep storage costs low.

## [Lambda](https://aws.amazon.com/lambda/)
As we only want to run this script once a day there is no need to have a VM running constantly. Instead we use Lambda to point to our image in ECR so we only run the script when we need to.

## [CloudWatch](https://aws.amazon.com/cloudwatch/)
Used as our task scheduler to trigger Lambda once a day.

Cloudwatch also captures log output of each run.

# Development

To run this repo locally, you will need [Docker](https://www.docker.com/) installed.

Once installed, you will need to build and run the image.

As this repo is deployed on Lambda you can only interact with the image via the [RIE layer](https://github.com/aws/aws-lambda-runtime-interface-emulator).

### Build
`docker build -t captain-cold .`

### Run
`docker run -p 9000:8080 --env-file .env captain-cold`

### Sending requests once running
`curl -XPOST "http://localhost:9000/2015-03-31/functions/function/invocations"`

*The project at current does not have discrete lambda function handlers, therefore you do not need a POST body.*
# Motivation

To be able to post in a Discord channel based on the overnight temperature.

Many aspects of this project are for learning.

# Overview

## [Go](https://go.dev/)
A concise, statically typed language. If it compiles, it works?

## [Discord Webhook](https://support.discord.com/hc/en-us/articles/228383668-Intro-to-Webhooks)
As this script does not require any particularly complex Discord functionality we use a simple webhook to post to our specified Discord channel.

## Weather API ([Open-Meteo](https://open-meteo.com/en))
A fantastic, free, open source, no sign-up weather API with high temporal resolution, allowing us to look up temperatures over the next day.

## [Docker multistage build](https://docs.docker.com/develop/develop-images/multistage-build/)
Utilising docker multistage builds to build our Go application but base the final image on [alpine](https://hub.docker.com/_/alpine/) to greatly reduce the size.

## [GitHub Actions](https://github.com/features/actions)
Triggered on push to master, this will build our image, push to ECR and point Lamda at our new image.

We could stop here and use GitHub Actions to run our script, as it's a once a day script, we could simply run `go run .` instead of building an image and running Lamda. However, the taken approach gives a better idea on how to deploy serverless applications correctly.

## [ECR](https://aws.amazon.com/ecr/)
Push our image to ECR to be used within Lamda.

We have a lifecycle policy to remove all old images to keep storage costs down.

## [Lamda](https://aws.amazon.com/lambda/)
As we only want to run this script once a day there is no need to have a VM running constantly. Instead we use Lamda to point to our image in ECR so we only run the script when we need to.

## [CloudWatch](https://aws.amazon.com/cloudwatch/)
Used as our task scheduler to trigger Lamda once a day and provides us log output of each run.
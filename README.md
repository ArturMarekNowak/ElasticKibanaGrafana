# ElasticKibanaGrafana

[![Build](https://github.com/ArturMarekNowak/ElasticKibanaGrafana/actions/workflows/workflow.yml/badge.svg)](https://github.com/ArturMarekNowak/ElasticKibanaGrafana/actions/workflows/workflow.yml/badge.svg) [![Trivy and dockler](https://github.com/ArturMarekNowak/ElasticKibanaGrafana/actions/workflows/image-scan.yml/badge.svg)](https://github.com/ArturMarekNowak/ElasticKibanaGrafana/actions/workflows/image-scan.yml/badge.svg) [![CodeFactor](https://www.codefactor.io/repository/github/arturmareknowak/elastickibanagrafana/badge)](https://www.codefactor.io/repository/github/arturmareknowak/elastickibanagrafana)

This project is a PoC of monitoring and metrics gathering in go 

## Table of contents
* [General info](#general-info)
* [Technologies](#technologies)
* [Setup](#setup)
* [Status](#status)
* [Inspiration](#inspiration)

## General info

After docker compose is run the project should be ready to go. Project consists of six docker containers:

<p align="center"><img src="./docs/network.drawio.png"/>
<p align="center">Pic.1 Visualization of the project</p>

Go app implements only one endpoint - _GET /helloWorld_. Each should be seen in kibana and be visualized in grafana. In kibana after the index is setup one can use filter with _path_ selector to extract API call from log:

<p align="center"><img src="./docs/kibanaFilter.png"/>
<p align="center">Pic.2 Usage of filter</p>

<p align="center"><img src="./docs/kibanaLog.png"/>
<p align="center">Pic.3 Filtered result</p>

In grafana with usage of simple promQL query such as `rate(http_requests_count[1m])`, data can be visualized:

<p align="center"><img src="./docs/grafana.png"/>
<p align="center">Pic.4 Number of calls to the specific endpoint</p>


## Technologies
* Go 1.22
* Docker

## Setup
1. Run docker compose in src folder: `docker-compose up`

## Status
Project is: _finished_

## Inspiration
Its good to have such small PoC 

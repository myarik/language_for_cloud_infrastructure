# Language for cloud infrastructure 

This repository contains the [presentation](https://www.slideshare.net/YaroslavMuravskiy/programming-language-for-the-cloud-infrastructure-238708126) and sample codes written in Go, Python, Node 
which demonstrates how to resolve some tasks using concurrency.  

## Downloader

### Description 

Create an asynchronous program for downloading and storing content.

### Workflow task schema

![schema downloader](doc/assets/downloader_schema.png)

### Execution

![gif downloader](doc/assets/downloader.gif)

### Benchmark

![benchmark downloader](doc/assets/benchmark_downloader.png)

## Slow consumer

### Description 

Create an asynchronous program for downloading and storing content. 
Consider that only three consumers can run simultaneously and 
a producer shouldn't start a download process if all consumers are busy.  

### Workflow task schema

![schema slow_consumer](doc/assets/slow_consumer.png)

### Execution

![gif slow_consumer](doc/assets/slow_consumer.gif)

### Benchmark

![benchmark slow_consumer](doc/assets/benchmark_slow_consumer.png)

## First response

### Description 

Create an asynchronous program for getting the first response from the multiple replicas.  

### Workflow task schema

![schema first_response](doc/assets/first_response.png)

### Execution

![gif slow_consumer](doc/assets/first_response.gif)

### Benchmark

![benchmark first_response](doc/assets/benchmark_first_response.png)


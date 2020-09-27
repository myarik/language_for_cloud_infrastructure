SHELL := /bin/bash # Use bash syntax
.PHONY: docker_build_downloader docker_build_slow_consumer run_downloader run_slow_consumer

# define standard colors
BLACK        := $(shell tput -Txterm setaf 0)
RED          := $(shell tput -Txterm setaf 1)
GREEN        := $(shell tput -Txterm setaf 2)
YELLOW       := $(shell tput -Txterm setaf 3)
LIGHTPURPLE  := $(shell tput -Txterm setaf 4)
PURPLE       := $(shell tput -Txterm setaf 5)
BLUE         := $(shell tput -Txterm setaf 6)
WHITE        := $(shell tput -Txterm setaf 7)

RESET := $(shell tput -Txterm sgr0)

.PHONY: help
## help: prints this help message
help:
	@echo "Usage:"
	@sed -n 's/^##//p' ${MAKEFILE_LIST} | column -t -s ':' |  sed -e 's/^/ /'

# Build
## docker_build_downloader: build the downloader application
docker_build_downloader:
	@echo "${GREEN}Building downloader application${RESET}"
	${MAKE} -C node-demo docker_build_downloader
	${MAKE} -C python-demo docker_build_downloader
	${MAKE} -C go-demo docker_build_downloader


## docker_build_slow_consumer: build the slow consumer application
docker_build_slow_consumer:
	@echo -n "${GREEN}Building slow consumer application${RESET}"
	${MAKE} -C node-demo docker_build_slow_consumer
	${MAKE} -C python-demo docker_build_slow_consumer
	${MAKE} -C go-demo docker_build_slow_consumer

## docker_build_first_response: build the first response application
docker_build_first_response:
	@echo -n "${GREEN}Building first response application${RESET}"
	${MAKE} -C node-demo docker_build_first_response
	${MAKE} -C python-demo docker_build_first_response
	${MAKE} -C go-demo docker_build_first_response

# Run
## run_downloader: runs the downloader applications
run_downloader:
	@echo "${YELLOW}Python${RESET}"
	${MAKE} -C python-demo run_downloader
	@echo "${GREEN}Node${RESET}"
	${MAKE} -C node-demo run_downloader
	@echo "${LIGHTPURPLE}Golang${RESET}"
	${MAKE} -C go-demo run_downloader

## run_slow_consumer: runs the slow consumer applications
run_slow_consumer:
	@echo "${YELLOW}Python${RESET}"
	${MAKE} -C python-demo run_slow_consumer
	@echo "${GREEN}Node${RESET}"
	${MAKE} -C node-demo run_slow_consumer
	@echo "${LIGHTPURPLE}Golang${RESET}"
	${MAKE} -C go-demo run_slow_consumer

## run_first_response: runs the first response applications
run_first_response:
	@echo "${YELLOW}Python${RESET}"
	${MAKE} -C python-demo run_first_response
	@echo "${GREEN}Node${RESET}"
	${MAKE} -C node-demo run_first_response
	@echo "${LIGHTPURPLE}Golang${RESET}"
	${MAKE} -C go-demo run_first_response

## benchmark_downloader: runs the downloader benchmark
benchmark_downloader:
	@echo "${YELLOW}Python${RESET}"
	number=1 ; while [[ $$number -le 5 ]] ; do \
        ${MAKE} -C python-demo run_downloader; \
        ((number = number + 1)) ; \
    done
	@echo "${GREEN}Node${RESET}"
	number=1 ; while [[ $$number -le 5 ]] ; do \
        ${MAKE} -C node-demo run_downloader; \
        ((number = number + 1)) ; \
    done
	@echo "${LIGHTPURPLE}Golang${RESET}"
	number=1 ; while [[ $$number -le 5 ]] ; do \
        ${MAKE} -C go-demo run_downloader; \
        ((number = number + 1)) ; \
    done

## benchmark_slow_consumer: runs the slow_consumer benchmark
benchmark_slow_consumer:
	@echo "${YELLOW}Python${RESET}"
	number=1 ; while [[ $$number -le 5 ]] ; do \
        ${MAKE} -C python-demo run_slow_consumer; \
        ((number = number + 1)) ; \
    done
	@echo "${GREEN}Node${RESET}"
	number=1 ; while [[ $$number -le 5 ]] ; do \
        ${MAKE} -C node-demo run_slow_consumer; \
        ((number = number + 1)) ; \
    done
	@echo "${LIGHTPURPLE}Golang${RESET}"
	number=1 ; while [[ $$number -le 5 ]] ; do \
        ${MAKE} -C go-demo run_slow_consumer; \
        ((number = number + 1)) ; \
    done

## benchmark_first_response: runs the first_response benchmark
benchmark_first_response:
	@echo "${YELLOW}Python${RESET}"
	number=1 ; while [[ $$number -le 5 ]] ; do \
        ${MAKE} -C python-demo run_first_response; \
        ((number = number + 1)) ; \
    done
	@echo "${GREEN}Node${RESET}"
	number=1 ; while [[ $$number -le 5 ]] ; do \
        ${MAKE} -C node-demo run_first_response; \
        ((number = number + 1)) ; \
    done
	@echo "${LIGHTPURPLE}Golang${RESET}"
	number=1 ; while [[ $$number -le 5 ]] ; do \
        ${MAKE} -C go-demo run_first_response; \
        ((number = number + 1)) ; \
    done


## show_size: displays the image sizes
show_size:
	docker images --format "{{.Tag}}\t{{.Size}}\t{{.Repository}}" | grep myarik/demo | sort -k 2 -h


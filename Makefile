ROOT_DIR    = $(shell pwd)
NAMESPACE   = "default"
DEPLOY_NAME = "template-single"
DOCKER_NAME = "hjm-certcheck"
DOCKER_PREFIX = "hjmcloud"
# Install/Update to the latest CLI tool.
.PHONY: cli
cli:
	@set -e; \
	wget -O gf https://github.com/gogf/gf/releases/latest/download/gf_$(shell go env GOOS)_$(shell go env GOARCH) && \
	chmod +x gf && \
	./gf install -y && \
	rm ./gf


# Check and install CLI tool.
.PHONY: cli.install
cli.install:
	@set -e; \
	gf -v > /dev/null 2>&1 || if [[ "$?" -ne "0" ]]; then \
  		echo "GoFame CLI is not installed, start proceeding auto installation..."; \
		make cli; \
	fi;


# Generate Go files for DAO/DO/Entity.
.PHONY: dao
dao: cli.install
	@gf gen dao

# Generate Go files for Service.
.PHONY: service
service: cli.install
	@gf gen service

# Build image, deploy image and yaml to current kubectl environment and make port forward to local machine.
.PHONY: start
start:
	@set -e; \
	make image; \
	make deploy; \
	make port;

# Build docker image.
.PHONY: image
image: cli.install
	$(eval _TAG  = $(shell git log -1 --format="%cd.%h" --date=format:"%Y%m%d%H%M%S"))
ifneq (, $(shell git status --porcelain 2>/dev/null))
	$(eval _TAG  = $(_TAG).dirty)
endif
	$(eval _TAG  = $(if ${TAG},  ${TAG}, dev))
	$(eval _PUSH = $(if ${PUSH}, ${PUSH}, ))
	@gf docker $(PUSH) -b "-a amd64 -s linux -p temp" -tn $(DOCKER_NAME):${_TAG} -tp $(DOCKER_PREFIX);


# Build docker image and automatically push to docker repo.
.PHONY: image.push
image.push:
	@make image PUSH=-p;


# Deploy image and yaml to current kubectl environment.
.PHONY: deploy
deploy:
	$(eval _TAG = $(if ${TAG},  ${TAG}, develop))

	@set -e; \
	mkdir -p $(ROOT_DIR)/temp/kustomize;\
	cd $(ROOT_DIR)/manifest/deploy/kustomize/overlays/${_TAG};\
	kustomize build > $(ROOT_DIR)/temp/kustomize.yaml;\
	kubectl   apply -f $(ROOT_DIR)/temp/kustomize.yaml; \
	kubectl   patch -n $(NAMESPACE) deployment/$(DEPLOY_NAME) -p "{\"spec\":{\"template\":{\"metadata\":{\"labels\":{\"date\":\"$(shell date +%s)\"}}}}}";

# Build binary files and publish.
.PHONY: bin.publish
bin.publish:
	@set -e; \
	rm -rf $(ROOT_DIR)/temp; \
	gf build -n $(DOCKER_NAME) -a 386,amd64,arm,arm64 -s linux,darwin,windows -p temp;\
	cd $(ROOT_DIR)/temp;\
	git init;\
	git add -A;\
	git commit -m "deploy | $(shell date +'%Y-%m-%d %H:%M:%S')";\
	echo "$(shell date +'%Y-%m-%d %H:%M:%S')|开始到gitee部署";\
	git push -f git@gitee.com:$(DOCKER_PREFIX)/$(DOCKER_NAME).git master:release;\
	echo "$(shell date +'%Y-%m-%d %H:%M:%S')|完成到gitee部署";\
	open https://gitee.com/$(DOCKER_PREFIX)/$(DOCKER_NAME)/pages;



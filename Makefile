# Set these to the desired values
ARTIFACT_ID=k8s-blueprint-operator
VERSION=0.1.0
## Image URL to use all building/pushing image targets
IMAGE_DEV=${K3CES_REGISTRY_URL_PREFIX}/${ARTIFACT_ID}:${VERSION}
IMAGE=cloudogu/${ARTIFACT_ID}:${VERSION}
GOTAG=1.21
MAKEFILES_VERSION=9.0.1
LINT_VERSION=v1.55.2
STAGE?=production


ADDITIONAL_CLEAN=dist-clean

include build/make/variables.mk
include build/make/self-update.mk
include build/make/dependencies-gomod.mk
include build/make/build.mk
include build/make/test-common.mk
include build/make/test-unit.mk
include build/make/static-analysis.mk
include build/make/clean.mk
include build/make/digital-signature.mk
include build/make/mocks.mk

PRE_COMPILE=generate
K8S_RESOURCE_TEMP_FOLDER ?= $(TARGET_DIR)
K8S_PRE_GENERATE_TARGETS=k8s-create-temporary-resource template-dev-only-image-pull-policy

K8S_CRD_COMPONENT_SOURCE=$(WORKDIR)/k8s/helm-crd/templates/k8s.cloudogu.com_blueprints.yaml
CRD_SRC_GO=$(WORKDIR)/pkg/api/v1/blueprint_types.go
K8S_COPY_CRD_TARGET_DIR=$(WORKDIR)/pkg/api/v1

include build/make/k8s-controller.mk

.PHONY: build-boot
build-boot: helm-apply kill-operator-pod ## Builds a new version of the operator and deploys it into the K8s-EcoSystem.

##@ Deployment

.PHONY: helm-values-update-image-version
helm-values-update-image-version: $(BINARY_YQ)
	@echo "Updating the image version in source values.yaml to ${VERSION}..."
	@$(BINARY_YQ) -i e ".manager.image.tag = \"${VERSION}\"" ${K8S_COMPONENT_SOURCE_VALUES}

.PHONY: helm-values-replace-image-repo
helm-values-replace-image-repo: $(BINARY_YQ)
	@if [[ ${STAGE} == "development" ]]; then \
      		echo "Setting dev image repo in target values.yaml!" ;\
    		$(BINARY_YQ) -i e ".manager.image.repository=\"${IMAGE_DEV}\"" "${K8S_COMPONENT_TARGET_VALUES}" ;\
    	fi

.PHONY: template-stage
template-stage: $(BINARY_YQ)
	@if [[ ${STAGE} == "development" ]]; then \
  		echo "Setting STAGE env in deployment to ${STAGE}!" ;\
		$(BINARY_YQ) -i e ".manager.env.stage=\"${STAGE}\"" ${K8S_COMPONENT_TARGET_VALUES} ;\
	fi

.PHONY: template-log-level
template-log-level: ${BINARY_YQ}
	@if [[ "${STAGE}" == "development" ]]; then \
      echo "Setting LOG_LEVEL env in deployment to ${LOG_LEVEL}!" ; \
      $(BINARY_YQ) -i e ".manager.env.logLevel=\"${LOG_LEVEL}\"" "${K8S_COMPONENT_TARGET_VALUES}" ; \
    fi

.PHONY: template-image-pull-policy
template-image-pull-policy: $(BINARY_YQ)
	@if [[ "${STAGE}" == "development" ]]; then \
          echo "Setting pull policy to always!" ; \
          $(BINARY_YQ) -i e ".manager.imagePullPolicy=\"Always\"" "${K8S_COMPONENT_TARGET_VALUES}" ; \
    fi

.PHONY: kill-operator-pod
kill-operator-pod:
	@echo "Restarting k8s-blueprint-operator!"
	@kubectl -n ${NAMESPACE} delete pods -l 'app.kubernetes.io/name=k8s-blueprint-operator'

##@ Debug

.PHONY: print-debug-info
print-debug-info: ## Generates indo and the list of environment variables required to start the operator in debug mode.
	@echo "The target generates a list of env variables required to start the operator in debug mode. These can be pasted directly into the 'go build' run configuration in IntelliJ to run and debug the operator on-demand."
	@echo "STAGE=$(STAGE);LOG_LEVEL=$(LOG_LEVEL);KUBECONFIG=$(KUBECONFIG);NAMESPACE=$(NAMESPACE);DOGU_REGISTRY_ENDPOINT=$(DOGU_REGISTRY_ENDPOINT);DOGU_REGISTRY_USERNAME=$(DOGU_REGISTRY_USERNAME);DOGU_REGISTRY_PASSWORD=$(DOGU_REGISTRY_PASSWORD);DOCKER_REGISTRY={\"auths\":{\"$(docker_registry_server)\":{\"username\":\"$(docker_registry_username)\",\"password\":\"$(docker_registry_password)\",\"email\":\"ignore@me.com\",\"auth\":\"ignoreMe\"}}}"

############
# DEFAULTS #
############

GOOS                 ?= $(shell go env GOOS)
GOARCH               ?= $(shell go env GOARCH)

#########
# TOOLS #
#########

TOOLS_DIR                          ?= $(PWD)/.tools
CONTROLLER_GEN                     := $(TOOLS_DIR)/controller-gen
CONTROLLER_GEN_VERSION             ?= v0.20.0
REGISTER_GEN                       ?= $(TOOLS_DIR)/register-gen
DEEPCOPY_GEN                       ?= $(TOOLS_DIR)/deepcopy-gen
CODE_GEN_VERSION                   ?= v0.35.0
GEN_CRD_API_REFERENCE_DOCS         ?= $(TOOLS_DIR)/gen-crd-api-reference-docs
GEN_CRD_API_REFERENCE_DOCS_VERSION ?= latest
GENREF                             ?= $(TOOLS_DIR)/genref
GENREF_VERSION                     ?= master
TOOLS                              := $(CONTROLLER_GEN) $(REGISTER_GEN) $(DEEPCOPY_GEN) $(GEN_CRD_API_REFERENCE_DOCS) $(GENREF)
ifeq ($(GOOS), darwin)
SED                                := gsed
else
SED                                := sed
endif

$(CONTROLLER_GEN):
	@echo Install controller-gen... >&2
	@cd ./hack/controller-gen && GOBIN=$(TOOLS_DIR) go install -buildvcs=false

$(REGISTER_GEN):
	@echo Install register-gen... >&2
	@GOBIN=$(TOOLS_DIR) go install k8s.io/code-generator/cmd/register-gen@$(CODE_GEN_VERSION)

$(DEEPCOPY_GEN):
	@echo Install deepcopy-gen... >&2
	@GOBIN=$(TOOLS_DIR) go install k8s.io/code-generator/cmd/deepcopy-gen@$(CODE_GEN_VERSION)

$(GEN_CRD_API_REFERENCE_DOCS):
	@echo Install gen-crd-api-reference-docs... >&2
	@GOBIN=$(TOOLS_DIR) go install github.com/ahmetb/gen-crd-api-reference-docs@$(GEN_CRD_API_REFERENCE_DOCS_VERSION)

$(GENREF):
	@echo Install genref... >&2
	@GOBIN=$(TOOLS_DIR) go install github.com/kubernetes-sigs/reference-docs/genref@$(GENREF_VERSION)

.PHONY: install-tools
install-tools: ## Install tools
install-tools: $(TOOLS)

.PHONY: clean-tools
clean-tools: ## Remove installed tools
	@echo Clean tools... >&2
	@rm -rf $(TOOLS_DIR)

###########
# CODEGEN #
###########

CONFIG_PATH := ./config
CRDS_PATH := $(CONFIG_PATH)/crds

.PHONY: register-gen
register-gen: ## Generate API types registrations
register-gen: $(REGISTER_GEN)
	@echo Generate registration... >&2
	@$(REGISTER_GEN) --go-header-file=./hack/boilerplate.go.txt --output-file zz_generated.register.go ./api/...

.PHONY: deepcopy-gen
deepcopy-gen: ## Generate API deep copy functions
deepcopy-gen: $(DEEPCOPY_GEN)
	@echo Generate deep copy functions... >&2
	@$(DEEPCOPY_GEN) --go-header-file ./hack/boilerplate.go.txt --output-file zz_generated.deepcopy.go ./api/...

.PHONY: controller-gen
controller-gen: ## Generate policies CRDs
controller-gen: $(CONTROLLER_GEN)
	@echo Generate policies crds... >&2
	@rm -rf $(CRDS_PATH) && mkdir -p $(CRDS_PATH)
	@$(CONTROLLER_GEN) \
		paths=./api/policies.kyverno.io/... \
		crd:crdVersions=v1,ignoreUnexportedFields=true,generateEmbeddedObjectMeta=false \
		output:dir=$(CRDS_PATH)
	@cat $(CRDS_PATH)/*.yaml > $(CONFIG_PATH)/crds.yaml

.PHONY: api-docs
api-docs: ## Generate API docs
api-docs: $(GEN_CRD_API_REFERENCE_DOCS)
api-docs: $(GENREF)
	@echo Generate api docs... >&2
	@rm -rf docs/user/crd && mkdir -p docs/user/crd
	@$(GEN_CRD_API_REFERENCE_DOCS) \
		-api-dir github.com/kyverno/api/api \
		-config docs/user/config.json \
		-template-dir docs/user/template \
		-out-file docs/user/crd/index.html
	@cd ./docs/user && $(GENREF) \
		-c config-api.yaml \
		-o crd \
		-f html

.PHONY: helm-chart
helm-chart: ## Generate helm CRDs
helm-chart: controller-gen
	@echo Generate helm crds... >&2
	@rm -rf charts/kyverno-api/templates/crds && mkdir -p charts/kyverno-api/templates/crds
	@cp $(CRDS_PATH)/*.yaml charts/kyverno-api/templates/crds/
	@$(SED) -i '/^  annotations:/a \ \ \ \ {{- include "kyverno-api.annotations" . | nindent 4 }}' charts/kyverno-api/templates/crds/*
	@$(SED) -i '/^  annotations:/i \ \ labels:' charts/kyverno-api/templates/crds/*
	@$(SED) -i '/^  labels:/a \ \ \ \ {{- include "kyverno-api.labels" . | nindent 4 }}' charts/kyverno-api/templates/crds/*
	@$(SED) -i '/controller-gen.kubebuilder.io/d' charts/kyverno-api/templates/crds/*

.PHONY: codegen
codegen: ## Generate all generated code
codegen: register-gen
codegen: deepcopy-gen
codegen: controller-gen
codegen: api-docs
codegen: helm-chart

##################
# VERIFY CODEGEN #
##################

.PHONY: codegen
verify-codegen: ## Verify all generated code and docs are up to date
verify-codegen: codegen
	@echo Checking git diff... >&2
	@echo 'If this test fails, it is because the git diff is non-empty after running "make codegen".' >&2
	@echo 'To correct this, locally run "make codegen", commit the changes, and re-run tests.' >&2
	@git diff --exit-code

########
# HELP #
########

.PHONY: help
help: ## Shows the available commands
	@grep -E '^[a-zA-Z_-]+:.*?## .*$$' $(MAKEFILE_LIST) | awk 'BEGIN {FS = ":.*?## "}; {printf "\033[36m%-40s\033[0m %s\n", $$1, $$2}'

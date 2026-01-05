#########
# TOOLS #
#########

TOOLS_DIR                          ?= $(PWD)/.tools
CONTROLLER_GEN                     := $(TOOLS_DIR)/controller-gen
CONTROLLER_GEN_VERSION             ?= v0.18.0
REGISTER_GEN                       ?= $(TOOLS_DIR)/register-gen
DEEPCOPY_GEN                       ?= $(TOOLS_DIR)/deepcopy-gen
CODE_GEN_VERSION                   ?= v0.34.1
TOOLS                              := $(CONTROLLER_GEN) $(REGISTER_GEN) $(DEEPCOPY_GEN)

$(CONTROLLER_GEN):
	@echo Install controller-gen... >&2
	@cd ./hack/controller-gen && GOBIN=$(TOOLS_DIR) go install -buildvcs=false

$(REGISTER_GEN):
	@echo Install register-gen... >&2
	@GOBIN=$(TOOLS_DIR) go install k8s.io/code-generator/cmd/register-gen@$(CODE_GEN_VERSION)

$(DEEPCOPY_GEN):
	@echo Install deepcopy-gen... >&2
	@GOBIN=$(TOOLS_DIR) go install k8s.io/code-generator/cmd/deepcopy-gen@$(CODE_GEN_VERSION)

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

CRDS_PATH := ./config/crds

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
		paths=./api/policies.kyverno.io/v1alpha1/... \
		paths=./api/policies.kyverno.io/v1beta1/... \
		crd:crdVersions=v1,ignoreUnexportedFields=true,generateEmbeddedObjectMeta=false \
		output:dir=$(CRDS_PATH)

.PHONY: codegen
codegen: ## Generate all generated code
codegen: register-gen
codegen: deepcopy-gen
codegen: controller-gen

##################
# VERIFY CODEGEN #
##################

.PHONY: codegen
verify-codegen: ## Verify all generated code and docs are up to date
verify-codegen: codegen-all
	@echo Checking git diff... >&2
	@echo 'If this test fails, it is because the git diff is non-empty after running "make codegen-fix-tests".' >&2
	@echo 'To correct this, locally run "make codegen-fix-tests", commit the changes, and re-run tests.' >&2
	@git diff --exit-code

########
# HELP #
########

.PHONY: help
help: ## Shows the available commands
	@grep -E '^[a-zA-Z_-]+:.*?## .*$$' $(MAKEFILE_LIST) | awk 'BEGIN {FS = ":.*?## "}; {printf "\033[36m%-40s\033[0m %s\n", $$1, $$2}'

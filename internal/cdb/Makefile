STAGE_DIR := staging

CLANG_BASE_URL := https://raw.githubusercontent.com/llvm/llvm-project/refs/heads/release/22.x
CLANG_DIR := $(STAGE_DIR)/clang
CLANG_FILES := \
	clang/include/clang/Options/Options.td \
	llvm/include/llvm/Option/OptParser.td
CLANG_TBLGEN := /opt/homebrew/opt/llvm/bin/clang-tblgen

GCC_BASE_URL := https://raw.githubusercontent.com/gcc-mirror/gcc/refs/heads/releases/gcc-16
GCC_DIR := $(STAGE_DIR)/gcc
GCC_FILES := \
	gcc/c-family/c.opt \
	gcc/common.opt \
	gcc/params.opt \
	gcc/analyzer/analyzer.opt

.PHONY: stage stage-clang stage-gcc gen gen-gcc gen-clang

stage: stage-clang stage-gcc

stage-clang: $(addprefix $(CLANG_DIR)/,$(CLANG_FILES)) $(CLANG_DIR)/dump.json
$(CLANG_DIR)/%:
	mkdir -p $(dir $@)
	curl -fsSL $(CLANG_BASE_URL)/$* -o $@
$(CLANG_DIR)/dump.json:
	cd $(CLANG_DIR) && \
		$(CLANG_TBLGEN) -dump-json \
			-I clang/include/clang/Options \
			-I llvm/include \
			clang/include/clang/Options/Options.td \
			-o dump.json

stage-gcc: $(addprefix $(GCC_DIR)/,$(GCC_FILES))
$(GCC_DIR)/%:
	mkdir -p $(dir $@)
	curl -fsSL $(GCC_BASE_URL)/$* -o $@

gen: stage ./config/gcc.go ./config/clang.go

./config/gcc.go: $(GCC_DIR)
	go run ./cmd/gccconfiggen -o $@ $(GCC_DIR)

./config/clang.go: $(CLANG_DIR)/dump.json
	go run ./cmd/clangconfiggen -o $@ $(CLANG_DIR)/dump.json
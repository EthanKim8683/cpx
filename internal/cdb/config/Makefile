GCC_HOME := ../../../third_party/gcc

x: gcc.awk
	awk -f $(GCC_HOME)/gcc/opt-gather.awk \
		$(GCC_HOME)/gcc/c-family/c.opt \
		$(GCC_HOME)/gcc/common.opt \
		$(GCC_HOME)/gcc/params.opt \
		$(GCC_HOME)/gcc/analyzer/analyzer.opt \
		| awk -f gcc.awk \
			-f $(GCC_HOME)/gcc/opt-read.awk \
			-f $(GCC_HOME)/gcc/opt-functions.awk \
			> gcc.yaml
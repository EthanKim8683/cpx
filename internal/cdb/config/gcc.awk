BEGIN {
	VISIBILITY_DRIVER = "driver"
	VISIBILITY_C = "language/c"
	VISIBILITY_CXX = "language/c++"
	VISIBILITY_OBJC = "language/objc"
	VISIBILITY_OBJCXX = "language/objc++"

	KIND_FLAG = "flag"
	KIND_JOINED = "joined"
	KIND_SEPARATE = "separate"
	KIND_JOINED_OR_SEPARATE = "joined-or-separate"
	KIND_JOINED_AND_SEPARATE = "joined-and-separate"
	KIND_COMMA_JOINED = "comma-joined"
	KIND_MULTI_ARG = "multi-arg"
	KIND_REMAINING_ARGS = "remaining-args"
	KIND_REMAINING_ARGS_JOINED = "remaining-args-joined"
}

function get_kind(flags) {
	if (opt_args("Args", flags) != "") {
		return KIND_MULTI_ARG
	}
	if (flag_set_p("Joined", flags) && flag_set_p("Separate", flags)) {
		return KIND_JOINED_OR_SEPARATE
	}
	if (flag_set_p("JoinedOrMissing", flags)) {
		return KIND_JOINED
	}
	if (flag_set_p("Joined", flags)) {
		return KIND_JOINED
	}
	if (flag_set_p("Separate", flags)) {
		return KIND_SEPARATE
	}
	return KIND_FLAG
}

function get_num_args(kind, flags) {
	if (kind == "multi-arg") {
		return opt_args("Args", flags)
	}
	if (kind == "flag") {
		return 0
	}
	return 1
}

END {
	n_merged_opts = 0
	for (i = 0; i < n_opts; i++) {
		j = n_merged_opts
		merged_opts[j] = opts[i]
		merged_flags[j] = flags[i]
		n_merged_opts++
		while (i + 1 < n_opts && opts[i] == opts[i + 1]) {
			merged_flags[j] = merged_flags[j] " " flags[i + 1]
			i++
		}
	}

	for (i = 0; i < n_merged_opts; i++) {
		print quote merged_opts[i] quote ":"

		kind = get_kind(merged_flags[i])
		num_args = get_num_args(kind, merged_flags[i])

		print "  kind: " quote kind quote
		print "  numArgs: " num_args
	}
}
package tools

const (
	// Error code format: x.xxx.xxx?x
	// First 0 is used by system.
	// Next 001 - 100 is used by lapi.
	// 101 - 150 is used by validation
	// Tools package is using 151 - 170
	// Last 3 digits for error code, it might be extended to
	// whatever number if necessary

	ERR_TOOLS_LOAD_INI_FAILED          = "0.151.001"
	ERR_TOOLS_LOAD_INI_INVALID_SECTION = "0.151.002"

	ERR_TOOLS_STRUCT_READ_INVALID_TYPE = "0.152.001"
)

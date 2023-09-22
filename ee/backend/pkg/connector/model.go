package connector

import "strconv"

var sessionColumns = []string{
	"sessionid",
	"user_agent",
	"user_browser",
	"user_browser_version",
	"user_country",
	"user_device",
	"user_device_heap_size",
	"user_device_memory_size",
	"user_device_type",
	"user_os",
	"user_os_version",
	"user_uuid",
	"connection_effective_bandwidth",
	"connection_type",
	"metadata_key",
	"metadata_value",
	"referrer",
	"user_anonymous_id",
	"user_id",
	"session_start_timestamp",
	"session_end_timestamp",
	"session_duration",
	"first_contentful_paint",
	"speed_index",
	"visually_complete",
	"timing_time_to_interactive",
	"avg_cpu",
	"avg_fps",
	"max_cpu",
	"max_fps",
	"max_total_js_heap_size",
	"max_used_js_heap_size",
	"js_exceptions_count",
	"inputs_count",
	"clicks_count",
	"issues_count",
	"urls_count",
}

var sessionInts = []string{
	"user_device_heap_size",
	"user_device_memory_size",
	"connection_effective_bandwidth",
	"first_contentful_paint",
	"speed_index",
	"visually_complete",
	"timing_time_to_interactive",
	"avg_cpu",
	"avg_fps",
	"max_cpu",
	"max_fps",
	"max_total_js_heap_size",
	"max_used_js_heap_size",
	"js_exceptions_count",
	"inputs_count",
	"clicks_count",
	"issues_count",
	"urls_count",
}

var eventColumns = []string{
	"sessionid",
	"consolelog_level",
	"consolelog_value",
	"customevent_name",
	"customevent_payload",
	"jsexception_message",
	"jsexception_name",
	"jsexception_payload",
	"jsexception_metadata",
	"networkrequest_type",
	"networkrequest_method",
	"networkrequest_url",
	"networkrequest_request",
	"networkrequest_response",
	"networkrequest_status",
	"networkrequest_timestamp",
	"networkrequest_duration",
	"issueevent_message_id",
	"issueevent_timestamp",
	"issueevent_type",
	"issueevent_context_string",
	"issueevent_context",
	"issueevent_payload",
	"issueevent_url",
	"customissue_name",
	"customissue_payload",
	"received_at",
	"batch_order_number",
}

func QUOTES(s string) string {
	return strconv.Quote(s)
}

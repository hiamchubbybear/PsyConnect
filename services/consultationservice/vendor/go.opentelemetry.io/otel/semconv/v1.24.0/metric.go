




package semconv 

const (

	
	
	
	
	
	
	
	DBClientConnectionsUsageName        = "db.client.connections.usage"
	DBClientConnectionsUsageUnit        = "{connection}"
	DBClientConnectionsUsageDescription = "The number of connections that are currently in state described by the `state` attribute"

	
	
	
	
	
	
	DBClientConnectionsIdleMaxName        = "db.client.connections.idle.max"
	DBClientConnectionsIdleMaxUnit        = "{connection}"
	DBClientConnectionsIdleMaxDescription = "The maximum number of idle open connections allowed"

	
	
	
	
	
	
	DBClientConnectionsIdleMinName        = "db.client.connections.idle.min"
	DBClientConnectionsIdleMinUnit        = "{connection}"
	DBClientConnectionsIdleMinDescription = "The minimum number of idle open connections allowed"

	
	
	
	
	
	
	DBClientConnectionsMaxName        = "db.client.connections.max"
	DBClientConnectionsMaxUnit        = "{connection}"
	DBClientConnectionsMaxDescription = "The maximum number of open connections allowed"

	
	
	
	
	
	
	
	DBClientConnectionsPendingRequestsName        = "db.client.connections.pending_requests"
	DBClientConnectionsPendingRequestsUnit        = "{request}"
	DBClientConnectionsPendingRequestsDescription = "The number of pending requests for an open connection, cumulative for the entire pool"

	
	
	
	
	
	
	
	DBClientConnectionsTimeoutsName        = "db.client.connections.timeouts"
	DBClientConnectionsTimeoutsUnit        = "{timeout}"
	DBClientConnectionsTimeoutsDescription = "The number of connection timeouts that have occurred trying to obtain a connection from the pool"

	
	
	
	
	
	
	DBClientConnectionsCreateTimeName        = "db.client.connections.create_time"
	DBClientConnectionsCreateTimeUnit        = "ms"
	DBClientConnectionsCreateTimeDescription = "The time it took to create a new connection"

	
	
	
	
	
	
	DBClientConnectionsWaitTimeName        = "db.client.connections.wait_time"
	DBClientConnectionsWaitTimeUnit        = "ms"
	DBClientConnectionsWaitTimeDescription = "The time it took to obtain an open connection from the pool"

	
	
	
	
	
	
	DBClientConnectionsUseTimeName        = "db.client.connections.use_time"
	DBClientConnectionsUseTimeUnit        = "ms"
	DBClientConnectionsUseTimeDescription = "The time between borrowing a connection and returning it to the pool"

	
	
	
	
	
	
	AspnetcoreRoutingMatchAttemptsName        = "aspnetcore.routing.match_attempts"
	AspnetcoreRoutingMatchAttemptsUnit        = "{match_attempt}"
	AspnetcoreRoutingMatchAttemptsDescription = "Number of requests that were attempted to be matched to an endpoint."

	
	
	
	
	
	
	AspnetcoreDiagnosticsExceptionsName        = "aspnetcore.diagnostics.exceptions"
	AspnetcoreDiagnosticsExceptionsUnit        = "{exception}"
	AspnetcoreDiagnosticsExceptionsDescription = "Number of exceptions caught by exception handling middleware."

	
	
	
	
	
	
	
	AspnetcoreRateLimitingActiveRequestLeasesName        = "aspnetcore.rate_limiting.active_request_leases"
	AspnetcoreRateLimitingActiveRequestLeasesUnit        = "{request}"
	AspnetcoreRateLimitingActiveRequestLeasesDescription = "Number of requests that are currently active on the server that hold a rate limiting lease."

	
	
	
	
	
	
	
	AspnetcoreRateLimitingRequestLeaseDurationName        = "aspnetcore.rate_limiting.request_lease.duration"
	AspnetcoreRateLimitingRequestLeaseDurationUnit        = "s"
	AspnetcoreRateLimitingRequestLeaseDurationDescription = "The duration of rate limiting lease held by requests on the server."

	
	
	
	
	
	
	
	AspnetcoreRateLimitingRequestTimeInQueueName        = "aspnetcore.rate_limiting.request.time_in_queue"
	AspnetcoreRateLimitingRequestTimeInQueueUnit        = "s"
	AspnetcoreRateLimitingRequestTimeInQueueDescription = "The time the request spent in a queue waiting to acquire a rate limiting lease."

	
	
	
	
	
	
	
	AspnetcoreRateLimitingQueuedRequestsName        = "aspnetcore.rate_limiting.queued_requests"
	AspnetcoreRateLimitingQueuedRequestsUnit        = "{request}"
	AspnetcoreRateLimitingQueuedRequestsDescription = "Number of requests that are currently queued, waiting to acquire a rate limiting lease."

	
	
	
	
	
	
	AspnetcoreRateLimitingRequestsName        = "aspnetcore.rate_limiting.requests"
	AspnetcoreRateLimitingRequestsUnit        = "{request}"
	AspnetcoreRateLimitingRequestsDescription = "Number of requests that tried to acquire a rate limiting lease."

	
	
	
	
	
	
	DNSLookupDurationName        = "dns.lookup.duration"
	DNSLookupDurationUnit        = "s"
	DNSLookupDurationDescription = "Measures the time taken to perform a DNS lookup."

	
	
	
	
	
	
	
	HTTPClientOpenConnectionsName        = "http.client.open_connections"
	HTTPClientOpenConnectionsUnit        = "{connection}"
	HTTPClientOpenConnectionsDescription = "Number of outbound HTTP connections that are currently active or idle on the client."

	
	
	
	
	
	
	HTTPClientConnectionDurationName        = "http.client.connection.duration"
	HTTPClientConnectionDurationUnit        = "s"
	HTTPClientConnectionDurationDescription = "The duration of the successfully established outbound HTTP connections."

	
	
	
	
	
	
	HTTPClientActiveRequestsName        = "http.client.active_requests"
	HTTPClientActiveRequestsUnit        = "{request}"
	HTTPClientActiveRequestsDescription = "Number of active HTTP requests."

	
	
	
	
	
	
	
	HTTPClientRequestTimeInQueueName        = "http.client.request.time_in_queue"
	HTTPClientRequestTimeInQueueUnit        = "s"
	HTTPClientRequestTimeInQueueDescription = "The amount of time requests spent on a queue waiting for an available connection."

	
	
	
	
	
	
	KestrelActiveConnectionsName        = "kestrel.active_connections"
	KestrelActiveConnectionsUnit        = "{connection}"
	KestrelActiveConnectionsDescription = "Number of connections that are currently active on the server."

	
	
	
	
	
	
	KestrelConnectionDurationName        = "kestrel.connection.duration"
	KestrelConnectionDurationUnit        = "s"
	KestrelConnectionDurationDescription = "The duration of connections on the server."

	
	
	
	
	
	
	KestrelRejectedConnectionsName        = "kestrel.rejected_connections"
	KestrelRejectedConnectionsUnit        = "{connection}"
	KestrelRejectedConnectionsDescription = "Number of connections rejected by the server."

	
	
	
	
	
	
	KestrelQueuedConnectionsName        = "kestrel.queued_connections"
	KestrelQueuedConnectionsUnit        = "{connection}"
	KestrelQueuedConnectionsDescription = "Number of connections that are currently queued and are waiting to start."

	
	
	
	
	
	
	
	KestrelQueuedRequestsName        = "kestrel.queued_requests"
	KestrelQueuedRequestsUnit        = "{request}"
	KestrelQueuedRequestsDescription = "Number of HTTP requests on multiplexed connections (HTTP/2 and HTTP/3) that are currently queued and are waiting to start."

	
	
	
	
	
	
	KestrelUpgradedConnectionsName        = "kestrel.upgraded_connections"
	KestrelUpgradedConnectionsUnit        = "{connection}"
	KestrelUpgradedConnectionsDescription = "Number of connections that are currently upgraded (WebSockets). ."

	
	
	
	
	
	
	KestrelTLSHandshakeDurationName        = "kestrel.tls_handshake.duration"
	KestrelTLSHandshakeDurationUnit        = "s"
	KestrelTLSHandshakeDurationDescription = "The duration of TLS handshakes on the server."

	
	
	
	
	
	
	KestrelActiveTLSHandshakesName        = "kestrel.active_tls_handshakes"
	KestrelActiveTLSHandshakesUnit        = "{handshake}"
	KestrelActiveTLSHandshakesDescription = "Number of TLS handshakes that are currently in progress on the server."

	
	
	
	
	
	
	SignalrServerConnectionDurationName        = "signalr.server.connection.duration"
	SignalrServerConnectionDurationUnit        = "s"
	SignalrServerConnectionDurationDescription = "The duration of connections on the server."

	
	
	
	
	
	
	SignalrServerActiveConnectionsName        = "signalr.server.active_connections"
	SignalrServerActiveConnectionsUnit        = "{connection}"
	SignalrServerActiveConnectionsDescription = "Number of connections that are currently active on the server."

	
	
	
	
	
	
	FaaSInvokeDurationName        = "faas.invoke_duration"
	FaaSInvokeDurationUnit        = "s"
	FaaSInvokeDurationDescription = "Measures the duration of the function's logic execution"

	
	
	
	
	
	
	FaaSInitDurationName        = "faas.init_duration"
	FaaSInitDurationUnit        = "s"
	FaaSInitDurationDescription = "Measures the duration of the function's initialization, such as a cold start"

	
	
	
	
	
	FaaSColdstartsName        = "faas.coldstarts"
	FaaSColdstartsUnit        = "{coldstart}"
	FaaSColdstartsDescription = "Number of invocation cold starts"

	
	
	
	
	
	FaaSErrorsName        = "faas.errors"
	FaaSErrorsUnit        = "{error}"
	FaaSErrorsDescription = "Number of invocation errors"

	
	
	
	
	
	FaaSInvocationsName        = "faas.invocations"
	FaaSInvocationsUnit        = "{invocation}"
	FaaSInvocationsDescription = "Number of successful invocations"

	
	
	
	
	
	FaaSTimeoutsName        = "faas.timeouts"
	FaaSTimeoutsUnit        = "{timeout}"
	FaaSTimeoutsDescription = "Number of invocation timeouts"

	
	
	
	
	
	
	FaaSMemUsageName        = "faas.mem_usage"
	FaaSMemUsageUnit        = "By"
	FaaSMemUsageDescription = "Distribution of max memory usage per invocation"

	
	
	
	
	
	FaaSCPUUsageName        = "faas.cpu_usage"
	FaaSCPUUsageUnit        = "s"
	FaaSCPUUsageDescription = "Distribution of CPU usage per invocation"

	
	
	
	
	
	FaaSNetIoName        = "faas.net_io"
	FaaSNetIoUnit        = "By"
	FaaSNetIoDescription = "Distribution of net I/O usage per invocation"

	
	
	
	
	
	
	HTTPServerRequestDurationName        = "http.server.request.duration"
	HTTPServerRequestDurationUnit        = "s"
	HTTPServerRequestDurationDescription = "Duration of HTTP server requests."

	
	
	
	
	
	
	HTTPServerActiveRequestsName        = "http.server.active_requests"
	HTTPServerActiveRequestsUnit        = "{request}"
	HTTPServerActiveRequestsDescription = "Number of active HTTP server requests."

	
	
	
	
	
	
	HTTPServerRequestBodySizeName        = "http.server.request.body.size"
	HTTPServerRequestBodySizeUnit        = "By"
	HTTPServerRequestBodySizeDescription = "Size of HTTP server request bodies."

	
	
	
	
	
	
	HTTPServerResponseBodySizeName        = "http.server.response.body.size"
	HTTPServerResponseBodySizeUnit        = "By"
	HTTPServerResponseBodySizeDescription = "Size of HTTP server response bodies."

	
	
	
	
	
	
	HTTPClientRequestDurationName        = "http.client.request.duration"
	HTTPClientRequestDurationUnit        = "s"
	HTTPClientRequestDurationDescription = "Duration of HTTP client requests."

	
	
	
	
	
	
	HTTPClientRequestBodySizeName        = "http.client.request.body.size"
	HTTPClientRequestBodySizeUnit        = "By"
	HTTPClientRequestBodySizeDescription = "Size of HTTP client request bodies."

	
	
	
	
	
	
	HTTPClientResponseBodySizeName        = "http.client.response.body.size"
	HTTPClientResponseBodySizeUnit        = "By"
	HTTPClientResponseBodySizeDescription = "Size of HTTP client response bodies."

	
	
	
	
	
	JvmMemoryInitName        = "jvm.memory.init"
	JvmMemoryInitUnit        = "By"
	JvmMemoryInitDescription = "Measure of initial memory requested."

	
	
	
	
	
	
	JvmSystemCPUUtilizationName        = "jvm.system.cpu.utilization"
	JvmSystemCPUUtilizationUnit        = "1"
	JvmSystemCPUUtilizationDescription = "Recent CPU utilization for the whole system as reported by the JVM."

	
	
	
	
	
	
	JvmSystemCPULoad1mName        = "jvm.system.cpu.load_1m"
	JvmSystemCPULoad1mUnit        = "{run_queue_item}"
	JvmSystemCPULoad1mDescription = "Average CPU load of the whole system for the last minute as reported by the JVM."

	
	
	
	
	
	
	JvmBufferMemoryUsageName        = "jvm.buffer.memory.usage"
	JvmBufferMemoryUsageUnit        = "By"
	JvmBufferMemoryUsageDescription = "Measure of memory used by buffers."

	
	
	
	
	
	
	JvmBufferMemoryLimitName        = "jvm.buffer.memory.limit"
	JvmBufferMemoryLimitUnit        = "By"
	JvmBufferMemoryLimitDescription = "Measure of total memory capacity of buffers."

	
	
	
	
	
	JvmBufferCountName        = "jvm.buffer.count"
	JvmBufferCountUnit        = "{buffer}"
	JvmBufferCountDescription = "Number of buffers in the pool."

	
	
	
	
	
	JvmMemoryUsedName        = "jvm.memory.used"
	JvmMemoryUsedUnit        = "By"
	JvmMemoryUsedDescription = "Measure of memory used."

	
	
	
	
	
	JvmMemoryCommittedName        = "jvm.memory.committed"
	JvmMemoryCommittedUnit        = "By"
	JvmMemoryCommittedDescription = "Measure of memory committed."

	
	
	
	
	
	JvmMemoryLimitName        = "jvm.memory.limit"
	JvmMemoryLimitUnit        = "By"
	JvmMemoryLimitDescription = "Measure of max obtainable memory."

	
	
	
	
	
	
	
	JvmMemoryUsedAfterLastGcName        = "jvm.memory.used_after_last_gc"
	JvmMemoryUsedAfterLastGcUnit        = "By"
	JvmMemoryUsedAfterLastGcDescription = "Measure of memory used, as measured after the most recent garbage collection event on this pool."

	
	
	
	
	
	JvmGcDurationName        = "jvm.gc.duration"
	JvmGcDurationUnit        = "s"
	JvmGcDurationDescription = "Duration of JVM garbage collection actions."

	
	
	
	
	
	JvmThreadCountName        = "jvm.thread.count"
	JvmThreadCountUnit        = "{thread}"
	JvmThreadCountDescription = "Number of executing platform threads."

	
	
	
	
	
	JvmClassLoadedName        = "jvm.class.loaded"
	JvmClassLoadedUnit        = "{class}"
	JvmClassLoadedDescription = "Number of classes loaded since JVM start."

	
	
	
	
	
	
	JvmClassUnloadedName        = "jvm.class.unloaded"
	JvmClassUnloadedUnit        = "{class}"
	JvmClassUnloadedDescription = "Number of classes unloaded since JVM start."

	
	
	
	
	
	JvmClassCountName        = "jvm.class.count"
	JvmClassCountUnit        = "{class}"
	JvmClassCountDescription = "Number of classes currently loaded."

	
	
	
	
	
	
	JvmCPUCountName        = "jvm.cpu.count"
	JvmCPUCountUnit        = "{cpu}"
	JvmCPUCountDescription = "Number of processors available to the Java virtual machine."

	
	
	
	
	
	
	JvmCPUTimeName        = "jvm.cpu.time"
	JvmCPUTimeUnit        = "s"
	JvmCPUTimeDescription = "CPU time used by the process as reported by the JVM."

	
	
	
	
	
	
	JvmCPURecentUtilizationName        = "jvm.cpu.recent_utilization"
	JvmCPURecentUtilizationUnit        = "1"
	JvmCPURecentUtilizationDescription = "Recent CPU utilization for the process as reported by the JVM."

	
	
	
	
	
	
	MessagingPublishDurationName        = "messaging.publish.duration"
	MessagingPublishDurationUnit        = "s"
	MessagingPublishDurationDescription = "Measures the duration of publish operation."

	
	
	
	
	
	
	MessagingReceiveDurationName        = "messaging.receive.duration"
	MessagingReceiveDurationUnit        = "s"
	MessagingReceiveDurationDescription = "Measures the duration of receive operation."

	
	
	
	
	
	
	MessagingDeliverDurationName        = "messaging.deliver.duration"
	MessagingDeliverDurationUnit        = "s"
	MessagingDeliverDurationDescription = "Measures the duration of deliver operation."

	
	
	
	
	
	
	MessagingPublishMessagesName        = "messaging.publish.messages"
	MessagingPublishMessagesUnit        = "{message}"
	MessagingPublishMessagesDescription = "Measures the number of published messages."

	
	
	
	
	
	
	MessagingReceiveMessagesName        = "messaging.receive.messages"
	MessagingReceiveMessagesUnit        = "{message}"
	MessagingReceiveMessagesDescription = "Measures the number of received messages."

	
	
	
	
	
	
	MessagingDeliverMessagesName        = "messaging.deliver.messages"
	MessagingDeliverMessagesUnit        = "{message}"
	MessagingDeliverMessagesDescription = "Measures the number of delivered messages."

	
	
	
	
	
	
	RPCServerDurationName        = "rpc.server.duration"
	RPCServerDurationUnit        = "ms"
	RPCServerDurationDescription = "Measures the duration of inbound RPC."

	
	
	
	
	
	
	RPCServerRequestSizeName        = "rpc.server.request.size"
	RPCServerRequestSizeUnit        = "By"
	RPCServerRequestSizeDescription = "Measures the size of RPC request messages (uncompressed)."

	
	
	
	
	
	
	RPCServerResponseSizeName        = "rpc.server.response.size"
	RPCServerResponseSizeUnit        = "By"
	RPCServerResponseSizeDescription = "Measures the size of RPC response messages (uncompressed)."

	
	
	
	
	
	
	RPCServerRequestsPerRPCName        = "rpc.server.requests_per_rpc"
	RPCServerRequestsPerRPCUnit        = "{count}"
	RPCServerRequestsPerRPCDescription = "Measures the number of messages received per RPC."

	
	
	
	
	
	
	RPCServerResponsesPerRPCName        = "rpc.server.responses_per_rpc"
	RPCServerResponsesPerRPCUnit        = "{count}"
	RPCServerResponsesPerRPCDescription = "Measures the number of messages sent per RPC."

	
	
	
	
	
	
	RPCClientDurationName        = "rpc.client.duration"
	RPCClientDurationUnit        = "ms"
	RPCClientDurationDescription = "Measures the duration of outbound RPC."

	
	
	
	
	
	
	RPCClientRequestSizeName        = "rpc.client.request.size"
	RPCClientRequestSizeUnit        = "By"
	RPCClientRequestSizeDescription = "Measures the size of RPC request messages (uncompressed)."

	
	
	
	
	
	
	RPCClientResponseSizeName        = "rpc.client.response.size"
	RPCClientResponseSizeUnit        = "By"
	RPCClientResponseSizeDescription = "Measures the size of RPC response messages (uncompressed)."

	
	
	
	
	
	
	RPCClientRequestsPerRPCName        = "rpc.client.requests_per_rpc"
	RPCClientRequestsPerRPCUnit        = "{count}"
	RPCClientRequestsPerRPCDescription = "Measures the number of messages received per RPC."

	
	
	
	
	
	
	RPCClientResponsesPerRPCName        = "rpc.client.responses_per_rpc"
	RPCClientResponsesPerRPCUnit        = "{count}"
	RPCClientResponsesPerRPCDescription = "Measures the number of messages sent per RPC."

	
	
	
	
	
	SystemCPUTimeName        = "system.cpu.time"
	SystemCPUTimeUnit        = "s"
	SystemCPUTimeDescription = "Seconds each logical CPU spent on each mode"

	
	
	
	
	
	
	
	SystemCPUUtilizationName        = "system.cpu.utilization"
	SystemCPUUtilizationUnit        = "1"
	SystemCPUUtilizationDescription = "Difference in system.cpu.time since the last measurement, divided by the elapsed time and number of logical CPUs"

	
	
	
	
	
	
	SystemCPUFrequencyName        = "system.cpu.frequency"
	SystemCPUFrequencyUnit        = "{Hz}"
	SystemCPUFrequencyDescription = "Reports the current frequency of the CPU in Hz"

	
	
	
	
	
	
	SystemCPUPhysicalCountName        = "system.cpu.physical.count"
	SystemCPUPhysicalCountUnit        = "{cpu}"
	SystemCPUPhysicalCountDescription = "Reports the number of actual physical processor cores on the hardware"

	
	
	
	
	
	
	
	SystemCPULogicalCountName        = "system.cpu.logical.count"
	SystemCPULogicalCountUnit        = "{cpu}"
	SystemCPULogicalCountDescription = "Reports the number of logical (virtual) processor cores created by the operating system to manage multitasking"

	
	
	
	
	
	SystemMemoryUsageName        = "system.memory.usage"
	SystemMemoryUsageUnit        = "By"
	SystemMemoryUsageDescription = "Reports memory in use by state."

	
	
	
	
	
	
	SystemMemoryLimitName        = "system.memory.limit"
	SystemMemoryLimitUnit        = "By"
	SystemMemoryLimitDescription = "Total memory available in the system."

	
	
	
	
	
	
	SystemMemoryUtilizationName = "system.memory.utilization"
	SystemMemoryUtilizationUnit = "1"

	
	
	
	
	
	SystemPagingUsageName        = "system.paging.usage"
	SystemPagingUsageUnit        = "By"
	SystemPagingUsageDescription = "Unix swap or windows pagefile usage"

	
	
	
	
	
	
	SystemPagingUtilizationName = "system.paging.utilization"
	SystemPagingUtilizationUnit = "1"

	
	
	
	
	
	
	SystemPagingFaultsName = "system.paging.faults"
	SystemPagingFaultsUnit = "{fault}"

	
	
	
	
	
	
	SystemPagingOperationsName = "system.paging.operations"
	SystemPagingOperationsUnit = "{operation}"

	
	
	
	
	
	
	SystemDiskIoName = "system.disk.io"
	SystemDiskIoUnit = "By"

	
	
	
	
	
	
	SystemDiskOperationsName = "system.disk.operations"
	SystemDiskOperationsUnit = "{operation}"

	
	
	
	
	
	SystemDiskIoTimeName        = "system.disk.io_time"
	SystemDiskIoTimeUnit        = "s"
	SystemDiskIoTimeDescription = "Time disk spent activated"

	
	
	
	
	
	
	SystemDiskOperationTimeName        = "system.disk.operation_time"
	SystemDiskOperationTimeUnit        = "s"
	SystemDiskOperationTimeDescription = "Sum of the time each operation took to complete"

	
	
	
	
	
	
	SystemDiskMergedName = "system.disk.merged"
	SystemDiskMergedUnit = "{operation}"

	
	
	
	
	
	
	SystemFilesystemUsageName = "system.filesystem.usage"
	SystemFilesystemUsageUnit = "By"

	
	
	
	
	
	
	SystemFilesystemUtilizationName = "system.filesystem.utilization"
	SystemFilesystemUtilizationUnit = "1"

	
	
	
	
	
	
	SystemNetworkDroppedName        = "system.network.dropped"
	SystemNetworkDroppedUnit        = "{packet}"
	SystemNetworkDroppedDescription = "Count of packets that are dropped or discarded even though there was no error"

	
	
	
	
	
	
	SystemNetworkPacketsName = "system.network.packets"
	SystemNetworkPacketsUnit = "{packet}"

	
	
	
	
	
	SystemNetworkErrorsName        = "system.network.errors"
	SystemNetworkErrorsUnit        = "{error}"
	SystemNetworkErrorsDescription = "Count of network errors detected"

	
	
	
	
	
	
	SystemNetworkIoName = "system.network.io"
	SystemNetworkIoUnit = "By"

	
	
	
	
	
	
	SystemNetworkConnectionsName = "system.network.connections"
	SystemNetworkConnectionsUnit = "{connection}"

	
	
	
	
	
	
	SystemProcessesCountName        = "system.processes.count"
	SystemProcessesCountUnit        = "{process}"
	SystemProcessesCountDescription = "Total number of processes in each state"

	
	
	
	
	
	
	SystemProcessesCreatedName        = "system.processes.created"
	SystemProcessesCreatedUnit        = "{process}"
	SystemProcessesCreatedDescription = "Total number of processes created over uptime of the host"

	
	
	
	
	
	
	
	SystemLinuxMemoryAvailableName        = "system.linux.memory.available"
	SystemLinuxMemoryAvailableUnit        = "By"
	SystemLinuxMemoryAvailableDescription = "An estimate of how much memory is available for starting new applications, without causing swapping"
)

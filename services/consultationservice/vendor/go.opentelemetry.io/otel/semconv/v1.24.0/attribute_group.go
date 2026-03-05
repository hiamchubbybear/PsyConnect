




package semconv 

import "go.opentelemetry.io/otel/attribute"


const (
	
	
	
	
	
	
	
	
	
	
	FaaSInvokedNameKey = attribute.Key("faas.invoked_name")

	
	
	
	
	
	
	
	
	
	FaaSInvokedProviderKey = attribute.Key("faas.invoked_provider")

	
	
	
	
	
	
	
	
	
	
	
	
	
	
	
	
	FaaSInvokedRegionKey = attribute.Key("faas.invoked_region")

	
	
	
	
	
	
	
	FaaSTriggerKey = attribute.Key("faas.trigger")
)

var (
	
	FaaSInvokedProviderAlibabaCloud = FaaSInvokedProviderKey.String("alibaba_cloud")
	
	FaaSInvokedProviderAWS = FaaSInvokedProviderKey.String("aws")
	
	FaaSInvokedProviderAzure = FaaSInvokedProviderKey.String("azure")
	
	FaaSInvokedProviderGCP = FaaSInvokedProviderKey.String("gcp")
	
	FaaSInvokedProviderTencentCloud = FaaSInvokedProviderKey.String("tencent_cloud")
)

var (
	
	FaaSTriggerDatasource = FaaSTriggerKey.String("datasource")
	
	FaaSTriggerHTTP = FaaSTriggerKey.String("http")
	
	FaaSTriggerPubsub = FaaSTriggerKey.String("pubsub")
	
	FaaSTriggerTimer = FaaSTriggerKey.String("timer")
	
	FaaSTriggerOther = FaaSTriggerKey.String("other")
)




func FaaSInvokedName(val string) attribute.KeyValue {
	return FaaSInvokedNameKey.String(val)
}




func FaaSInvokedRegion(val string) attribute.KeyValue {
	return FaaSInvokedRegionKey.String(val)
}


const (
	
	
	
	
	
	
	
	
	
	
	
	
	
	EventNameKey = attribute.Key("event.name")
)




func EventName(val string) attribute.KeyValue {
	return EventNameKey.String(val)
}



const (
	
	
	
	
	
	
	
	
	
	
	
	
	
	
	LogRecordUIDKey = attribute.Key("log.record.uid")
)




func LogRecordUID(val string) attribute.KeyValue {
	return LogRecordUIDKey.String(val)
}


const (
	
	
	
	
	
	
	
	LogIostreamKey = attribute.Key("log.iostream")
)

var (
	
	LogIostreamStdout = LogIostreamKey.String("stdout")
	
	LogIostreamStderr = LogIostreamKey.String("stderr")
)


const (
	
	
	
	
	
	
	
	LogFileNameKey = attribute.Key("log.file.name")

	
	
	
	
	
	
	
	
	LogFileNameResolvedKey = attribute.Key("log.file.name_resolved")

	
	
	
	
	
	
	
	LogFilePathKey = attribute.Key("log.file.path")

	
	
	
	
	
	
	
	
	LogFilePathResolvedKey = attribute.Key("log.file.path_resolved")
)




func LogFileName(val string) attribute.KeyValue {
	return LogFileNameKey.String(val)
}




func LogFileNameResolved(val string) attribute.KeyValue {
	return LogFileNameResolvedKey.String(val)
}




func LogFilePath(val string) attribute.KeyValue {
	return LogFilePathKey.String(val)
}




func LogFilePathResolved(val string) attribute.KeyValue {
	return LogFilePathResolvedKey.String(val)
}


const (
	
	
	
	
	
	
	
	
	
	
	
	PoolNameKey = attribute.Key("pool.name")

	
	
	
	
	
	
	
	StateKey = attribute.Key("state")
)

var (
	
	StateIdle = StateKey.String("idle")
	
	StateUsed = StateKey.String("used")
)







func PoolName(val string) attribute.KeyValue {
	return PoolNameKey.String(val)
}


const (
	
	
	
	
	
	
	
	
	
	
	
	AspnetcoreDiagnosticsHandlerTypeKey = attribute.Key("aspnetcore.diagnostics.handler.type")

	
	
	
	
	
	
	
	
	
	AspnetcoreRateLimitingPolicyKey = attribute.Key("aspnetcore.rate_limiting.policy")

	
	
	
	
	
	
	
	
	
	AspnetcoreRateLimitingResultKey = attribute.Key("aspnetcore.rate_limiting.result")

	
	
	
	
	
	
	
	
	
	AspnetcoreRequestIsUnhandledKey = attribute.Key("aspnetcore.request.is_unhandled")

	
	
	
	
	
	
	
	
	
	AspnetcoreRoutingIsFallbackKey = attribute.Key("aspnetcore.routing.is_fallback")
)

var (
	
	AspnetcoreRateLimitingResultAcquired = AspnetcoreRateLimitingResultKey.String("acquired")
	
	AspnetcoreRateLimitingResultEndpointLimiter = AspnetcoreRateLimitingResultKey.String("endpoint_limiter")
	
	AspnetcoreRateLimitingResultGlobalLimiter = AspnetcoreRateLimitingResultKey.String("global_limiter")
	
	AspnetcoreRateLimitingResultRequestCanceled = AspnetcoreRateLimitingResultKey.String("request_canceled")
)






func AspnetcoreDiagnosticsHandlerType(val string) attribute.KeyValue {
	return AspnetcoreDiagnosticsHandlerTypeKey.String(val)
}




func AspnetcoreRateLimitingPolicy(val string) attribute.KeyValue {
	return AspnetcoreRateLimitingPolicyKey.String(val)
}




func AspnetcoreRequestIsUnhandled(val bool) attribute.KeyValue {
	return AspnetcoreRequestIsUnhandledKey.Bool(val)
}




func AspnetcoreRoutingIsFallback(val bool) attribute.KeyValue {
	return AspnetcoreRoutingIsFallbackKey.Bool(val)
}


const (
	
	
	
	
	
	
	
	
	SignalrConnectionStatusKey = attribute.Key("signalr.connection.status")

	
	
	
	
	
	
	
	
	
	SignalrTransportKey = attribute.Key("signalr.transport")
)

var (
	
	SignalrConnectionStatusNormalClosure = SignalrConnectionStatusKey.String("normal_closure")
	
	SignalrConnectionStatusTimeout = SignalrConnectionStatusKey.String("timeout")
	
	SignalrConnectionStatusAppShutdown = SignalrConnectionStatusKey.String("app_shutdown")
)

var (
	
	SignalrTransportServerSentEvents = SignalrTransportKey.String("server_sent_events")
	
	SignalrTransportLongPolling = SignalrTransportKey.String("long_polling")
	
	SignalrTransportWebSockets = SignalrTransportKey.String("web_sockets")
)


const (
	
	
	
	
	
	
	
	
	
	
	JvmBufferPoolNameKey = attribute.Key("jvm.buffer.pool.name")
)




func JvmBufferPoolName(val string) attribute.KeyValue {
	return JvmBufferPoolNameKey.String(val)
}


const (
	
	
	
	
	
	
	
	
	
	
	JvmMemoryPoolNameKey = attribute.Key("jvm.memory.pool.name")

	
	
	
	
	
	
	
	
	JvmMemoryTypeKey = attribute.Key("jvm.memory.type")
)

var (
	
	JvmMemoryTypeHeap = JvmMemoryTypeKey.String("heap")
	
	JvmMemoryTypeNonHeap = JvmMemoryTypeKey.String("non_heap")
)




func JvmMemoryPoolName(val string) attribute.KeyValue {
	return JvmMemoryPoolNameKey.String(val)
}


const (
	
	
	
	
	
	
	
	SystemDeviceKey = attribute.Key("system.device")
)



func SystemDevice(val string) attribute.KeyValue {
	return SystemDeviceKey.String(val)
}


const (
	
	
	
	
	
	
	
	
	SystemCPULogicalNumberKey = attribute.Key("system.cpu.logical_number")

	
	
	
	
	
	
	
	
	SystemCPUStateKey = attribute.Key("system.cpu.state")
)

var (
	
	SystemCPUStateUser = SystemCPUStateKey.String("user")
	
	SystemCPUStateSystem = SystemCPUStateKey.String("system")
	
	SystemCPUStateNice = SystemCPUStateKey.String("nice")
	
	SystemCPUStateIdle = SystemCPUStateKey.String("idle")
	
	SystemCPUStateIowait = SystemCPUStateKey.String("iowait")
	
	SystemCPUStateInterrupt = SystemCPUStateKey.String("interrupt")
	
	SystemCPUStateSteal = SystemCPUStateKey.String("steal")
)




func SystemCPULogicalNumber(val int) attribute.KeyValue {
	return SystemCPULogicalNumberKey.Int(val)
}


const (
	
	
	
	
	
	
	
	
	SystemMemoryStateKey = attribute.Key("system.memory.state")
)

var (
	
	SystemMemoryStateUsed = SystemMemoryStateKey.String("used")
	
	SystemMemoryStateFree = SystemMemoryStateKey.String("free")
	
	SystemMemoryStateShared = SystemMemoryStateKey.String("shared")
	
	SystemMemoryStateBuffers = SystemMemoryStateKey.String("buffers")
	
	SystemMemoryStateCached = SystemMemoryStateKey.String("cached")
)


const (
	
	
	
	
	
	
	
	
	SystemPagingDirectionKey = attribute.Key("system.paging.direction")

	
	
	
	
	
	
	
	
	SystemPagingStateKey = attribute.Key("system.paging.state")

	
	
	
	
	
	
	
	
	SystemPagingTypeKey = attribute.Key("system.paging.type")
)

var (
	
	SystemPagingDirectionIn = SystemPagingDirectionKey.String("in")
	
	SystemPagingDirectionOut = SystemPagingDirectionKey.String("out")
)

var (
	
	SystemPagingStateUsed = SystemPagingStateKey.String("used")
	
	SystemPagingStateFree = SystemPagingStateKey.String("free")
)

var (
	
	SystemPagingTypeMajor = SystemPagingTypeKey.String("major")
	
	SystemPagingTypeMinor = SystemPagingTypeKey.String("minor")
)


const (
	
	
	
	
	
	
	
	
	SystemFilesystemModeKey = attribute.Key("system.filesystem.mode")

	
	
	
	
	
	
	
	
	SystemFilesystemMountpointKey = attribute.Key("system.filesystem.mountpoint")

	
	
	
	
	
	
	
	
	SystemFilesystemStateKey = attribute.Key("system.filesystem.state")

	
	
	
	
	
	
	
	
	SystemFilesystemTypeKey = attribute.Key("system.filesystem.type")
)

var (
	
	SystemFilesystemStateUsed = SystemFilesystemStateKey.String("used")
	
	SystemFilesystemStateFree = SystemFilesystemStateKey.String("free")
	
	SystemFilesystemStateReserved = SystemFilesystemStateKey.String("reserved")
)

var (
	
	SystemFilesystemTypeFat32 = SystemFilesystemTypeKey.String("fat32")
	
	SystemFilesystemTypeExfat = SystemFilesystemTypeKey.String("exfat")
	
	SystemFilesystemTypeNtfs = SystemFilesystemTypeKey.String("ntfs")
	
	SystemFilesystemTypeRefs = SystemFilesystemTypeKey.String("refs")
	
	SystemFilesystemTypeHfsplus = SystemFilesystemTypeKey.String("hfsplus")
	
	SystemFilesystemTypeExt4 = SystemFilesystemTypeKey.String("ext4")
)




func SystemFilesystemMode(val string) attribute.KeyValue {
	return SystemFilesystemModeKey.String(val)
}




func SystemFilesystemMountpoint(val string) attribute.KeyValue {
	return SystemFilesystemMountpointKey.String(val)
}


const (
	
	
	
	
	
	
	
	
	SystemNetworkStateKey = attribute.Key("system.network.state")
)

var (
	
	SystemNetworkStateClose = SystemNetworkStateKey.String("close")
	
	SystemNetworkStateCloseWait = SystemNetworkStateKey.String("close_wait")
	
	SystemNetworkStateClosing = SystemNetworkStateKey.String("closing")
	
	SystemNetworkStateDelete = SystemNetworkStateKey.String("delete")
	
	SystemNetworkStateEstablished = SystemNetworkStateKey.String("established")
	
	SystemNetworkStateFinWait1 = SystemNetworkStateKey.String("fin_wait_1")
	
	SystemNetworkStateFinWait2 = SystemNetworkStateKey.String("fin_wait_2")
	
	SystemNetworkStateLastAck = SystemNetworkStateKey.String("last_ack")
	
	SystemNetworkStateListen = SystemNetworkStateKey.String("listen")
	
	SystemNetworkStateSynRecv = SystemNetworkStateKey.String("syn_recv")
	
	SystemNetworkStateSynSent = SystemNetworkStateKey.String("syn_sent")
	
	SystemNetworkStateTimeWait = SystemNetworkStateKey.String("time_wait")
)


const (
	
	
	
	
	
	
	
	
	
	SystemProcessesStatusKey = attribute.Key("system.processes.status")
)

var (
	
	SystemProcessesStatusRunning = SystemProcessesStatusKey.String("running")
	
	SystemProcessesStatusSleeping = SystemProcessesStatusKey.String("sleeping")
	
	SystemProcessesStatusStopped = SystemProcessesStatusKey.String("stopped")
	
	SystemProcessesStatusDefunct = SystemProcessesStatusKey.String("defunct")
)









const (
	
	
	
	
	
	
	
	
	
	
	
	
	ClientAddressKey = attribute.Key("client.address")

	
	
	
	
	
	
	
	
	
	
	ClientPortKey = attribute.Key("client.port")
)





func ClientAddress(val string) attribute.KeyValue {
	return ClientAddressKey.String(val)
}



func ClientPort(val int) attribute.KeyValue {
	return ClientPortKey.Int(val)
}


const (
	
	
	
	
	
	
	
	
	DBCassandraConsistencyLevelKey = attribute.Key("db.cassandra.consistency_level")

	
	
	
	
	
	
	
	
	DBCassandraCoordinatorDCKey = attribute.Key("db.cassandra.coordinator.dc")

	
	
	
	
	
	
	
	
	DBCassandraCoordinatorIDKey = attribute.Key("db.cassandra.coordinator.id")

	
	
	
	
	
	
	
	DBCassandraIdempotenceKey = attribute.Key("db.cassandra.idempotence")

	
	
	
	
	
	
	
	
	DBCassandraPageSizeKey = attribute.Key("db.cassandra.page_size")

	
	
	
	
	
	
	
	
	
	DBCassandraSpeculativeExecutionCountKey = attribute.Key("db.cassandra.speculative_execution_count")

	
	
	
	
	
	
	
	
	
	
	
	
	
	
	
	DBCassandraTableKey = attribute.Key("db.cassandra.table")

	
	
	
	
	
	
	
	
	
	DBConnectionStringKey = attribute.Key("db.connection_string")

	
	
	
	
	
	
	
	
	DBCosmosDBClientIDKey = attribute.Key("db.cosmosdb.client_id")

	
	
	
	
	
	
	
	DBCosmosDBConnectionModeKey = attribute.Key("db.cosmosdb.connection_mode")

	
	
	
	
	
	
	
	
	DBCosmosDBContainerKey = attribute.Key("db.cosmosdb.container")

	
	
	
	
	
	
	
	DBCosmosDBOperationTypeKey = attribute.Key("db.cosmosdb.operation_type")

	
	
	
	
	
	
	
	
	DBCosmosDBRequestChargeKey = attribute.Key("db.cosmosdb.request_charge")

	
	
	
	
	
	
	
	DBCosmosDBRequestContentLengthKey = attribute.Key("db.cosmosdb.request_content_length")

	
	
	
	
	
	
	
	
	DBCosmosDBStatusCodeKey = attribute.Key("db.cosmosdb.status_code")

	
	
	
	
	
	
	
	
	DBCosmosDBSubStatusCodeKey = attribute.Key("db.cosmosdb.sub_status_code")

	
	
	
	
	
	
	
	
	DBElasticsearchClusterNameKey = attribute.Key("db.elasticsearch.cluster.name")

	
	
	
	
	
	
	
	
	
	DBElasticsearchNodeNameKey = attribute.Key("db.elasticsearch.node.name")

	
	
	
	
	
	
	
	
	
	
	
	
	
	DBInstanceIDKey = attribute.Key("db.instance.id")

	
	
	
	
	
	
	
	
	
	
	
	DBJDBCDriverClassnameKey = attribute.Key("db.jdbc.driver_classname")

	
	
	
	
	
	
	
	
	DBMongoDBCollectionKey = attribute.Key("db.mongodb.collection")

	
	
	
	
	
	
	
	
	
	
	
	
	
	DBMSSQLInstanceNameKey = attribute.Key("db.mssql.instance_name")

	
	
	
	
	
	
	
	
	
	
	
	
	
	
	DBNameKey = attribute.Key("db.name")

	
	
	
	
	
	
	
	
	
	
	
	
	
	
	
	
	DBOperationKey = attribute.Key("db.operation")

	
	
	
	
	
	
	
	
	
	
	DBRedisDBIndexKey = attribute.Key("db.redis.database_index")

	
	
	
	
	
	
	
	
	
	
	
	
	
	
	DBSQLTableKey = attribute.Key("db.sql.table")

	
	
	
	
	
	
	
	
	DBStatementKey = attribute.Key("db.statement")

	
	
	
	
	
	
	
	
	DBSystemKey = attribute.Key("db.system")

	
	
	
	
	
	
	
	DBUserKey = attribute.Key("db.user")
)

var (
	
	DBCassandraConsistencyLevelAll = DBCassandraConsistencyLevelKey.String("all")
	
	DBCassandraConsistencyLevelEachQuorum = DBCassandraConsistencyLevelKey.String("each_quorum")
	
	DBCassandraConsistencyLevelQuorum = DBCassandraConsistencyLevelKey.String("quorum")
	
	DBCassandraConsistencyLevelLocalQuorum = DBCassandraConsistencyLevelKey.String("local_quorum")
	
	DBCassandraConsistencyLevelOne = DBCassandraConsistencyLevelKey.String("one")
	
	DBCassandraConsistencyLevelTwo = DBCassandraConsistencyLevelKey.String("two")
	
	DBCassandraConsistencyLevelThree = DBCassandraConsistencyLevelKey.String("three")
	
	DBCassandraConsistencyLevelLocalOne = DBCassandraConsistencyLevelKey.String("local_one")
	
	DBCassandraConsistencyLevelAny = DBCassandraConsistencyLevelKey.String("any")
	
	DBCassandraConsistencyLevelSerial = DBCassandraConsistencyLevelKey.String("serial")
	
	DBCassandraConsistencyLevelLocalSerial = DBCassandraConsistencyLevelKey.String("local_serial")
)

var (
	
	DBCosmosDBConnectionModeGateway = DBCosmosDBConnectionModeKey.String("gateway")
	
	DBCosmosDBConnectionModeDirect = DBCosmosDBConnectionModeKey.String("direct")
)

var (
	
	DBCosmosDBOperationTypeInvalid = DBCosmosDBOperationTypeKey.String("Invalid")
	
	DBCosmosDBOperationTypeCreate = DBCosmosDBOperationTypeKey.String("Create")
	
	DBCosmosDBOperationTypePatch = DBCosmosDBOperationTypeKey.String("Patch")
	
	DBCosmosDBOperationTypeRead = DBCosmosDBOperationTypeKey.String("Read")
	
	DBCosmosDBOperationTypeReadFeed = DBCosmosDBOperationTypeKey.String("ReadFeed")
	
	DBCosmosDBOperationTypeDelete = DBCosmosDBOperationTypeKey.String("Delete")
	
	DBCosmosDBOperationTypeReplace = DBCosmosDBOperationTypeKey.String("Replace")
	
	DBCosmosDBOperationTypeExecute = DBCosmosDBOperationTypeKey.String("Execute")
	
	DBCosmosDBOperationTypeQuery = DBCosmosDBOperationTypeKey.String("Query")
	
	DBCosmosDBOperationTypeHead = DBCosmosDBOperationTypeKey.String("Head")
	
	DBCosmosDBOperationTypeHeadFeed = DBCosmosDBOperationTypeKey.String("HeadFeed")
	
	DBCosmosDBOperationTypeUpsert = DBCosmosDBOperationTypeKey.String("Upsert")
	
	DBCosmosDBOperationTypeBatch = DBCosmosDBOperationTypeKey.String("Batch")
	
	DBCosmosDBOperationTypeQueryPlan = DBCosmosDBOperationTypeKey.String("QueryPlan")
	
	DBCosmosDBOperationTypeExecuteJavascript = DBCosmosDBOperationTypeKey.String("ExecuteJavaScript")
)

var (
	
	DBSystemOtherSQL = DBSystemKey.String("other_sql")
	
	DBSystemMSSQL = DBSystemKey.String("mssql")
	
	DBSystemMssqlcompact = DBSystemKey.String("mssqlcompact")
	
	DBSystemMySQL = DBSystemKey.String("mysql")
	
	DBSystemOracle = DBSystemKey.String("oracle")
	
	DBSystemDB2 = DBSystemKey.String("db2")
	
	DBSystemPostgreSQL = DBSystemKey.String("postgresql")
	
	DBSystemRedshift = DBSystemKey.String("redshift")
	
	DBSystemHive = DBSystemKey.String("hive")
	
	DBSystemCloudscape = DBSystemKey.String("cloudscape")
	
	DBSystemHSQLDB = DBSystemKey.String("hsqldb")
	
	DBSystemProgress = DBSystemKey.String("progress")
	
	DBSystemMaxDB = DBSystemKey.String("maxdb")
	
	DBSystemHanaDB = DBSystemKey.String("hanadb")
	
	DBSystemIngres = DBSystemKey.String("ingres")
	
	DBSystemFirstSQL = DBSystemKey.String("firstsql")
	
	DBSystemEDB = DBSystemKey.String("edb")
	
	DBSystemCache = DBSystemKey.String("cache")
	
	DBSystemAdabas = DBSystemKey.String("adabas")
	
	DBSystemFirebird = DBSystemKey.String("firebird")
	
	DBSystemDerby = DBSystemKey.String("derby")
	
	DBSystemFilemaker = DBSystemKey.String("filemaker")
	
	DBSystemInformix = DBSystemKey.String("informix")
	
	DBSystemInstantDB = DBSystemKey.String("instantdb")
	
	DBSystemInterbase = DBSystemKey.String("interbase")
	
	DBSystemMariaDB = DBSystemKey.String("mariadb")
	
	DBSystemNetezza = DBSystemKey.String("netezza")
	
	DBSystemPervasive = DBSystemKey.String("pervasive")
	
	DBSystemPointbase = DBSystemKey.String("pointbase")
	
	DBSystemSqlite = DBSystemKey.String("sqlite")
	
	DBSystemSybase = DBSystemKey.String("sybase")
	
	DBSystemTeradata = DBSystemKey.String("teradata")
	
	DBSystemVertica = DBSystemKey.String("vertica")
	
	DBSystemH2 = DBSystemKey.String("h2")
	
	DBSystemColdfusion = DBSystemKey.String("coldfusion")
	
	DBSystemCassandra = DBSystemKey.String("cassandra")
	
	DBSystemHBase = DBSystemKey.String("hbase")
	
	DBSystemMongoDB = DBSystemKey.String("mongodb")
	
	DBSystemRedis = DBSystemKey.String("redis")
	
	DBSystemCouchbase = DBSystemKey.String("couchbase")
	
	DBSystemCouchDB = DBSystemKey.String("couchdb")
	
	DBSystemCosmosDB = DBSystemKey.String("cosmosdb")
	
	DBSystemDynamoDB = DBSystemKey.String("dynamodb")
	
	DBSystemNeo4j = DBSystemKey.String("neo4j")
	
	DBSystemGeode = DBSystemKey.String("geode")
	
	DBSystemElasticsearch = DBSystemKey.String("elasticsearch")
	
	DBSystemMemcached = DBSystemKey.String("memcached")
	
	DBSystemCockroachdb = DBSystemKey.String("cockroachdb")
	
	DBSystemOpensearch = DBSystemKey.String("opensearch")
	
	DBSystemClickhouse = DBSystemKey.String("clickhouse")
	
	DBSystemSpanner = DBSystemKey.String("spanner")
	
	DBSystemTrino = DBSystemKey.String("trino")
)




func DBCassandraCoordinatorDC(val string) attribute.KeyValue {
	return DBCassandraCoordinatorDCKey.String(val)
}




func DBCassandraCoordinatorID(val string) attribute.KeyValue {
	return DBCassandraCoordinatorIDKey.String(val)
}




func DBCassandraIdempotence(val bool) attribute.KeyValue {
	return DBCassandraIdempotenceKey.Bool(val)
}




func DBCassandraPageSize(val int) attribute.KeyValue {
	return DBCassandraPageSizeKey.Int(val)
}





func DBCassandraSpeculativeExecutionCount(val int) attribute.KeyValue {
	return DBCassandraSpeculativeExecutionCountKey.Int(val)
}





func DBCassandraTable(val string) attribute.KeyValue {
	return DBCassandraTableKey.String(val)
}





func DBConnectionString(val string) attribute.KeyValue {
	return DBConnectionStringKey.String(val)
}




func DBCosmosDBClientID(val string) attribute.KeyValue {
	return DBCosmosDBClientIDKey.String(val)
}




func DBCosmosDBContainer(val string) attribute.KeyValue {
	return DBCosmosDBContainerKey.String(val)
}




func DBCosmosDBRequestCharge(val float64) attribute.KeyValue {
	return DBCosmosDBRequestChargeKey.Float64(val)
}




func DBCosmosDBRequestContentLength(val int) attribute.KeyValue {
	return DBCosmosDBRequestContentLengthKey.Int(val)
}




func DBCosmosDBStatusCode(val int) attribute.KeyValue {
	return DBCosmosDBStatusCodeKey.Int(val)
}




func DBCosmosDBSubStatusCode(val int) attribute.KeyValue {
	return DBCosmosDBSubStatusCodeKey.Int(val)
}




func DBElasticsearchClusterName(val string) attribute.KeyValue {
	return DBElasticsearchClusterNameKey.String(val)
}





func DBElasticsearchNodeName(val string) attribute.KeyValue {
	return DBElasticsearchNodeNameKey.String(val)
}









func DBInstanceID(val string) attribute.KeyValue {
	return DBInstanceIDKey.String(val)
}






func DBJDBCDriverClassname(val string) attribute.KeyValue {
	return DBJDBCDriverClassnameKey.String(val)
}




func DBMongoDBCollection(val string) attribute.KeyValue {
	return DBMongoDBCollectionKey.String(val)
}






func DBMSSQLInstanceName(val string) attribute.KeyValue {
	return DBMSSQLInstanceNameKey.String(val)
}





func DBName(val string) attribute.KeyValue {
	return DBNameKey.String(val)
}






func DBOperation(val string) attribute.KeyValue {
	return DBOperationKey.String(val)
}






func DBRedisDBIndex(val int) attribute.KeyValue {
	return DBRedisDBIndexKey.Int(val)
}




func DBSQLTable(val string) attribute.KeyValue {
	return DBSQLTableKey.String(val)
}




func DBStatement(val string) attribute.KeyValue {
	return DBStatementKey.String(val)
}



func DBUser(val string) attribute.KeyValue {
	return DBUserKey.String(val)
}


const (
	
	
	
	
	
	
	
	HTTPFlavorKey = attribute.Key("http.flavor")

	
	
	
	
	
	
	
	
	HTTPMethodKey = attribute.Key("http.method")

	
	
	
	
	
	
	
	
	HTTPRequestContentLengthKey = attribute.Key("http.request_content_length")

	
	
	
	
	
	
	
	
	HTTPResponseContentLengthKey = attribute.Key("http.response_content_length")

	
	
	
	
	
	
	
	
	HTTPSchemeKey = attribute.Key("http.scheme")

	
	
	
	
	
	
	
	
	HTTPStatusCodeKey = attribute.Key("http.status_code")

	
	
	
	
	
	
	
	
	HTTPTargetKey = attribute.Key("http.target")

	
	
	
	
	
	
	
	
	HTTPURLKey = attribute.Key("http.url")

	
	
	
	
	
	
	
	
	
	
	HTTPUserAgentKey = attribute.Key("http.user_agent")
)

var (
	
	
	
	HTTPFlavorHTTP10 = HTTPFlavorKey.String("1.0")
	
	
	
	HTTPFlavorHTTP11 = HTTPFlavorKey.String("1.1")
	
	
	
	HTTPFlavorHTTP20 = HTTPFlavorKey.String("2.0")
	
	
	
	HTTPFlavorHTTP30 = HTTPFlavorKey.String("3.0")
	
	
	
	HTTPFlavorSPDY = HTTPFlavorKey.String("SPDY")
	
	
	
	HTTPFlavorQUIC = HTTPFlavorKey.String("QUIC")
)





func HTTPMethod(val string) attribute.KeyValue {
	return HTTPMethodKey.String(val)
}





func HTTPRequestContentLength(val int) attribute.KeyValue {
	return HTTPRequestContentLengthKey.Int(val)
}





func HTTPResponseContentLength(val int) attribute.KeyValue {
	return HTTPResponseContentLengthKey.Int(val)
}





func HTTPScheme(val string) attribute.KeyValue {
	return HTTPSchemeKey.String(val)
}





func HTTPStatusCode(val int) attribute.KeyValue {
	return HTTPStatusCodeKey.Int(val)
}





func HTTPTarget(val string) attribute.KeyValue {
	return HTTPTargetKey.String(val)
}





func HTTPURL(val string) attribute.KeyValue {
	return HTTPURLKey.String(val)
}





func HTTPUserAgent(val string) attribute.KeyValue {
	return HTTPUserAgentKey.String(val)
}


const (
	
	
	
	
	
	
	
	
	NetHostNameKey = attribute.Key("net.host.name")

	
	
	
	
	
	
	
	
	NetHostPortKey = attribute.Key("net.host.port")

	
	
	
	
	
	
	
	
	
	NetPeerNameKey = attribute.Key("net.peer.name")

	
	
	
	
	
	
	
	
	
	NetPeerPortKey = attribute.Key("net.peer.port")

	
	
	
	
	
	
	
	
	NetProtocolNameKey = attribute.Key("net.protocol.name")

	
	
	
	
	
	
	
	
	NetProtocolVersionKey = attribute.Key("net.protocol.version")

	
	
	
	
	
	
	
	NetSockFamilyKey = attribute.Key("net.sock.family")

	
	
	
	
	
	
	
	
	NetSockHostAddrKey = attribute.Key("net.sock.host.addr")

	
	
	
	
	
	
	
	
	NetSockHostPortKey = attribute.Key("net.sock.host.port")

	
	
	
	
	
	
	
	
	NetSockPeerAddrKey = attribute.Key("net.sock.peer.addr")

	
	
	
	
	
	
	
	
	NetSockPeerNameKey = attribute.Key("net.sock.peer.name")

	
	
	
	
	
	
	
	
	NetSockPeerPortKey = attribute.Key("net.sock.peer.port")

	
	
	
	
	
	
	
	NetTransportKey = attribute.Key("net.transport")
)

var (
	
	
	
	NetSockFamilyInet = NetSockFamilyKey.String("inet")
	
	
	
	NetSockFamilyInet6 = NetSockFamilyKey.String("inet6")
	
	
	
	NetSockFamilyUnix = NetSockFamilyKey.String("unix")
)

var (
	
	
	
	NetTransportTCP = NetTransportKey.String("ip_tcp")
	
	
	
	NetTransportUDP = NetTransportKey.String("ip_udp")
	
	
	
	NetTransportPipe = NetTransportKey.String("pipe")
	
	
	
	NetTransportInProc = NetTransportKey.String("inproc")
	
	
	
	NetTransportOther = NetTransportKey.String("other")
)





func NetHostName(val string) attribute.KeyValue {
	return NetHostNameKey.String(val)
}





func NetHostPort(val int) attribute.KeyValue {
	return NetHostPortKey.Int(val)
}






func NetPeerName(val string) attribute.KeyValue {
	return NetPeerNameKey.String(val)
}






func NetPeerPort(val int) attribute.KeyValue {
	return NetPeerPortKey.Int(val)
}





func NetProtocolName(val string) attribute.KeyValue {
	return NetProtocolNameKey.String(val)
}





func NetProtocolVersion(val string) attribute.KeyValue {
	return NetProtocolVersionKey.String(val)
}





func NetSockHostAddr(val string) attribute.KeyValue {
	return NetSockHostAddrKey.String(val)
}





func NetSockHostPort(val int) attribute.KeyValue {
	return NetSockHostPortKey.Int(val)
}





func NetSockPeerAddr(val string) attribute.KeyValue {
	return NetSockPeerAddrKey.String(val)
}





func NetSockPeerName(val string) attribute.KeyValue {
	return NetSockPeerNameKey.String(val)
}





func NetSockPeerPort(val int) attribute.KeyValue {
	return NetSockPeerPortKey.Int(val)
}









const (
	
	
	
	
	
	
	
	
	
	
	
	
	
	DestinationAddressKey = attribute.Key("destination.address")

	
	
	
	
	
	
	
	
	DestinationPortKey = attribute.Key("destination.port")
)





func DestinationAddress(val string) attribute.KeyValue {
	return DestinationAddressKey.String(val)
}




func DestinationPort(val int) attribute.KeyValue {
	return DestinationPortKey.Int(val)
}


const (
	
	
	
	
	
	
	
	
	DiskIoDirectionKey = attribute.Key("disk.io.direction")
)

var (
	
	DiskIoDirectionRead = DiskIoDirectionKey.String("read")
	
	DiskIoDirectionWrite = DiskIoDirectionKey.String("write")
)


const (
	
	
	
	
	
	
	
	
	
	
	
	
	
	
	
	
	
	
	
	
	
	
	
	
	
	
	
	
	
	
	
	ErrorTypeKey = attribute.Key("error.type")
)

var (
	
	ErrorTypeOther = ErrorTypeKey.String("_OTHER")
)



const (
	
	
	
	
	
	
	
	
	
	
	
	
	
	
	
	
	
	
	
	
	
	
	
	
	
	
	
	
	
	
	
	ExceptionEscapedKey = attribute.Key("exception.escaped")

	
	
	
	
	
	
	
	
	
	ExceptionMessageKey = attribute.Key("exception.message")

	
	
	
	
	
	
	
	
	
	
	
	
	
	ExceptionStacktraceKey = attribute.Key("exception.stacktrace")

	
	
	
	
	
	
	
	
	
	
	ExceptionTypeKey = attribute.Key("exception.type")
)





func ExceptionEscaped(val bool) attribute.KeyValue {
	return ExceptionEscapedKey.Bool(val)
}




func ExceptionMessage(val string) attribute.KeyValue {
	return ExceptionMessageKey.String(val)
}





func ExceptionStacktrace(val string) attribute.KeyValue {
	return ExceptionStacktraceKey.String(val)
}






func ExceptionType(val string) attribute.KeyValue {
	return ExceptionTypeKey.String(val)
}


const (
	
	
	
	
	
	
	
	
	
	
	
	
	
	HTTPRequestBodySizeKey = attribute.Key("http.request.body.size")

	
	
	
	
	
	
	
	
	
	
	
	
	
	
	
	
	
	
	
	
	
	
	
	
	
	
	
	
	
	
	
	
	
	HTTPRequestMethodKey = attribute.Key("http.request.method")

	
	
	
	
	
	
	
	
	HTTPRequestMethodOriginalKey = attribute.Key("http.request.method_original")

	
	
	
	
	
	
	
	
	
	
	
	
	
	HTTPRequestResendCountKey = attribute.Key("http.request.resend_count")

	
	
	
	
	
	
	
	
	
	
	
	
	
	HTTPResponseBodySizeKey = attribute.Key("http.response.body.size")

	
	
	
	
	
	
	
	
	
	HTTPResponseStatusCodeKey = attribute.Key("http.response.status_code")

	
	
	
	
	
	
	
	
	
	
	
	
	
	HTTPRouteKey = attribute.Key("http.route")
)

var (
	
	HTTPRequestMethodConnect = HTTPRequestMethodKey.String("CONNECT")
	
	HTTPRequestMethodDelete = HTTPRequestMethodKey.String("DELETE")
	
	HTTPRequestMethodGet = HTTPRequestMethodKey.String("GET")
	
	HTTPRequestMethodHead = HTTPRequestMethodKey.String("HEAD")
	
	HTTPRequestMethodOptions = HTTPRequestMethodKey.String("OPTIONS")
	
	HTTPRequestMethodPatch = HTTPRequestMethodKey.String("PATCH")
	
	HTTPRequestMethodPost = HTTPRequestMethodKey.String("POST")
	
	HTTPRequestMethodPut = HTTPRequestMethodKey.String("PUT")
	
	HTTPRequestMethodTrace = HTTPRequestMethodKey.String("TRACE")
	
	HTTPRequestMethodOther = HTTPRequestMethodKey.String("_OTHER")
)








func HTTPRequestBodySize(val int) attribute.KeyValue {
	return HTTPRequestBodySizeKey.Int(val)
}




func HTTPRequestMethodOriginal(val string) attribute.KeyValue {
	return HTTPRequestMethodOriginalKey.String(val)
}




func HTTPRequestResendCount(val int) attribute.KeyValue {
	return HTTPRequestResendCountKey.Int(val)
}








func HTTPResponseBodySize(val int) attribute.KeyValue {
	return HTTPResponseBodySizeKey.Int(val)
}




func HTTPResponseStatusCode(val int) attribute.KeyValue {
	return HTTPResponseStatusCodeKey.Int(val)
}




func HTTPRoute(val string) attribute.KeyValue {
	return HTTPRouteKey.String(val)
}



const (
	
	
	
	
	
	
	
	
	
	
	
	
	
	
	MessagingBatchMessageCountKey = attribute.Key("messaging.batch.message_count")

	
	
	
	
	
	
	
	
	MessagingClientIDKey = attribute.Key("messaging.client_id")

	
	
	
	
	
	
	
	
	MessagingDestinationAnonymousKey = attribute.Key("messaging.destination.anonymous")

	
	
	
	
	
	
	
	
	
	
	
	
	MessagingDestinationNameKey = attribute.Key("messaging.destination.name")

	
	
	
	
	
	
	
	
	
	
	
	
	
	MessagingDestinationTemplateKey = attribute.Key("messaging.destination.template")

	
	
	
	
	
	
	
	
	MessagingDestinationTemporaryKey = attribute.Key("messaging.destination.temporary")

	
	
	
	
	
	
	
	
	MessagingDestinationPublishAnonymousKey = attribute.Key("messaging.destination_publish.anonymous")

	
	
	
	
	
	
	
	
	
	
	
	
	
	MessagingDestinationPublishNameKey = attribute.Key("messaging.destination_publish.name")

	
	
	
	
	
	
	
	
	
	MessagingGCPPubsubMessageOrderingKeyKey = attribute.Key("messaging.gcp_pubsub.message.ordering_key")

	
	
	
	
	
	
	
	
	
	MessagingKafkaConsumerGroupKey = attribute.Key("messaging.kafka.consumer.group")

	
	
	
	
	
	
	
	
	MessagingKafkaDestinationPartitionKey = attribute.Key("messaging.kafka.destination.partition")

	
	
	
	
	
	
	
	
	
	
	
	
	
	
	MessagingKafkaMessageKeyKey = attribute.Key("messaging.kafka.message.key")

	
	
	
	
	
	
	
	
	MessagingKafkaMessageOffsetKey = attribute.Key("messaging.kafka.message.offset")

	
	
	
	
	
	
	
	MessagingKafkaMessageTombstoneKey = attribute.Key("messaging.kafka.message.tombstone")

	
	
	
	
	
	
	
	
	
	
	
	MessagingMessageBodySizeKey = attribute.Key("messaging.message.body.size")

	
	
	
	
	
	
	
	
	
	MessagingMessageConversationIDKey = attribute.Key("messaging.message.conversation_id")

	
	
	
	
	
	
	
	
	
	
	
	MessagingMessageEnvelopeSizeKey = attribute.Key("messaging.message.envelope.size")

	
	
	
	
	
	
	
	
	
	MessagingMessageIDKey = attribute.Key("messaging.message.id")

	
	
	
	
	
	
	
	
	MessagingOperationKey = attribute.Key("messaging.operation")

	
	
	
	
	
	
	
	
	MessagingRabbitmqDestinationRoutingKeyKey = attribute.Key("messaging.rabbitmq.destination.routing_key")

	
	
	
	
	
	
	
	
	
	MessagingRocketmqClientGroupKey = attribute.Key("messaging.rocketmq.client_group")

	
	
	
	
	
	
	
	
	MessagingRocketmqConsumptionModelKey = attribute.Key("messaging.rocketmq.consumption_model")

	
	
	
	
	
	
	
	
	
	MessagingRocketmqMessageDelayTimeLevelKey = attribute.Key("messaging.rocketmq.message.delay_time_level")

	
	
	
	
	
	
	
	
	
	MessagingRocketmqMessageDeliveryTimestampKey = attribute.Key("messaging.rocketmq.message.delivery_timestamp")

	
	
	
	
	
	
	
	
	
	
	MessagingRocketmqMessageGroupKey = attribute.Key("messaging.rocketmq.message.group")

	
	
	
	
	
	
	
	
	MessagingRocketmqMessageKeysKey = attribute.Key("messaging.rocketmq.message.keys")

	
	
	
	
	
	
	
	
	MessagingRocketmqMessageTagKey = attribute.Key("messaging.rocketmq.message.tag")

	
	
	
	
	
	
	
	MessagingRocketmqMessageTypeKey = attribute.Key("messaging.rocketmq.message.type")

	
	
	
	
	
	
	
	
	
	MessagingRocketmqNamespaceKey = attribute.Key("messaging.rocketmq.namespace")

	
	
	
	
	
	
	
	
	MessagingSystemKey = attribute.Key("messaging.system")
)

var (
	
	MessagingOperationPublish = MessagingOperationKey.String("publish")
	
	MessagingOperationCreate = MessagingOperationKey.String("create")
	
	MessagingOperationReceive = MessagingOperationKey.String("receive")
	
	MessagingOperationDeliver = MessagingOperationKey.String("deliver")
)

var (
	
	MessagingRocketmqConsumptionModelClustering = MessagingRocketmqConsumptionModelKey.String("clustering")
	
	MessagingRocketmqConsumptionModelBroadcasting = MessagingRocketmqConsumptionModelKey.String("broadcasting")
)

var (
	
	MessagingRocketmqMessageTypeNormal = MessagingRocketmqMessageTypeKey.String("normal")
	
	MessagingRocketmqMessageTypeFifo = MessagingRocketmqMessageTypeKey.String("fifo")
	
	MessagingRocketmqMessageTypeDelay = MessagingRocketmqMessageTypeKey.String("delay")
	
	MessagingRocketmqMessageTypeTransaction = MessagingRocketmqMessageTypeKey.String("transaction")
)

var (
	
	MessagingSystemActivemq = MessagingSystemKey.String("activemq")
	
	MessagingSystemAWSSqs = MessagingSystemKey.String("aws_sqs")
	
	MessagingSystemAzureEventgrid = MessagingSystemKey.String("azure_eventgrid")
	
	MessagingSystemAzureEventhubs = MessagingSystemKey.String("azure_eventhubs")
	
	MessagingSystemAzureServicebus = MessagingSystemKey.String("azure_servicebus")
	
	MessagingSystemGCPPubsub = MessagingSystemKey.String("gcp_pubsub")
	
	MessagingSystemJms = MessagingSystemKey.String("jms")
	
	MessagingSystemKafka = MessagingSystemKey.String("kafka")
	
	MessagingSystemRabbitmq = MessagingSystemKey.String("rabbitmq")
	
	MessagingSystemRocketmq = MessagingSystemKey.String("rocketmq")
)





func MessagingBatchMessageCount(val int) attribute.KeyValue {
	return MessagingBatchMessageCountKey.Int(val)
}




func MessagingClientID(val string) attribute.KeyValue {
	return MessagingClientIDKey.String(val)
}





func MessagingDestinationAnonymous(val bool) attribute.KeyValue {
	return MessagingDestinationAnonymousKey.Bool(val)
}




func MessagingDestinationName(val string) attribute.KeyValue {
	return MessagingDestinationNameKey.String(val)
}




func MessagingDestinationTemplate(val string) attribute.KeyValue {
	return MessagingDestinationTemplateKey.String(val)
}





func MessagingDestinationTemporary(val bool) attribute.KeyValue {
	return MessagingDestinationTemporaryKey.Bool(val)
}





func MessagingDestinationPublishAnonymous(val bool) attribute.KeyValue {
	return MessagingDestinationPublishAnonymousKey.Bool(val)
}




func MessagingDestinationPublishName(val string) attribute.KeyValue {
	return MessagingDestinationPublishNameKey.String(val)
}





func MessagingGCPPubsubMessageOrderingKey(val string) attribute.KeyValue {
	return MessagingGCPPubsubMessageOrderingKeyKey.String(val)
}





func MessagingKafkaConsumerGroup(val string) attribute.KeyValue {
	return MessagingKafkaConsumerGroupKey.String(val)
}




func MessagingKafkaDestinationPartition(val int) attribute.KeyValue {
	return MessagingKafkaDestinationPartitionKey.Int(val)
}







func MessagingKafkaMessageKey(val string) attribute.KeyValue {
	return MessagingKafkaMessageKeyKey.String(val)
}




func MessagingKafkaMessageOffset(val int) attribute.KeyValue {
	return MessagingKafkaMessageOffsetKey.Int(val)
}




func MessagingKafkaMessageTombstone(val bool) attribute.KeyValue {
	return MessagingKafkaMessageTombstoneKey.Bool(val)
}




func MessagingMessageBodySize(val int) attribute.KeyValue {
	return MessagingMessageBodySizeKey.Int(val)
}





func MessagingMessageConversationID(val string) attribute.KeyValue {
	return MessagingMessageConversationIDKey.String(val)
}




func MessagingMessageEnvelopeSize(val int) attribute.KeyValue {
	return MessagingMessageEnvelopeSizeKey.Int(val)
}





func MessagingMessageID(val string) attribute.KeyValue {
	return MessagingMessageIDKey.String(val)
}




func MessagingRabbitmqDestinationRoutingKey(val string) attribute.KeyValue {
	return MessagingRabbitmqDestinationRoutingKeyKey.String(val)
}





func MessagingRocketmqClientGroup(val string) attribute.KeyValue {
	return MessagingRocketmqClientGroupKey.String(val)
}





func MessagingRocketmqMessageDelayTimeLevel(val int) attribute.KeyValue {
	return MessagingRocketmqMessageDelayTimeLevelKey.Int(val)
}





func MessagingRocketmqMessageDeliveryTimestamp(val int) attribute.KeyValue {
	return MessagingRocketmqMessageDeliveryTimestampKey.Int(val)
}






func MessagingRocketmqMessageGroup(val string) attribute.KeyValue {
	return MessagingRocketmqMessageGroupKey.String(val)
}




func MessagingRocketmqMessageKeys(val ...string) attribute.KeyValue {
	return MessagingRocketmqMessageKeysKey.StringSlice(val)
}




func MessagingRocketmqMessageTag(val string) attribute.KeyValue {
	return MessagingRocketmqMessageTagKey.String(val)
}





func MessagingRocketmqNamespace(val string) attribute.KeyValue {
	return MessagingRocketmqNamespaceKey.String(val)
}


const (
	
	
	
	
	
	
	
	
	
	NetworkCarrierIccKey = attribute.Key("network.carrier.icc")

	
	
	
	
	
	
	
	
	NetworkCarrierMccKey = attribute.Key("network.carrier.mcc")

	
	
	
	
	
	
	
	
	NetworkCarrierMncKey = attribute.Key("network.carrier.mnc")

	
	
	
	
	
	
	
	
	NetworkCarrierNameKey = attribute.Key("network.carrier.name")

	
	
	
	
	
	
	
	
	
	
	NetworkConnectionSubtypeKey = attribute.Key("network.connection.subtype")

	
	
	
	
	
	
	
	
	NetworkConnectionTypeKey = attribute.Key("network.connection.type")

	
	
	
	
	
	
	
	
	NetworkIoDirectionKey = attribute.Key("network.io.direction")

	
	
	
	
	
	
	
	
	
	NetworkLocalAddressKey = attribute.Key("network.local.address")

	
	
	
	
	
	
	
	
	NetworkLocalPortKey = attribute.Key("network.local.port")

	
	
	
	
	
	
	
	
	
	NetworkPeerAddressKey = attribute.Key("network.peer.address")

	
	
	
	
	
	
	
	
	NetworkPeerPortKey = attribute.Key("network.peer.port")

	
	
	
	
	
	
	
	
	
	
	NetworkProtocolNameKey = attribute.Key("network.protocol.name")

	
	
	
	
	
	
	
	
	
	
	
	
	NetworkProtocolVersionKey = attribute.Key("network.protocol.version")

	
	
	
	
	
	
	
	
	
	
	
	
	
	
	
	
	NetworkTransportKey = attribute.Key("network.transport")

	
	
	
	
	
	
	
	
	
	NetworkTypeKey = attribute.Key("network.type")
)

var (
	
	NetworkConnectionSubtypeGprs = NetworkConnectionSubtypeKey.String("gprs")
	
	NetworkConnectionSubtypeEdge = NetworkConnectionSubtypeKey.String("edge")
	
	NetworkConnectionSubtypeUmts = NetworkConnectionSubtypeKey.String("umts")
	
	NetworkConnectionSubtypeCdma = NetworkConnectionSubtypeKey.String("cdma")
	
	NetworkConnectionSubtypeEvdo0 = NetworkConnectionSubtypeKey.String("evdo_0")
	
	NetworkConnectionSubtypeEvdoA = NetworkConnectionSubtypeKey.String("evdo_a")
	
	NetworkConnectionSubtypeCdma20001xrtt = NetworkConnectionSubtypeKey.String("cdma2000_1xrtt")
	
	NetworkConnectionSubtypeHsdpa = NetworkConnectionSubtypeKey.String("hsdpa")
	
	NetworkConnectionSubtypeHsupa = NetworkConnectionSubtypeKey.String("hsupa")
	
	NetworkConnectionSubtypeHspa = NetworkConnectionSubtypeKey.String("hspa")
	
	NetworkConnectionSubtypeIden = NetworkConnectionSubtypeKey.String("iden")
	
	NetworkConnectionSubtypeEvdoB = NetworkConnectionSubtypeKey.String("evdo_b")
	
	NetworkConnectionSubtypeLte = NetworkConnectionSubtypeKey.String("lte")
	
	NetworkConnectionSubtypeEhrpd = NetworkConnectionSubtypeKey.String("ehrpd")
	
	NetworkConnectionSubtypeHspap = NetworkConnectionSubtypeKey.String("hspap")
	
	NetworkConnectionSubtypeGsm = NetworkConnectionSubtypeKey.String("gsm")
	
	NetworkConnectionSubtypeTdScdma = NetworkConnectionSubtypeKey.String("td_scdma")
	
	NetworkConnectionSubtypeIwlan = NetworkConnectionSubtypeKey.String("iwlan")
	
	NetworkConnectionSubtypeNr = NetworkConnectionSubtypeKey.String("nr")
	
	NetworkConnectionSubtypeNrnsa = NetworkConnectionSubtypeKey.String("nrnsa")
	
	NetworkConnectionSubtypeLteCa = NetworkConnectionSubtypeKey.String("lte_ca")
)

var (
	
	NetworkConnectionTypeWifi = NetworkConnectionTypeKey.String("wifi")
	
	NetworkConnectionTypeWired = NetworkConnectionTypeKey.String("wired")
	
	NetworkConnectionTypeCell = NetworkConnectionTypeKey.String("cell")
	
	NetworkConnectionTypeUnavailable = NetworkConnectionTypeKey.String("unavailable")
	
	NetworkConnectionTypeUnknown = NetworkConnectionTypeKey.String("unknown")
)

var (
	
	NetworkIoDirectionTransmit = NetworkIoDirectionKey.String("transmit")
	
	NetworkIoDirectionReceive = NetworkIoDirectionKey.String("receive")
)

var (
	
	NetworkTransportTCP = NetworkTransportKey.String("tcp")
	
	NetworkTransportUDP = NetworkTransportKey.String("udp")
	
	NetworkTransportPipe = NetworkTransportKey.String("pipe")
	
	NetworkTransportUnix = NetworkTransportKey.String("unix")
)

var (
	
	NetworkTypeIpv4 = NetworkTypeKey.String("ipv4")
	
	NetworkTypeIpv6 = NetworkTypeKey.String("ipv6")
)




func NetworkCarrierIcc(val string) attribute.KeyValue {
	return NetworkCarrierIccKey.String(val)
}




func NetworkCarrierMcc(val string) attribute.KeyValue {
	return NetworkCarrierMccKey.String(val)
}




func NetworkCarrierMnc(val string) attribute.KeyValue {
	return NetworkCarrierMncKey.String(val)
}




func NetworkCarrierName(val string) attribute.KeyValue {
	return NetworkCarrierNameKey.String(val)
}




func NetworkLocalAddress(val string) attribute.KeyValue {
	return NetworkLocalAddressKey.String(val)
}




func NetworkLocalPort(val int) attribute.KeyValue {
	return NetworkLocalPortKey.Int(val)
}




func NetworkPeerAddress(val string) attribute.KeyValue {
	return NetworkPeerAddressKey.String(val)
}




func NetworkPeerPort(val int) attribute.KeyValue {
	return NetworkPeerPortKey.Int(val)
}





func NetworkProtocolName(val string) attribute.KeyValue {
	return NetworkProtocolNameKey.String(val)
}




func NetworkProtocolVersion(val string) attribute.KeyValue {
	return NetworkProtocolVersionKey.String(val)
}


const (
	
	
	
	
	
	
	
	
	RPCConnectRPCErrorCodeKey = attribute.Key("rpc.connect_rpc.error_code")

	
	
	
	
	
	
	
	
	
	RPCGRPCStatusCodeKey = attribute.Key("rpc.grpc.status_code")

	
	
	
	
	
	
	
	
	RPCJsonrpcErrorCodeKey = attribute.Key("rpc.jsonrpc.error_code")

	
	
	
	
	
	
	
	
	RPCJsonrpcErrorMessageKey = attribute.Key("rpc.jsonrpc.error_message")

	
	
	
	
	
	
	
	
	
	
	
	RPCJsonrpcRequestIDKey = attribute.Key("rpc.jsonrpc.request_id")

	
	
	
	
	
	
	
	
	
	RPCJsonrpcVersionKey = attribute.Key("rpc.jsonrpc.version")

	
	
	
	
	
	
	
	
	
	
	
	
	
	RPCMethodKey = attribute.Key("rpc.method")

	
	
	
	
	
	
	
	
	
	
	
	
	
	
	RPCServiceKey = attribute.Key("rpc.service")

	
	
	
	
	
	
	
	RPCSystemKey = attribute.Key("rpc.system")
)

var (
	
	RPCConnectRPCErrorCodeCancelled = RPCConnectRPCErrorCodeKey.String("cancelled")
	
	RPCConnectRPCErrorCodeUnknown = RPCConnectRPCErrorCodeKey.String("unknown")
	
	RPCConnectRPCErrorCodeInvalidArgument = RPCConnectRPCErrorCodeKey.String("invalid_argument")
	
	RPCConnectRPCErrorCodeDeadlineExceeded = RPCConnectRPCErrorCodeKey.String("deadline_exceeded")
	
	RPCConnectRPCErrorCodeNotFound = RPCConnectRPCErrorCodeKey.String("not_found")
	
	RPCConnectRPCErrorCodeAlreadyExists = RPCConnectRPCErrorCodeKey.String("already_exists")
	
	RPCConnectRPCErrorCodePermissionDenied = RPCConnectRPCErrorCodeKey.String("permission_denied")
	
	RPCConnectRPCErrorCodeResourceExhausted = RPCConnectRPCErrorCodeKey.String("resource_exhausted")
	
	RPCConnectRPCErrorCodeFailedPrecondition = RPCConnectRPCErrorCodeKey.String("failed_precondition")
	
	RPCConnectRPCErrorCodeAborted = RPCConnectRPCErrorCodeKey.String("aborted")
	
	RPCConnectRPCErrorCodeOutOfRange = RPCConnectRPCErrorCodeKey.String("out_of_range")
	
	RPCConnectRPCErrorCodeUnimplemented = RPCConnectRPCErrorCodeKey.String("unimplemented")
	
	RPCConnectRPCErrorCodeInternal = RPCConnectRPCErrorCodeKey.String("internal")
	
	RPCConnectRPCErrorCodeUnavailable = RPCConnectRPCErrorCodeKey.String("unavailable")
	
	RPCConnectRPCErrorCodeDataLoss = RPCConnectRPCErrorCodeKey.String("data_loss")
	
	RPCConnectRPCErrorCodeUnauthenticated = RPCConnectRPCErrorCodeKey.String("unauthenticated")
)

var (
	
	RPCGRPCStatusCodeOk = RPCGRPCStatusCodeKey.Int(0)
	
	RPCGRPCStatusCodeCancelled = RPCGRPCStatusCodeKey.Int(1)
	
	RPCGRPCStatusCodeUnknown = RPCGRPCStatusCodeKey.Int(2)
	
	RPCGRPCStatusCodeInvalidArgument = RPCGRPCStatusCodeKey.Int(3)
	
	RPCGRPCStatusCodeDeadlineExceeded = RPCGRPCStatusCodeKey.Int(4)
	
	RPCGRPCStatusCodeNotFound = RPCGRPCStatusCodeKey.Int(5)
	
	RPCGRPCStatusCodeAlreadyExists = RPCGRPCStatusCodeKey.Int(6)
	
	RPCGRPCStatusCodePermissionDenied = RPCGRPCStatusCodeKey.Int(7)
	
	RPCGRPCStatusCodeResourceExhausted = RPCGRPCStatusCodeKey.Int(8)
	
	RPCGRPCStatusCodeFailedPrecondition = RPCGRPCStatusCodeKey.Int(9)
	
	RPCGRPCStatusCodeAborted = RPCGRPCStatusCodeKey.Int(10)
	
	RPCGRPCStatusCodeOutOfRange = RPCGRPCStatusCodeKey.Int(11)
	
	RPCGRPCStatusCodeUnimplemented = RPCGRPCStatusCodeKey.Int(12)
	
	RPCGRPCStatusCodeInternal = RPCGRPCStatusCodeKey.Int(13)
	
	RPCGRPCStatusCodeUnavailable = RPCGRPCStatusCodeKey.Int(14)
	
	RPCGRPCStatusCodeDataLoss = RPCGRPCStatusCodeKey.Int(15)
	
	RPCGRPCStatusCodeUnauthenticated = RPCGRPCStatusCodeKey.Int(16)
)

var (
	
	RPCSystemGRPC = RPCSystemKey.String("grpc")
	
	RPCSystemJavaRmi = RPCSystemKey.String("java_rmi")
	
	RPCSystemDotnetWcf = RPCSystemKey.String("dotnet_wcf")
	
	RPCSystemApacheDubbo = RPCSystemKey.String("apache_dubbo")
	
	RPCSystemConnectRPC = RPCSystemKey.String("connect_rpc")
)




func RPCJsonrpcErrorCode(val int) attribute.KeyValue {
	return RPCJsonrpcErrorCodeKey.Int(val)
}




func RPCJsonrpcErrorMessage(val string) attribute.KeyValue {
	return RPCJsonrpcErrorMessageKey.String(val)
}







func RPCJsonrpcRequestID(val string) attribute.KeyValue {
	return RPCJsonrpcRequestIDKey.String(val)
}





func RPCJsonrpcVersion(val string) attribute.KeyValue {
	return RPCJsonrpcVersionKey.String(val)
}




func RPCMethod(val string) attribute.KeyValue {
	return RPCMethodKey.String(val)
}




func RPCService(val string) attribute.KeyValue {
	return RPCServiceKey.String(val)
}









const (
	
	
	
	
	
	
	
	
	
	
	
	
	ServerAddressKey = attribute.Key("server.address")

	
	
	
	
	
	
	
	
	
	
	ServerPortKey = attribute.Key("server.port")
)





func ServerAddress(val string) attribute.KeyValue {
	return ServerAddressKey.String(val)
}



func ServerPort(val int) attribute.KeyValue {
	return ServerPortKey.Int(val)
}









const (
	
	
	
	
	
	
	
	
	
	
	
	
	
	SourceAddressKey = attribute.Key("source.address")

	
	
	
	
	
	
	
	SourcePortKey = attribute.Key("source.port")
)





func SourceAddress(val string) attribute.KeyValue {
	return SourceAddressKey.String(val)
}



func SourcePort(val int) attribute.KeyValue {
	return SourcePortKey.Int(val)
}


const (
	
	
	
	
	
	
	
	
	
	
	
	
	
	TLSCipherKey = attribute.Key("tls.cipher")

	
	
	
	
	
	
	
	
	
	
	TLSClientCertificateKey = attribute.Key("tls.client.certificate")

	
	
	
	
	
	
	
	
	
	
	
	TLSClientCertificateChainKey = attribute.Key("tls.client.certificate_chain")

	
	
	
	
	
	
	
	
	
	
	TLSClientHashMd5Key = attribute.Key("tls.client.hash.md5")

	
	
	
	
	
	
	
	
	
	
	TLSClientHashSha1Key = attribute.Key("tls.client.hash.sha1")

	
	
	
	
	
	
	
	
	
	
	
	TLSClientHashSha256Key = attribute.Key("tls.client.hash.sha256")

	
	
	
	
	
	
	
	
	
	
	
	TLSClientIssuerKey = attribute.Key("tls.client.issuer")

	
	
	
	
	
	
	
	
	TLSClientJa3Key = attribute.Key("tls.client.ja3")

	
	
	
	
	
	
	
	
	TLSClientNotAfterKey = attribute.Key("tls.client.not_after")

	
	
	
	
	
	
	
	
	TLSClientNotBeforeKey = attribute.Key("tls.client.not_before")

	
	
	
	
	
	
	
	
	
	TLSClientServerNameKey = attribute.Key("tls.client.server_name")

	
	
	
	
	
	
	
	
	
	TLSClientSubjectKey = attribute.Key("tls.client.subject")

	
	
	
	
	
	
	
	
	
	TLSClientSupportedCiphersKey = attribute.Key("tls.client.supported_ciphers")

	
	
	
	
	
	
	
	
	TLSCurveKey = attribute.Key("tls.curve")

	
	
	
	
	
	
	
	
	
	TLSEstablishedKey = attribute.Key("tls.established")

	
	
	
	
	
	
	
	
	
	
	TLSNextProtocolKey = attribute.Key("tls.next_protocol")

	
	
	
	
	
	
	
	
	
	TLSProtocolNameKey = attribute.Key("tls.protocol.name")

	
	
	
	
	
	
	
	
	
	
	TLSProtocolVersionKey = attribute.Key("tls.protocol.version")

	
	
	
	
	
	
	
	
	TLSResumedKey = attribute.Key("tls.resumed")

	
	
	
	
	
	
	
	
	
	
	TLSServerCertificateKey = attribute.Key("tls.server.certificate")

	
	
	
	
	
	
	
	
	
	
	
	TLSServerCertificateChainKey = attribute.Key("tls.server.certificate_chain")

	
	
	
	
	
	
	
	
	
	
	TLSServerHashMd5Key = attribute.Key("tls.server.hash.md5")

	
	
	
	
	
	
	
	
	
	
	TLSServerHashSha1Key = attribute.Key("tls.server.hash.sha1")

	
	
	
	
	
	
	
	
	
	
	
	TLSServerHashSha256Key = attribute.Key("tls.server.hash.sha256")

	
	
	
	
	
	
	
	
	
	
	
	TLSServerIssuerKey = attribute.Key("tls.server.issuer")

	
	
	
	
	
	
	
	
	TLSServerJa3sKey = attribute.Key("tls.server.ja3s")

	
	
	
	
	
	
	
	
	TLSServerNotAfterKey = attribute.Key("tls.server.not_after")

	
	
	
	
	
	
	
	
	TLSServerNotBeforeKey = attribute.Key("tls.server.not_before")

	
	
	
	
	
	
	
	
	
	TLSServerSubjectKey = attribute.Key("tls.server.subject")
)

var (
	
	TLSProtocolNameSsl = TLSProtocolNameKey.String("ssl")
	
	TLSProtocolNameTLS = TLSProtocolNameKey.String("tls")
)





func TLSCipher(val string) attribute.KeyValue {
	return TLSCipherKey.String(val)
}






func TLSClientCertificate(val string) attribute.KeyValue {
	return TLSClientCertificateKey.String(val)
}






func TLSClientCertificateChain(val ...string) attribute.KeyValue {
	return TLSClientCertificateChainKey.StringSlice(val)
}






func TLSClientHashMd5(val string) attribute.KeyValue {
	return TLSClientHashMd5Key.String(val)
}






func TLSClientHashSha1(val string) attribute.KeyValue {
	return TLSClientHashSha1Key.String(val)
}






func TLSClientHashSha256(val string) attribute.KeyValue {
	return TLSClientHashSha256Key.String(val)
}






func TLSClientIssuer(val string) attribute.KeyValue {
	return TLSClientIssuerKey.String(val)
}




func TLSClientJa3(val string) attribute.KeyValue {
	return TLSClientJa3Key.String(val)
}




func TLSClientNotAfter(val string) attribute.KeyValue {
	return TLSClientNotAfterKey.String(val)
}




func TLSClientNotBefore(val string) attribute.KeyValue {
	return TLSClientNotBeforeKey.String(val)
}





func TLSClientServerName(val string) attribute.KeyValue {
	return TLSClientServerNameKey.String(val)
}




func TLSClientSubject(val string) attribute.KeyValue {
	return TLSClientSubjectKey.String(val)
}




func TLSClientSupportedCiphers(val ...string) attribute.KeyValue {
	return TLSClientSupportedCiphersKey.StringSlice(val)
}




func TLSCurve(val string) attribute.KeyValue {
	return TLSCurveKey.String(val)
}





func TLSEstablished(val bool) attribute.KeyValue {
	return TLSEstablishedKey.Bool(val)
}






func TLSNextProtocol(val string) attribute.KeyValue {
	return TLSNextProtocolKey.String(val)
}






func TLSProtocolVersion(val string) attribute.KeyValue {
	return TLSProtocolVersionKey.String(val)
}




func TLSResumed(val bool) attribute.KeyValue {
	return TLSResumedKey.Bool(val)
}






func TLSServerCertificate(val string) attribute.KeyValue {
	return TLSServerCertificateKey.String(val)
}






func TLSServerCertificateChain(val ...string) attribute.KeyValue {
	return TLSServerCertificateChainKey.StringSlice(val)
}






func TLSServerHashMd5(val string) attribute.KeyValue {
	return TLSServerHashMd5Key.String(val)
}






func TLSServerHashSha1(val string) attribute.KeyValue {
	return TLSServerHashSha1Key.String(val)
}






func TLSServerHashSha256(val string) attribute.KeyValue {
	return TLSServerHashSha256Key.String(val)
}






func TLSServerIssuer(val string) attribute.KeyValue {
	return TLSServerIssuerKey.String(val)
}




func TLSServerJa3s(val string) attribute.KeyValue {
	return TLSServerJa3sKey.String(val)
}




func TLSServerNotAfter(val string) attribute.KeyValue {
	return TLSServerNotAfterKey.String(val)
}




func TLSServerNotBefore(val string) attribute.KeyValue {
	return TLSServerNotBeforeKey.String(val)
}




func TLSServerSubject(val string) attribute.KeyValue {
	return TLSServerSubjectKey.String(val)
}


const (
	
	
	
	
	
	
	
	
	URLFragmentKey = attribute.Key("url.fragment")

	
	
	
	
	
	
	
	
	
	
	
	
	
	
	
	
	
	
	
	
	URLFullKey = attribute.Key("url.full")

	
	
	
	
	
	
	
	
	URLPathKey = attribute.Key("url.path")

	
	
	
	
	
	
	
	
	
	
	URLQueryKey = attribute.Key("url.query")

	
	
	
	
	
	
	
	
	
	URLSchemeKey = attribute.Key("url.scheme")
)




func URLFragment(val string) attribute.KeyValue {
	return URLFragmentKey.String(val)
}




func URLFull(val string) attribute.KeyValue {
	return URLFullKey.String(val)
}




func URLPath(val string) attribute.KeyValue {
	return URLPathKey.String(val)
}




func URLQuery(val string) attribute.KeyValue {
	return URLQueryKey.String(val)
}





func URLScheme(val string) attribute.KeyValue {
	return URLSchemeKey.String(val)
}


const (
	
	
	
	
	
	
	
	
	
	
	
	
	UserAgentOriginalKey = attribute.Key("user_agent.original")
)






func UserAgentOriginal(val string) attribute.KeyValue {
	return UserAgentOriginalKey.String(val)
}












const (
	
	
	
	
	
	
	
	SessionIDKey = attribute.Key("session.id")

	
	
	
	
	
	
	
	
	SessionPreviousIDKey = attribute.Key("session.previous_id")
)



func SessionID(val string) attribute.KeyValue {
	return SessionIDKey.String(val)
}




func SessionPreviousID(val string) attribute.KeyValue {
	return SessionPreviousIDKey.String(val)
}

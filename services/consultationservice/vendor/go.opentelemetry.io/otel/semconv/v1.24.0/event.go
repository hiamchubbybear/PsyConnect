




package semconv 

import "go.opentelemetry.io/otel/attribute"



const (
	
	
	
	
	
	
	
	
	
	
	IosStateKey = attribute.Key("ios.state")
)

var (
	
	IosStateActive = IosStateKey.String("active")
	
	IosStateInactive = IosStateKey.String("inactive")
	
	IosStateBackground = IosStateKey.String("background")
	
	IosStateForeground = IosStateKey.String("foreground")
	
	IosStateTerminate = IosStateKey.String("terminate")
)



const (
	
	
	
	
	
	
	
	
	
	
	
	AndroidStateKey = attribute.Key("android.state")
)

var (
	
	AndroidStateCreated = AndroidStateKey.String("created")
	
	AndroidStateBackground = AndroidStateKey.String("background")
	
	AndroidStateForeground = AndroidStateKey.String("foreground")
)



const (
	
	
	
	
	
	
	
	
	FeatureFlagKeyKey = attribute.Key("feature_flag.key")

	
	
	
	
	
	
	
	
	FeatureFlagProviderNameKey = attribute.Key("feature_flag.provider_name")

	
	
	
	
	
	
	
	
	
	
	
	
	
	
	
	
	
	
	FeatureFlagVariantKey = attribute.Key("feature_flag.variant")
)




func FeatureFlagKey(val string) attribute.KeyValue {
	return FeatureFlagKeyKey.String(val)
}




func FeatureFlagProviderName(val string) attribute.KeyValue {
	return FeatureFlagProviderNameKey.String(val)
}





func FeatureFlagVariant(val string) attribute.KeyValue {
	return FeatureFlagVariantKey.String(val)
}


const (
	
	
	
	
	
	
	
	MessageCompressedSizeKey = attribute.Key("message.compressed_size")

	
	
	
	
	
	
	
	
	
	
	MessageIDKey = attribute.Key("message.id")

	
	
	
	
	
	
	
	MessageTypeKey = attribute.Key("message.type")

	
	
	
	
	
	
	
	MessageUncompressedSizeKey = attribute.Key("message.uncompressed_size")
)

var (
	
	MessageTypeSent = MessageTypeKey.String("SENT")
	
	MessageTypeReceived = MessageTypeKey.String("RECEIVED")
)




func MessageCompressedSize(val int) attribute.KeyValue {
	return MessageCompressedSizeKey.Int(val)
}





func MessageID(val int) attribute.KeyValue {
	return MessageIDKey.Int(val)
}




func MessageUncompressedSize(val int) attribute.KeyValue {
	return MessageUncompressedSizeKey.Int(val)
}

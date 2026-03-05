package protocol

import (
	"fmt"
	"io"
)

func ReadRequest(r io.Reader) (apiVersion int16, correlationID int32, clientID string, msg Message, err error) {
	d := &decoder{reader: r, remain: 4}
	size := d.readInt32()

	if err = d.err; err != nil {
		err = dontExpectEOF(err)
		return
	}

	d.remain = int(size)
	apiKey := ApiKey(d.readInt16())
	apiVersion = d.readInt16()
	correlationID = d.readInt32()
	clientID = d.readString()

	if i := int(apiKey); i < 0 || i >= len(apiTypes) {
		err = fmt.Errorf("unsupported api key: %d", i)
		return
	}

	if err = d.err; err != nil {
		err = dontExpectEOF(err)
		return
	}

	t := &apiTypes[apiKey]
	if t == nil {
		err = fmt.Errorf("unsupported api: %s", apiNames[apiKey])
		return
	}

	minVersion := t.minVersion()
	maxVersion := t.maxVersion()

	if apiVersion < minVersion || apiVersion > maxVersion {
		err = fmt.Errorf("unsupported %s version: v%d not in range v%d-v%d", apiKey, apiVersion, minVersion, maxVersion)
		return
	}

	req := &t.requests[apiVersion-minVersion]

	if req.flexible {
		
		taggedCount := int(d.readUnsignedVarInt())
		for i := 0; i < taggedCount; i++ {
			d.readUnsignedVarInt() 
			size := d.readUnsignedVarInt()

			
			d.read(int(size))
		}
	}

	msg = req.new()
	req.decode(d, valueOf(msg))
	d.discardAll()

	if err = d.err; err != nil {
		err = dontExpectEOF(err)
	}

	return
}

func WriteRequest(w io.Writer, apiVersion int16, correlationID int32, clientID string, msg Message) error {
	apiKey := msg.ApiKey()

	if i := int(apiKey); i < 0 || i >= len(apiTypes) {
		return fmt.Errorf("unsupported api key: %d", i)
	}

	t := &apiTypes[apiKey]
	if t == nil {
		return fmt.Errorf("unsupported api: %s", apiNames[apiKey])
	}

	if typedMessage, ok := msg.(OverrideTypeMessage); ok {
		typeKey := typedMessage.TypeKey()
		overrideType := overrideApiTypes[apiKey][typeKey]
		t = &overrideType
	}

	minVersion := t.minVersion()
	maxVersion := t.maxVersion()

	if apiVersion < minVersion || apiVersion > maxVersion {
		return fmt.Errorf("unsupported %s version: v%d not in range v%d-v%d", apiKey, apiVersion, minVersion, maxVersion)
	}

	r := &t.requests[apiVersion-minVersion]
	v := valueOf(msg)
	b := newPageBuffer()
	defer b.unref()

	e := &encoder{writer: b}
	e.writeInt32(0) 
	e.writeInt16(int16(apiKey))
	e.writeInt16(apiVersion)
	e.writeInt32(correlationID)

	if r.flexible {
		
		
		
		
		
		
		
		e.writeNullString(clientID)
		e.writeUnsignedVarInt(0)
	} else {
		
		
		
		e.writeString(clientID)
	}
	r.encode(e, v)
	err := e.err

	if err == nil {
		size := packUint32(uint32(b.Size()) - 4)
		b.WriteAt(size[:], 0)
		_, err = b.WriteTo(w)
	}

	return err
}

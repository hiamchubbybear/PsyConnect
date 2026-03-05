



package sse

import (
	"bytes"
	"io"
	"io/ioutil"
)

type decoder struct {
	events []Event
}

func Decode(r io.Reader) ([]Event, error) {
	var dec decoder
	return dec.decode(r)
}

func (d *decoder) dispatchEvent(event Event, data string) {
	dataLength := len(data)
	if dataLength > 0 {
		
		data = data[:dataLength-1]
		dataLength--
	}
	if dataLength == 0 && event.Event == "" {
		return
	}
	if event.Event == "" {
		event.Event = "message"
	}
	event.Data = data
	d.events = append(d.events, event)
}

func (d *decoder) decode(r io.Reader) ([]Event, error) {
	buf, err := ioutil.ReadAll(r)
	if err != nil {
		return nil, err
	}

	var currentEvent Event
	var dataBuffer *bytes.Buffer = new(bytes.Buffer)
	
	
	
	
	lines := bytes.Split(buf, []byte{'\n'})
	for _, line := range lines {
		if len(line) == 0 {
			
			d.dispatchEvent(currentEvent, dataBuffer.String())

			
			currentEvent = Event{}
			dataBuffer.Reset()
			continue
		}
		if line[0] == byte(':') {
			
			continue
		}

		var field, value []byte
		colonIndex := bytes.IndexRune(line, ':')
		if colonIndex != -1 {
			
			
			
			field = line[:colonIndex]
			
			
			value = line[colonIndex+1:]
			
			if len(value) > 0 && value[0] == ' ' {
				value = value[1:]
			}
		} else {
			
			
			field = line
			value = []byte{}
		}
		
		
		
		switch string(field) {
		case "event":
			
			currentEvent.Event = string(value)
		case "id":
			
			currentEvent.Id = string(value)
		case "retry":
			
			
			
			currentEvent.Id = string(value)
		case "data":
			
			dataBuffer.Write(value)
			
			dataBuffer.WriteString("\n")
		default:
			
			continue
		}
	}
	
	d.dispatchEvent(currentEvent, dataBuffer.String())

	return d.events, nil
}

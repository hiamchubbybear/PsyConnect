package magic

import (
	"bytes"
	"strings"
	"time"

	"github.com/gabriel-vasile/mimetype/internal/charset"
	"github.com/gabriel-vasile/mimetype/internal/json"
)

var (
	
	HTML = markup(
		[]byte("<!DOCTYPE HTML"),
		[]byte("<HTML"),
		[]byte("<HEAD"),
		[]byte("<SCRIPT"),
		[]byte("<IFRAME"),
		[]byte("<H1"),
		[]byte("<DIV"),
		[]byte("<FONT"),
		[]byte("<TABLE"),
		[]byte("<A"),
		[]byte("<STYLE"),
		[]byte("<TITLE"),
		[]byte("<B"),
		[]byte("<BODY"),
		[]byte("<BR"),
		[]byte("<P"),
	)
	
	XML = markup([]byte("<?XML"))
	
	Owl2 = xml(newXMLSig("Ontology", `xmlns="http://www.w3.org/2002/07/owl#"`))
	
	Rss = xml(newXMLSig("rss", ""))
	
	Atom = xml(newXMLSig("feed", `xmlns="http://www.w3.org/2005/Atom"`))
	
	Kml = xml(
		newXMLSig("kml", `xmlns="http://www.opengis.net/kml/2.2"`),
		newXMLSig("kml", `xmlns="http://earth.google.com/kml/2.0"`),
		newXMLSig("kml", `xmlns="http://earth.google.com/kml/2.1"`),
		newXMLSig("kml", `xmlns="http://earth.google.com/kml/2.2"`),
	)
	
	Xliff = xml(newXMLSig("xliff", `xmlns="urn:oasis:names:tc:xliff:document:1.2"`))
	
	Collada = xml(newXMLSig("COLLADA", `xmlns="http://www.collada.org/2005/11/COLLADASchema"`))
	
	Gml = xml(
		newXMLSig("", `xmlns:gml="http://www.opengis.net/gml"`),
		newXMLSig("", `xmlns:gml="http://www.opengis.net/gml/3.2"`),
		newXMLSig("", `xmlns:gml="http://www.opengis.net/gml/3.3/exr"`),
	)
	
	Gpx = xml(newXMLSig("gpx", `xmlns="http://www.topografix.com/GPX/1/1"`))
	
	Tcx = xml(newXMLSig("TrainingCenterDatabase", `xmlns="http://www.garmin.com/xmlschemas/TrainingCenterDatabase/v2"`))
	
	X3d = xml(newXMLSig("X3D", `xmlns:xsd="http://www.w3.org/2001/XMLSchema-instance"`))
	
	Amf = xml(newXMLSig("amf", ""))
	
	Threemf = xml(newXMLSig("model", `xmlns="http://schemas.microsoft.com/3dmanufacturing/core/2015/02"`))
	
	Xfdf = xml(newXMLSig("xfdf", `xmlns="http://ns.adobe.com/xfdf/"`))
	
	VCard = ciPrefix([]byte("BEGIN:VCARD\n"), []byte("BEGIN:VCARD\r\n"))
	
	ICalendar = ciPrefix([]byte("BEGIN:VCALENDAR\n"), []byte("BEGIN:VCALENDAR\r\n"))
	phpPageF  = ciPrefix(
		[]byte("<?PHP"),
		[]byte("<?\n"),
		[]byte("<?\r"),
		[]byte("<? "),
	)
	phpScriptF = shebang(
		[]byte("/usr/local/bin/php"),
		[]byte("/usr/bin/php"),
		[]byte("/usr/bin/env php"),
	)
	
	Js = shebang(
		[]byte("/bin/node"),
		[]byte("/usr/bin/node"),
		[]byte("/bin/nodejs"),
		[]byte("/usr/bin/nodejs"),
		[]byte("/usr/bin/env node"),
		[]byte("/usr/bin/env nodejs"),
	)
	
	Lua = shebang(
		[]byte("/usr/bin/lua"),
		[]byte("/usr/local/bin/lua"),
		[]byte("/usr/bin/env lua"),
	)
	
	Perl = shebang(
		[]byte("/usr/bin/perl"),
		[]byte("/usr/bin/env perl"),
	)
	
	Python = shebang(
		[]byte("/usr/bin/python"),
		[]byte("/usr/local/bin/python"),
		[]byte("/usr/bin/env python"),
	)
	
	Tcl = shebang(
		[]byte("/usr/bin/tcl"),
		[]byte("/usr/local/bin/tcl"),
		[]byte("/usr/bin/env tcl"),
		[]byte("/usr/bin/tclsh"),
		[]byte("/usr/local/bin/tclsh"),
		[]byte("/usr/bin/env tclsh"),
		[]byte("/usr/bin/wish"),
		[]byte("/usr/local/bin/wish"),
		[]byte("/usr/bin/env wish"),
	)
	
	Rtf = prefix([]byte("{\\rtf"))
)





func Text(raw []byte, limit uint32) bool {
	
	if cset := charset.FromBOM(raw); cset != "" {
		return true
	}
	
	for _, b := range raw {
		if b <= 0x08 ||
			b == 0x0B ||
			0x0E <= b && b <= 0x1A ||
			0x1C <= b && b <= 0x1F {
			return false
		}
	}
	return true
}


func Php(raw []byte, limit uint32) bool {
	if res := phpPageF(raw, limit); res {
		return res
	}
	return phpScriptF(raw, limit)
}


func JSON(raw []byte, limit uint32) bool {
	raw = trimLWS(raw)
	
	
	if len(raw) < 2 || (raw[0] != '[' && raw[0] != '{') {
		return false
	}
	parsed, err := json.Scan(raw)
	
	if limit == 0 || len(raw) < int(limit) {
		return err == nil
	}

	
	return parsed == len(raw) && len(raw) > 0
}






func GeoJSON(raw []byte, limit uint32) bool {
	raw = trimLWS(raw)
	if len(raw) == 0 {
		return false
	}
	
	if raw[0] != '{' {
		return false
	}

	s := []byte(`"type"`)
	si, sl := bytes.Index(raw, s), len(s)

	if si == -1 {
		return false
	}

	
	
	if si+sl == len(raw) {
		return false
	}
	
	raw = raw[si+sl:]
	
	raw = trimLWS(raw)
	
	if len(raw) == 0 || raw[0] != ':' {
		return false
	}
	
	raw = trimLWS(raw[1:])

	geoJSONTypes := [][]byte{
		[]byte(`"Feature"`),
		[]byte(`"FeatureCollection"`),
		[]byte(`"Point"`),
		[]byte(`"LineString"`),
		[]byte(`"Polygon"`),
		[]byte(`"MultiPoint"`),
		[]byte(`"MultiLineString"`),
		[]byte(`"MultiPolygon"`),
		[]byte(`"GeometryCollection"`),
	}
	for _, t := range geoJSONTypes {
		if bytes.HasPrefix(raw, t) {
			return true
		}
	}

	return false
}




func NdJSON(raw []byte, limit uint32) bool {
	lCount, hasObjOrArr := 0, false
	raw = dropLastLine(raw, limit)
	var l []byte
	for len(raw) != 0 {
		l, raw = scanLine(raw)
		
		if l = trimRWS(trimLWS(l)); len(l) == 0 {
			continue
		}
		_, err := json.Scan(l)
		if err != nil {
			return false
		}
		if l[0] == '[' || l[0] == '{' {
			hasObjOrArr = true
		}
		lCount++
	}

	return lCount > 1 && hasObjOrArr
}



func HAR(raw []byte, limit uint32) bool {
	s := []byte(`"log"`)
	si, sl := bytes.Index(raw, s), len(s)

	if si == -1 {
		return false
	}

	
	
	if si+sl == len(raw) {
		return false
	}
	
	raw = raw[si+sl:]
	
	raw = trimLWS(raw)
	
	if len(raw) == 0 || raw[0] != ':' {
		return false
	}
	
	raw = trimLWS(raw[1:])

	harJSONTypes := [][]byte{
		[]byte(`"version"`),
		[]byte(`"creator"`),
		[]byte(`"entries"`),
	}
	for _, t := range harJSONTypes {
		si := bytes.Index(raw, t)
		if si > -1 {
			return true
		}
	}

	return false
}


func Svg(raw []byte, limit uint32) bool {
	return bytes.Contains(raw, []byte("<svg"))
}


func Srt(raw []byte, _ uint32) bool {
	line, raw := scanLine(raw)

	
	if string(line) != "1" {
		return false
	}
	line, raw = scanLine(raw)
	secondLine := string(line)
	
	
	if len(secondLine) != 29 {
		return false
	}
	
	
	if strings.Contains(secondLine, ".") {
		return false
	}
	
	ts := strings.Split(secondLine, " --> ")
	if len(ts) != 2 {
		return false
	}
	const layout = "15:04:05,000"
	t0, err := time.Parse(layout, ts[0])
	if err != nil {
		return false
	}
	t1, err := time.Parse(layout, ts[1])
	if err != nil {
		return false
	}
	if t0.After(t1) {
		return false
	}

	line, _ = scanLine(raw)
	
	return len(line) != 0
}



func Vtt(raw []byte, limit uint32) bool {
	
	prefixes := [][]byte{
		{0xEF, 0xBB, 0xBF, 0x57, 0x45, 0x42, 0x56, 0x54, 0x54, 0x0A}, 
		{0xEF, 0xBB, 0xBF, 0x57, 0x45, 0x42, 0x56, 0x54, 0x54, 0x0D}, 
		{0xEF, 0xBB, 0xBF, 0x57, 0x45, 0x42, 0x56, 0x54, 0x54, 0x20}, 
		{0xEF, 0xBB, 0xBF, 0x57, 0x45, 0x42, 0x56, 0x54, 0x54, 0x09}, 
		{0x57, 0x45, 0x42, 0x56, 0x54, 0x54, 0x0A},                   
		{0x57, 0x45, 0x42, 0x56, 0x54, 0x54, 0x0D},                   
		{0x57, 0x45, 0x42, 0x56, 0x54, 0x54, 0x20},                   
		{0x57, 0x45, 0x42, 0x56, 0x54, 0x54, 0x09},                   
	}
	for _, p := range prefixes {
		if bytes.HasPrefix(raw, p) {
			return true
		}
	}

	
	return bytes.Equal(raw, []byte{0xEF, 0xBB, 0xBF, 0x57, 0x45, 0x42, 0x56, 0x54, 0x54}) || 
		bytes.Equal(raw, []byte{0x57, 0x45, 0x42, 0x56, 0x54, 0x54}) 
}


func dropCR(data []byte) []byte {
	if len(data) > 0 && data[len(data)-1] == '\r' {
		return data[0 : len(data)-1]
	}
	return data
}
func scanLine(b []byte) (line, remainder []byte) {
	line, remainder, _ = bytes.Cut(b, []byte("\n"))
	return dropCR(line), remainder
}

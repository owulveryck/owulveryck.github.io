import "strings"

// CloudEvents Specification JSON Schema
@jsonschema(schema="http://json-schema.org/draft-07/schema#")

// Identifies the event.
id: #iddef

// Identifies the context in which an event happened.
source: #sourcedef

// The version of the CloudEvents specification which the event
// uses.
specversion: #specversiondef

// Describes the type of event related to the originating
// occurrence.
type: #typedef

// Content type of the data value. Must adhere to RFC 2046 format.
datacontenttype?: #datacontenttypedef

// Identifies the schema that data adheres to.
dataschema?: #dataschemadef

// Describes the subject of the event in the context of the event
// producer (identified by source).
subject?: #subjectdef

// Timestamp of when the occurrence happened. Must adhere to RFC
// 3339.
time?: #timedef

// The event payload.
data?: #datadef

// Base64 encoded event payload. Must adhere to RFC4648.
data_base64?: #data_base64def

#iddef: strings.MinRunes(1)

#sourcedef: strings.MinRunes(1)

#specversiondef: strings.MinRunes(1)

#typedef: strings.MinRunes(1)

#datacontenttypedef: null | strings.MinRunes(1)

#dataschemadef: null | strings.MinRunes(1)

#subjectdef: null | strings.MinRunes(1)

#timedef: null | strings.MinRunes(1)

#datadef: _ | null | bool | number | string | [...] | {
	...
}

#data_base64def: null | string
...

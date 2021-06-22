package main

import (
	"bytes"
	"context"
	"io"
	"log"
	"net/http"
	"net/http/httptest"
	"reflect"
	"testing"

	"cuelang.org/go/cue"
	"cuelang.org/go/cue/cuecontext"
	"github.com/Shopify/sarama"
	"github.com/cloudevents/sdk-go/protocol/kafka_sarama/v2"
	cloudevents "github.com/cloudevents/sdk-go/v2"
	"github.com/cloudevents/sdk-go/v2/event"
)

// TestHandler is triggering a http test server with the DataProduct's handler,
// send a post request with various payload and check the result
func TestHandler(t *testing.T) {
	ctx := cuecontext.New()
	val := ctx.CompileString(`
	 {
		// first name of the person
		first: =~ "[A-Z].*"
		// Last name of the person
		Last: =~ "[A-Z].*"
		// Age of the person
		Age?: number & < 130
	}
	`)
	ts := httptest.NewServer(&DataProduct{
		definition: val,
	})
	defer ts.Close()
	tests := []struct {
		name           string
		payload        io.Reader
		expectedStatus int
	}{
		{
			name: "valid",
			payload: bytes.NewBufferString(`
			{
				"first": "John",
				"Last": "Doe",
				"Age": 40
			}
`),
			expectedStatus: http.StatusOK,
		},
		{
			name: "invalid",
			payload: bytes.NewBufferString(`
			{
				"first": "John",
				"Last": "Doe",
				"Age": 140
			}
`),
			expectedStatus: http.StatusBadRequest,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			req, _ := http.NewRequest("POST", ts.URL, tt.payload)
			res, err := http.DefaultClient.Do(req)
			if err != nil {
				t.Error(err)
			}
			if res.StatusCode != tt.expectedStatus {
				t.Errorf("expected %v, got %v", tt.expectedStatus, res.StatusCode)
			}
		})
	}

}

func TestDataProduct_CreateEvent(t *testing.T) {
	ctx := cuecontext.New()
	val := ctx.CompileString(`{
		"first": "John",
		"Last": "Doe",
		"Age": 40
	}`)
	expectedEventJSON := []byte(`{
		"specversion": "1.0",
		"id": "",
		"source": "mysource",
		"type": "mytype",
		"subject": "subject",
		"datacontenttype": "application/json",
		"time": "2021-06-14T09:34:34.297881Z",
		"data_base64": "eyJmaXJzdCI6IkpvaG4iLCJMYXN0IjoiRG9lIiwiQWdlIjo0MH0="
	}`)
	expectedEvent := cloudevents.NewEvent()
	err := expectedEvent.UnmarshalJSON(expectedEventJSON)
	if err != nil {
		t.Fatal(err)
	}
	type args struct {
		sub     string
		typ     string
		payload cue.Value
	}
	tests := []struct {
		name    string
		d       *DataProduct
		args    args
		want    event.Event
		wantErr bool
	}{
		{
			"simple valid",
			&DataProduct{
				source: "mysource",
			},
			args{
				sub:     "subject",
				typ:     "mytype",
				payload: val,
			},
			expectedEvent,
			false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := tt.d.CreateEvent(tt.args.sub, tt.args.typ, tt.args.payload)
			if (err != nil) != tt.wantErr {
				t.Errorf("DataProduct.CreateEvent() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got.Data(), tt.want.Data()) {
				t.Errorf("DataProduct.CreateEvent() = %v, want %v", got, tt.want)
			}
			if !reflect.DeepEqual(got.Source(), tt.want.Source()) {
				t.Errorf("DataProduct.CreateEvent() = %v, want %v", got, tt.want)
			}
			if !reflect.DeepEqual(got.Type(), tt.want.Type()) {
				t.Errorf("DataProduct.CreateEvent() = %v, want %v", got, tt.want)
			}
		})
	}
}

func BenchmarkRun(b *testing.B) {
	def, err := extractValue("testdata/definition.cue")
	if err != nil {
		b.Fatal(err)
	}
	// Configure the channel
	saramaConfig := sarama.NewConfig()
	saramaConfig.Version = sarama.V2_0_0_0

	var c cloudevents.Client
	sender, err := kafka_sarama.NewSender([]string{"127.0.0.1:9092"}, saramaConfig, "test-topic")
	if err == nil {

		defer sender.Close(context.Background())

		c, err = cloudevents.NewClient(sender, cloudevents.WithTimeNow(), cloudevents.WithUUIDs())
		if err != nil {
			log.Fatalf("failed to create client, %v", err)
		}
	}
	mux := http.NewServeMux()
	mux.Handle("/", &DataProduct{
		definition: def,
		source:     "sample",
		dest:       c,
	})
	ts := httptest.NewServer(mux)
	defer ts.Close()
	for i := 0; i < b.N; i++ {
		req, err := http.NewRequest("POST", ts.URL, bytes.NewBufferString(`{ "first": "John", "Last": "Doe", "Age": 40 }`))
		if err != nil {
			b.Fatal(err)
		}
		res, err := http.DefaultClient.Do(req)
		if err != nil || res.StatusCode != http.StatusOK {
			b.Fail()
		}
	}

}

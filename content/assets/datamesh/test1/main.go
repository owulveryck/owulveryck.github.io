package main

import (
	"bytes"
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"time"

	"cuelang.org/go/cue"
	"cuelang.org/go/cue/cuecontext"
	"cuelang.org/go/cue/load"
	"cuelang.org/go/encoding/openapi"
	"github.com/Shopify/sarama"
	"github.com/cloudevents/sdk-go/protocol/kafka_sarama/v2"
	cloudevents "github.com/cloudevents/sdk-go/v2"
	"github.com/cloudevents/sdk-go/v2/event"
	"github.com/google/uuid"
)

// DataProduct is a structure that holds the definition as a CUE value
type DataProduct struct {
	definition cue.Value
	source     string // the source of the event
	dest       cloudevents.Client
	// ...
}

// ServeHTTP to make the *DataProduct a http handler
// This is an example, we do not handle the method properly nor we check the content type
// This methods reads the payload from the request an calls the ExtractData method for validation
func (d *DataProduct) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	b, _ := ioutil.ReadAll(r.Body)
	defer r.Body.Close()
	data, err := d.ExtractData(b)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	e, err := d.CreateEvent("mysubject", "mytype", data)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	if d.dest != nil {
		if result := d.dest.Send(
			// Set the producer message key
			kafka_sarama.WithMessageKey(context.Background(), sarama.StringEncoder(e.ID())),
			e,
		); cloudevents.IsUndelivered(result) {
			http.Error(w, result.Error(), http.StatusInternalServerError)
			return
		}
		fmt.Fprint(w, "sent to the channel")
	}
	fmt.Fprint(w, "ok")
}

// ExtractData tries to reconstruct a data from the payload b, and
// unifies it with the data definition.
// then it validates the resulting value
func (d *DataProduct) ExtractData(b []byte) (cue.Value, error) {
	data := d.definition.Context().CompileBytes(b)
	unified := d.definition.Unify(data)
	opts := []cue.Option{
		cue.Attributes(true),
		cue.Definitions(true),
		cue.Hidden(true),
	}
	return data, unified.Validate(opts...)
}

// CreateEvent with the payload encoded in JSON and d.source as source event
func (d *DataProduct) CreateEvent(sub, typ string, payload cue.Value) (event.Event, error) {
	e := cloudevents.NewEvent()
	b, err := payload.MarshalJSON()
	if err != nil {
		return e, err
	}
	e.SetType(typ)
	e.SetSource(d.source)
	e.SetSubject(sub)
	e.SetTime(time.Now())
	e.SetID(uuid.Must(uuid.NewRandom()).String())
	err = e.SetData(cloudevents.ApplicationJSON, b)
	if err != nil {
		return e, err
	}
	_, err = json.MarshalIndent(e, " ", " ")
	return e, err
}
func generateOpenAPI(defFile string, config *load.Config) ([]byte, error) {
	buildInstances := load.Instances([]string{defFile}, config)
	insts := cue.Build(buildInstances)
	b, err := openapi.Gen(insts[0], nil)
	if err != nil {
		return nil, err
	}
	var out bytes.Buffer
	err = json.Indent(&out, b, "", "   ")
	if err != nil {
		return nil, err
	}
	return out.Bytes(), nil
}

func extractValue(defFile string) (cue.Value, error) {
	ctx := cuecontext.New()

	buildInstances := load.Instances([]string{defFile}, nil)
	values, err := ctx.BuildInstances(buildInstances)
	v := values[0]
	it, err := v.Fields(cue.All())
	if err != nil {
		return cue.Value{}, err
	}
	for it.Next() {
		if it.Selector().IsDefinition() {
			return it.Value(), nil
		}
	}
	return cue.Value{}, errors.New("nothing found")
}

func main() {
	run()
}

func run() {

	// Read the definition file name
	defFile := os.Args[1]
	openapi, err := generateOpenAPI(defFile, nil)
	if err != nil {
		log.Fatal(err)
	}
	def, err := extractValue(defFile)
	if err != nil {
		log.Fatal(err)
	}
	// Configure the channel
	saramaConfig := sarama.NewConfig()
	saramaConfig.Version = sarama.V2_0_0_0

	var c cloudevents.Client
	sender, err := kafka_sarama.NewSender([]string{"127.0.0.1:9092"}, saramaConfig, "test-topic")
	if err != nil {
		log.Printf("failed to create protocol: %s", err.Error())
	} else {

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
	mux.HandleFunc("/openapi", func(w http.ResponseWriter, _ *http.Request) {
		io.Copy(w, bytes.NewBuffer(openapi))
	})
	log.Fatal(http.ListenAndServe(":8181", mux))
}

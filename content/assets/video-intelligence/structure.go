package main

import (
	"github.com/cathalgarvey/fmtless/encoding/json"
	"sort"
	"time"
)

type segment struct {
	EndTimeOffset   time.Duration `json:"endTimeOffset"`
	StartTimeOffset time.Duration `json:"startTimeOffset"`
}

func (s *segment) UnmarshalJSON(b []byte) error {
	type seg struct {
		EndTimeOffset   string `json:"endTimeOffset"`
		StartTimeOffset string `json:"startTimeOffset"`
	}
	var tmp seg
	err := json.Unmarshal(b, &tmp)
	if err != nil {
		return err
	}
	if tmp.EndTimeOffset == "-1" {
		tmp.EndTimeOffset = "0"
	}
	if tmp.StartTimeOffset == "-1" {
		tmp.StartTimeOffset = "0"
	}
	s.EndTimeOffset, err = time.ParseDuration(tmp.EndTimeOffset + "us")
	if err != nil {
		return err
	}
	s.StartTimeOffset, err = time.ParseDuration(tmp.StartTimeOffset + "us")
	if err != nil {
		return err
	}
	return nil
}

type videoIntelligence struct {
	Response struct {
		AnnotationResults []struct {
			ShotAnnotations []struct {
				EndTimeOffset   string `json:"endTimeOffset"`
				StartTimeOffset string `json:"startTimeOffset,omitempty"`
			} `json:"shotAnnotations"`
			LabelAnnotations []struct {
				Locations []struct {
					Level      string   `json:"level"`
					Confidence float64  `json:"confidence"`
					Segment    *segment `json:"segment"`
				} `json:"locations"`
				LanguageCode string `json:"languageCode"`
				Description  string `json:"description"`
			} `json:"labelAnnotations"`
			InputURI string `json:"inputUri"`
		} `json:"annotationResults"`
		Type string `json:"@type"`
	} `json:"response"`
	Done     bool `json:"done"`
	Metadata struct {
		AnnotationProgress []struct {
			UpdateTime      time.Time `json:"updateTime"`
			StartTime       time.Time `json:"startTime"`
			ProgressPercent int       `json:"progressPercent"`
			InputURI        string    `json:"inputUri"`
		} `json:"annotationProgress"`
		Type string `json:"@type"`
	} `json:"metadata"`
	Name string `json:"name"`
}

type annotation struct {
	t      time.Duration
	add    string
	remove string
	labels []string
}
type annotations struct {
	annotations []annotation
	labels      []string
}

func (a *annotations) Len() int {
	return len(a.annotations)
}

func (a *annotations) Swap(i, j int) {
	a.annotations[i], a.annotations[j] = a.annotations[j], a.annotations[i]
}

func (a *annotations) Less(i, j int) bool {
	if a.annotations[i].t < a.annotations[j].t {
		return true
	}
	return false
}

func processData() []annotation {
	var vi videoIntelligence
	json.Unmarshal(data, &vi)
	var m []annotation

	for _, la := range vi.Response.AnnotationResults[0].LabelAnnotations {
		for _, l := range la.Locations {
			if l.Segment.StartTimeOffset == 0 || l.Segment.EndTimeOffset == 0 {
				continue
			}

			m = append(m, annotation{t: l.Segment.StartTimeOffset, add: la.Description})
			m = append(m, annotation{t: l.Segment.EndTimeOffset, remove: la.Description})
		}
	}
	a := &annotations{m, []string{}}
	var result []annotation
	sort.Sort(a)
	var last time.Duration
	for i, v := range a.annotations {
		if v.add != "" {
			a.labels = append(a.labels, v.add)
		}
		if v.remove != "" {
			for j, vv := range a.labels {
				if vv == v.remove {
					a.labels = append(a.labels[:j], a.labels[j+1:]...)
					break
				}
			}
		}
		a.annotations[i].labels = a.labels
		if a.annotations[i].t != last {
			var lbl []string
			for _, lb := range a.labels {
				lbl = append(lbl, lb)
			}
			result = append(result, annotation{t: a.annotations[i].t, labels: lbl})
			last = a.annotations[i].t
		}
	}

	//return a.annotations
	return result
}

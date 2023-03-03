package main

// go:generate gopherjs build main.go -o app.js -m
// +build ignore

import (
	"time"

	"github.com/gopherjs/gopherjs/js"
)

var anns []annotation

func main() {

	// Process the data
	anns = processData()
	// var tag = document.createElement("script");
	tag := js.Global.Get("document").Call("createElement", "script")
	// tag.src = "//www.youtube.com/iframe_api";
	tag.Set("src", "//www.youtube.com/iframe_api")
	// var firstScriptTag = document.getElementsByTagName("script")[0];
	// firstScriptTag.parentNode.insertBefore(tag, firstScriptTag);
	scriptTags := js.Global.Get("document").Call("getElementsByTagName", "script")
	scriptTags.Index(0).Get("parentNode").Call("insertBefore", tag, scriptTags.Index(0))
	// // This function creates an <iframe> (and YouTube player)
	// // after the API code downloads.
	// var player;
	// window.onYouTubeIframeAPIReady = function() {
	//   player = new YT.Player("player", {
	//       "height": "315",
	//       "width": "560",
	//       "videoId": "A0yQ0dPhkOg",
	//       "events": {
	//       "onReady": onPlayerReady,
	//       "onStateChange": onPlayerStateChange
	//       }
	//       });
	// }
	// Create one configuration object that will be transpiled in json Object
	// and passed to the constructor of the player
	// See https://github.com/gopherjs/gopherjs/wiki/JavaScript-Tips-and-Gotchas
	type ytConfig struct {
		*js.Object        // so far so good
		Height     string `js:"height"`
		Width      string `js:"width"`
		VideoID    string `js:"videoId"`
	}
	// Create the configuration
	config := &ytConfig{Object: js.Global.Get("Object").New()}
	//config.Height = "315"
	//config.Width = "560"
	//config.Height = "315"
	config.Width = "100%"
	config.VideoID = "O1xitYlzDsc"
	js.Global.Get("window").Set("onYouTubeIframeAPIReady", func() {
		// Then create a new Player instance called "player", actually creating an iFrame "player" instead of the
		// Div identified by "player"
		js.Global.Get("launchyt").Call("addEventListener", "click", func() {
			go func() {
				var player *ytPlayer
				player = &ytPlayer{*(js.Global.Get("YT").Get("Player").New("player", config)), make(chan string)}
				player.Call("addEventListener", "onReady", player.onPlayerReady)
				player.Call("addEventListener", "onStateChange", player.onPlayerStateChange)
			}()
		})
	})
}

type ytPlayer struct {
	js.Object
	state chan string
}

func (yt *ytPlayer) onPlayerReady(event *js.Object) {
	// Trigger the goroutine that will display the current time of the video
	go func() {
		var state string
		var currLabels []string
		for {
			select {
			case state = <-yt.state:
			default:
			}

			if state == "1" {
				t, _ := yt.getCurrentTime()
				for k := 0; k < len(anns)-1; k++ {
					if t >= anns[k].t && t < anns[k+1].t {
						if !testEq(currLabels, anns[k].labels) {
							currLabels = anns[k].labels
							displayLabels(anns[k].t, currLabels)
						}
						break
					}
				}
			}
			time.Sleep(10 * time.Millisecond)
		}
	}()
}
func displayLabels(t time.Duration, labels []string) {
	l1 := js.Global.Get("document").Call("getElementById", "result1")
	l2 := js.Global.Get("document").Call("getElementById", "result2")
	l3 := js.Global.Get("document").Call("getElementById", "result3")
	l4 := js.Global.Get("document").Call("getElementById", "result4")
	l5 := js.Global.Get("document").Call("getElementById", "result5")

	l5.Set("innerHTML", l4.Get("innerHTML").String())
	l4.Set("innerHTML", l3.Get("innerHTML").String())
	l3.Set("innerHTML", l2.Get("innerHTML").String())
	l2.Set("innerHTML", l1.Get("innerHTML").String())

	var ss string
	for _, l := range labels {
		ss = ss + " | " + l
	}
	l1.Set("innerHTML", ss)

}
func (yt *ytPlayer) onPlayerStateChange(event *js.Object) {
	go func() {
		yt.state <- event.Get("data").String()
	}()
}

func (yt *ytPlayer) getCurrentTime() (time.Duration, error) {
	return time.ParseDuration(yt.Call("getCurrentTime").String() + "s")
}

func testEq(a, b []string) bool {

	if a == nil && b == nil {
		return true
	}

	if a == nil || b == nil {
		return false
	}

	if len(a) != len(b) {
		return false
	}

	for i := range a {
		if a[i] != b[i] {
			return false
		}
	}

	return true
}

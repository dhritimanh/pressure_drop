package engine

import (
	"fmt"
	"log"

	"github.com/hajimehoshi/ebiten/v2/audio"
)

type AudioPlayer struct {
	Context *audio.Context
	// Will hold *audio.Player references once physical .wav files exist.
}

func NewAudioPlayer() *AudioPlayer {
	// Initialize audio context at standard 44100 hz rate
	ctx := audio.NewContext(44100)
	log.Println("[Audio] Context Initialized. Ready for Phase 6 .wav attachments in assets/audio/...")
	return &AudioPlayer{Context: ctx}
}

func (a *AudioPlayer) PlayBGM() {
	fmt.Println("[Audio] background_loop.wav triggered")
}

func (a *AudioPlayer) PlayPing() {
	fmt.Println("[Audio] SFX: ping.wav")
}

func (a *AudioPlayer) PlaySnap() {
	fmt.Println("[Audio] SFX: cable_snap.wav")
}

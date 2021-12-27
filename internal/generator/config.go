/*
 * Copyright (c) 2021 Andreas Signer <asigner@gmail.com>
 *
 * This file is part of gosfxr.
 *
 * gosfxr is free software: you can redistribute it and/or
 * modify it under the terms of the GNU General Public License as
 * published by the Free Software Foundation, either version 3 of the
 * License, or (at your option) any later version.
 *
 * gosfxr is distributed in the hope that it will be useful,
 * but WITHOUT ANY WARRANTY; without even the implied warranty of
 * MERCHANTABILITY or FITNESS FOR A PARTICULAR PURPOSE.  See the
 * GNU General Public License for more details.
 *
 * You should have received a copy of the GNU General Public License
 * along with gosfxr.  If not, see <http://www.gnu.org/licenses/>.
 */
package generator

import (
	"encoding/json"
	"math"
	"math/rand"
)

type Waveform int

const (
	WaveformSquare Waveform = iota
	WaveformSawtooth
	WaveformSine
	WaveformNoise
)

/*
cat generator.go |
sed -e 's/p_env_attack/EnvelopeAttack/g' |
sed -e 's/p_env_punch/EnvelopeSustainPunch/g' |
sed -e 's/p_env_sustain/EnvelopeSustain/g' |
sed -e 's/p_env_decay/EnvelopeDecay/g' |
sed -e 's/p_base_freq/FreqStart/g' |
sed -e 's/p_freq_limit/FreqMinCutoff/g' |
sed -e 's/p_freq_ramp/FreqSlide/g' |
sed -e 's/p_freq_dramp/FreqDeltaSlide/g' |
sed -e 's/p_vib_strength/VibDepth/g' |
sed -e 's/p_vib_speed/VibSpeed/g' |
sed -e 's/p_vib_delay/VibDelay/g' |
sed -e 's/p_arp_mod/ArpFreqMult/g' |
sed -e 's/p_arp_speed/ArpChangeSpeed/g' |
sed -e 's/p_duty_ramp/DutyCycleSweep/g' |
sed -e 's/p_duty/DutyCycle/g' |
sed -e 's/p_repeat_speed/RepeatRate/g' |
sed -e 's/p_pha_offset/PhaserOffset/g' |
sed -e 's/p_pha_ramp/PhaserSweep/g' |
sed -e 's/p_lpf_freq/LPCutoffFreq/g' |
sed -e 's/p_lpf_ramp/LPCutoffSweep/g' |
sed -e 's/p_lpf_resonance/LPResonance/g' |
sed -e 's/p_hpf_freq/HPCutoffFreq/g' |
sed -e 's/p_hpf_ramp/HPCutoffSweep/g' \
> generator.go.NEW

*/

func frnd(r float64) float64 {
	return float64(rand.Int31n(10000)) / 10000.0 * r
}

func rnd(n int32) int32 {
	return rand.Int31n(n + 1)
}

func brnd() bool {
	return rand.Int31n(2) > 0
}

type Config struct {
	Waveform Waveform  `json:"waveform"`
	Volume float64 `json:"volume"`

	// Envelope
	EnvelopeAttack       float64 `json:"env_attack"`
	EnvelopeSustain      float64 `json:"env_sustain"`
	EnvelopeSustainPunch float64 `json:"env_punch"`
	EnvelopeDecay        float64 `json:"env_decay"`

	// Tone
	FreqStart      float64 `json:"base_freq"`
	FreqMinCutoff  float64 `json:"freq_limit"`
	FreqSlide      float64 `json:"freq_ramp"`
	FreqDeltaSlide float64 `json:"freq_dramp"`

	// Vibrato
	VibDepth float64 `json:"vib_strength"`
	VibSpeed float64 `json:"vib_speed"`
	VibDelay float64 `json:"vid_delay"`

    // Arpaggio
	ArpFreqMult    float64 `json:"arp_mod"`
	ArpChangeSpeed float64 `json:"arp_speed"`

	// Square wave duty (proportion of time signal is high vs. low)
	DutyCycle      float64 `json:"duty"`
	DutyCycleSweep float64 `json:"duty_ramp"`

	RepeatRate float64 `json:"repeat_speed"`

	// Phaser
	PhaserOffset float64 `json:"pha_offset"`
	PhaserSweep  float64 `json:"pha_ramp"`

	// Low-Pass Filter
	LPCutoffFreq  float64 `json:"lpf_freq"`
	LPCutoffSweep float64 `json:"lpf_ramp"`
	LPResonance   float64 `json:"lpf_resonance"`

	// High-Pass Filter
	HPCutoffFreq  float64 `json:"hpf_freq"`
	HPCutoffSweep float64 `json:"hpf_ramp"`
}

func NewConfig() *Config {
	g := &Config{}
	g.Reset()
	return g
}

func (g *Config) InitFromJson(j []byte) {
	json.Unmarshal(j, g)
}

func (g *Config) ToJson() []byte {
	content, _ := json.MarshalIndent(*g, "", "    ")
	return content
}

func (g *Config) Reset() {
	g.Waveform = WaveformSquare
	g.Volume = 1.0

	g.FreqStart = 0.3
	g.FreqMinCutoff = 0.0
	g.FreqSlide = 0.0
	g.FreqDeltaSlide = 0.0
	g.DutyCycle = 0.0
	g.DutyCycleSweep = 0.0

	g.VibDepth = 0.0
	g.VibSpeed = 0.0
	g.VibDelay = 0.0

	g.EnvelopeAttack = 0.0
	g.EnvelopeSustain = 0.3
	g.EnvelopeDecay = 0.4
	g.EnvelopeSustainPunch = 0.0

	g.LPResonance = 0.0
	g.LPCutoffFreq = 1.0
	g.LPCutoffSweep = 0.0
	g.HPCutoffFreq = 0.0
	g.HPCutoffSweep = 0.0

	g.PhaserOffset = 0.0
	g.PhaserSweep = 0.0

	g.RepeatRate = 0.0

	g.ArpChangeSpeed = 0.0
	g.ArpFreqMult = 0.0
}

func (g *Config) PresetPickup() {
	g.Reset()
	g.FreqStart = 0.4 + frnd(0.5)
	g.EnvelopeAttack = 0.0
	g.EnvelopeSustain = frnd(0.1)
	g.EnvelopeDecay = 0.1 + frnd(0.4)
	g.EnvelopeSustainPunch = 0.3 + frnd(0.3)
	if brnd() {
		g.ArpChangeSpeed = 0.5 + frnd(0.2)
		g.ArpFreqMult = 0.2 + frnd(0.4)
	}
}

func (g *Config) PresetLaser() {
	g.Reset()
	g.Waveform = Waveform(rnd(2))
	if g.Waveform == WaveformSine && rnd(1) == 1 {
		// Make Sine less propable? But why?
		g.Waveform = Waveform(rnd(1))
	}

	g.FreqStart = 0.5 + frnd(0.5)
	g.FreqMinCutoff = g.FreqStart - 0.2 - frnd(0.6)
	if g.FreqMinCutoff < 0.2 {
		g.FreqMinCutoff = 0.2
	}
	g.FreqSlide = -0.15 - frnd(0.2)
	if rnd(2) == 0 {
		g.FreqStart = 0.3 + frnd(0.6)
		g.FreqMinCutoff = frnd(0.1)
		g.FreqSlide = -0.35 - frnd(0.3)
	}
	if brnd() {
		g.DutyCycle = frnd(0.5)
		g.DutyCycleSweep = frnd(0.2)
	} else {
		g.DutyCycle = 0.4 + frnd(0.5)
		g.DutyCycleSweep = -frnd(0.7)
	}
	g.EnvelopeAttack = 0.0
	g.EnvelopeSustain = 0.1 + frnd(0.2)
	g.EnvelopeDecay = frnd(0.4)
	if brnd() {
		g.EnvelopeSustainPunch = frnd(0.3)
	}
	if rnd(2) == 0 {
		g.PhaserOffset = frnd(0.2)
		g.PhaserSweep = -frnd(0.2)
	}
	if brnd() {
		g.HPCutoffFreq = frnd(0.3)
	}
}

func (g *Config) PresetExplosion() {
	g.Reset()
	g.Waveform = WaveformNoise
	if brnd() {
		g.FreqStart = 0.1 + frnd(0.4)
		g.FreqSlide = -0.1 + frnd(0.4)
	} else {
		g.FreqStart = 0.2 + frnd(0.7)
		g.FreqSlide = -0.2 - frnd(0.2)
	}
	g.FreqStart *= g.FreqStart
	if rnd(4) == 0 {
		g.FreqSlide = 0.0
	}
	if rnd(2) == 0 {
		g.RepeatRate = 0.3 + frnd(0.5)
	}
	g.EnvelopeAttack = 0.0
	g.EnvelopeSustain = 0.1 + frnd(0.3)
	g.EnvelopeDecay = frnd(0.5)
	if rnd(1) == 0 {
		g.PhaserOffset = -0.3 + frnd(0.9)
		g.PhaserSweep = -frnd(0.3)
	}
	g.EnvelopeSustainPunch = 0.2 + frnd(0.6)
	if brnd() {
		g.VibDepth = frnd(0.7)
		g.VibSpeed = frnd(0.6)
	}
	if rnd(2) == 0 {
		g.ArpChangeSpeed = 0.6 + frnd(0.3)
		g.ArpFreqMult = 0.8 - frnd(1.6)
	}
}

func (g *Config) PresetPowerup() {
	g.Reset()
	if brnd() {
		g.Waveform = WaveformSawtooth
	} else {
		g.DutyCycle = frnd(0.6)
	}

	if brnd() {
		g.FreqStart = 0.2 + frnd(0.3)
		g.FreqSlide = 0.1 + frnd(0.4)
		g.RepeatRate = 0.4 + frnd(0.4)
	} else {
		g.FreqStart = 0.2 + frnd(0.3)
		g.FreqSlide = 0.05 + frnd(0.2)
		if brnd() {
			g.VibDepth = frnd(0.7)
			g.VibSpeed = frnd(0.6)
		}
	}
	g.EnvelopeAttack = 0.0
	g.EnvelopeSustain = frnd(0.4)
	g.EnvelopeDecay = 0.1 + frnd(0.4)
}

func (g *Config) PresetHit() {
	g.Reset()
	g.Waveform = Waveform(rnd(2))
	if g.Waveform == WaveformSine {
		g.Waveform = WaveformNoise
	}
	if g.Waveform == WaveformSquare {
		g.DutyCycle = frnd(0.6)
	}
	g.FreqStart = 0.2 + frnd(0.6)
	g.FreqSlide = -0.3 - frnd(0.4)
	g.EnvelopeAttack = 0.0
	g.EnvelopeSustain = frnd(0.1)
	g.EnvelopeDecay = 0.1 + frnd(0.2)
	if brnd() {
		g.HPCutoffFreq = frnd(0.3)
	}
}

func (g *Config) PresetJump() {
	g.Reset()
	g.Waveform = WaveformSquare
	g.DutyCycle = frnd(0.6)
	g.FreqStart = 0.3 + frnd(0.3)
	g.FreqSlide = 0.1 + frnd(0.2)
	g.EnvelopeAttack = 0.0
	g.EnvelopeSustain = 0.1 + frnd(0.3)
	g.EnvelopeDecay = 0.1 + frnd(0.2)
	if brnd() {
		g.HPCutoffFreq = frnd(0.3)
	}
	if brnd() {
		g.LPCutoffFreq = 1.0 - frnd(0.6)
	}
}

func (g *Config) PresetBlip() {
	g.Reset()
	g.Waveform = Waveform(rnd(1))
	if g.Waveform == WaveformSquare {
		g.DutyCycle = frnd(0.6)
	}
	g.FreqStart = 0.2 + frnd(0.4)
	g.EnvelopeAttack = 0.0
	g.EnvelopeSustain = 0.1 + frnd(0.1)
	g.EnvelopeDecay = frnd(0.2)
	g.HPCutoffFreq = 0.1
}

func (g *Config) Mutate() {
	if brnd() {
		g.FreqStart += frnd(0.1) - 0.05
	}
	//		if brnd() {g.FreqMinCutoff+=frnd(0.1)-0.05;       }
	if brnd() {
		g.FreqSlide += frnd(0.1) - 0.05
	}
	if brnd() {
		g.FreqDeltaSlide += frnd(0.1) - 0.05
	}
	if brnd() {
		g.DutyCycle += frnd(0.1) - 0.05
	}
	if brnd() {
		g.DutyCycleSweep += frnd(0.1) - 0.05
	}
	if brnd() {
		g.VibDepth += frnd(0.1) - 0.05
	}
	if brnd() {
		g.VibSpeed += frnd(0.1) - 0.05
	}
	if brnd() {
		g.VibDelay += frnd(0.1) - 0.05
	}
	if brnd() {
		g.EnvelopeAttack += frnd(0.1) - 0.05
	}
	if brnd() {
		g.EnvelopeSustain += frnd(0.1) - 0.05
	}
	if brnd() {
		g.EnvelopeDecay += frnd(0.1) - 0.05
	}
	if brnd() {
		g.EnvelopeSustainPunch += frnd(0.1) - 0.05
	}
	if brnd() {
		g.LPResonance += frnd(0.1) - 0.05
	}
	if brnd() {
		g.LPCutoffFreq += frnd(0.1) - 0.05
	}
	if brnd() {
		g.LPCutoffSweep += frnd(0.1) - 0.05
	}
	if brnd() {
		g.HPCutoffFreq += frnd(0.1) - 0.05
	}
	if brnd() {
		g.HPCutoffSweep += frnd(0.1) - 0.05
	}
	if brnd() {
		g.PhaserOffset += frnd(0.1) - 0.05
	}
	if brnd() {
		g.PhaserSweep += frnd(0.1) - 0.05
	}
	if brnd() {
		g.RepeatRate += frnd(0.1) - 0.05
	}
	if brnd() {
		g.ArpChangeSpeed += frnd(0.1) - 0.05
	}
	if brnd() {
		g.ArpFreqMult += frnd(0.1) - 0.05
	}
}

func (g *Config) Randomize() {
	g.FreqStart = math.Pow(frnd(2.0)-1.0, 2.0)
	if brnd() {
		g.FreqStart = math.Pow(frnd(2.0)-1.0, 3.0) + 0.5
	}
	g.FreqMinCutoff = 0.0
	g.FreqSlide = math.Pow(frnd(2.0)-1.0, 5.0)
	if g.FreqStart > 0.7 && g.FreqSlide > 0.2 {
		g.FreqSlide = -g.FreqSlide
	}
	if g.FreqStart < 0.2 && g.FreqSlide < -0.05 {
		g.FreqSlide = -g.FreqSlide
	}
	g.FreqDeltaSlide = math.Pow(frnd(2.0)-1.0, 3.0)
	g.DutyCycle = frnd(2.0) - 1.0
	g.DutyCycleSweep = math.Pow(frnd(2.0)-1.0, 3.0)
	g.VibDepth = math.Pow(frnd(2.0)-1.0, 3.0)
	g.VibSpeed = frnd(2.0) - 1.0
	g.VibDelay = frnd(2.0) - 1.0
	g.EnvelopeAttack = math.Pow(frnd(2.0)-1.0, 3.0)
	g.EnvelopeSustain = math.Pow(frnd(2.0)-1.0, 2.0)
	g.EnvelopeDecay = frnd(2.0) - 1.0
	g.EnvelopeSustainPunch = math.Pow(frnd(0.8), 2.0)
	if g.EnvelopeAttack+g.EnvelopeSustain+g.EnvelopeDecay < 0.2 {
		g.EnvelopeSustain += 0.2 + frnd(0.3)
		g.EnvelopeDecay += 0.2 + frnd(0.3)
	}
	g.LPResonance = frnd(2.0) - 1.0
	g.LPCutoffFreq = 1.0 - math.Pow(frnd(1.0), 3.0)
	g.LPCutoffSweep = math.Pow(frnd(2.0)-1.0, 3.0)
	if g.LPCutoffFreq < 0.1 && g.LPCutoffSweep < -0.05 {
		g.LPCutoffSweep = -g.LPCutoffSweep
	}
	g.HPCutoffFreq = math.Pow(frnd(1.0), 5.0)
	g.HPCutoffSweep = math.Pow(frnd(2.0)-1.0, 5.0)
	g.PhaserOffset = math.Pow(frnd(2.0)-1.0, 3.0)
	g.PhaserSweep = math.Pow(frnd(2.0)-1.0, 3.0)
	g.RepeatRate = frnd(2.0) - 1.0
	g.ArpChangeSpeed = frnd(2.0) - 1.0
	g.ArpFreqMult = frnd(2.0) - 1.0
}

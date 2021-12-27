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
	"math"
)

type Generator struct {
	cfg Config

	phase         int
	fperiod       float64
	fmaxperiod    float64
	fslide        float64
	fdslide       float64
	period        int
	square_duty   float64
	square_slide  float64
	env_stage     int
	env_time      int
	env_length    [3]int
	env_vol       float64
	fphase        float64
	fdphase       float64
	iphase        int
	phaser_buffer [1024]float64
	ipp           int
	noise_buffer  [32]float64
	fltp          float64
	fltdp         float64
	fltw          float64
	fltw_d        float64
	fltdmp        float64
	fltphp        float64
	flthp         float64
	flthp_d       float64
	vib_phase     float64
	vib_speed     float64
	vib_amp       float64
	rep_time      int
	rep_limit     int
	arp_time      int
	arp_limit     int
	arp_mod       float64
}

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

func (g *Generator) init() {
	g.phase = 0
	// reset filter
	g.fltp = 0.0
	g.fltdp = 0.0
	g.fltw = math.Pow(g.cfg.LPCutoffFreq, 3.0) * 0.1
	g.fltw_d = 1.0 + g.cfg.LPCutoffSweep*0.0001
	g.fltdmp = 5.0 / (1.0 + math.Pow(g.cfg.LPResonance, 2.0)*20.0) * (0.01 + g.fltw)
	if g.fltdmp > 0.8 {
		g.fltdmp = 0.8
	}
	g.fltphp = 0.0
	g.flthp = math.Pow(g.cfg.HPCutoffFreq, 2.0) * 0.1
	g.flthp_d = 1.0 + g.cfg.HPCutoffSweep*0.0003

	// reset vibrato
	g.vib_phase = 0.0
	g.vib_speed = math.Pow(g.cfg.VibSpeed, 2.0) * 0.01
	g.vib_amp = g.cfg.VibDepth * 0.5

	// reset envelope
	g.env_vol = 0.0
	g.env_stage = 0
	g.env_time = 0
	g.env_length[0] = (int)(g.cfg.EnvelopeAttack * g.cfg.EnvelopeAttack * 100000.0)
	g.env_length[1] = (int)(g.cfg.EnvelopeSustain * g.cfg.EnvelopeSustain * 100000.0)
	g.env_length[2] = (int)(g.cfg.EnvelopeDecay * g.cfg.EnvelopeDecay * 100000.0)

	g.fphase = math.Pow(g.cfg.PhaserOffset, 2.0) * 1020.0
	if g.cfg.PhaserOffset < 0.0 {
		g.fphase = -g.fphase
	}
	g.fdphase = math.Pow(g.cfg.PhaserSweep, 2.0) * 1.0
	if g.cfg.PhaserSweep < 0 {
		g.fdphase = -g.fdphase
	}
	g.iphase = int(math.Abs(g.fphase)) // abs((int)fphase)
	g.ipp = 0
	for i := 0; i < 1024; i++ {
		g.phaser_buffer[i] = 0
	}

	for i := 0; i < 32; i++ {
		g.noise_buffer[i] = frnd(2.0) - 1.0
	}

	g.rep_time = 0
	g.rep_limit = int(math.Pow(1.0-g.cfg.RepeatRate, 2.0)*20000 + 32)
	if g.cfg.RepeatRate == 0.0 {
		g.rep_limit = 0
	}
}

func (g *Generator) initForRepeat() {
	cfg := g.cfg
	g.fperiod = 100.0 / (cfg.FreqStart*cfg.FreqStart + 0.001)
	g.period = int(g.fperiod)
	g.fmaxperiod = 100.0 / (cfg.FreqMinCutoff*cfg.FreqMinCutoff + 0.001)
	g.fslide = 1.0 - math.Pow(cfg.FreqSlide, 3)*0.01
	g.fdslide = -math.Pow(cfg.FreqDeltaSlide, 3) * 0.000001

	g.square_duty = 0.5 - cfg.DutyCycle*0.5
	g.square_slide = -cfg.DutyCycleSweep * 0.00005

	if cfg.ArpFreqMult >= 0 {
		g.arp_mod = 1.0 - math.Pow(cfg.ArpFreqMult, 2.0)*0.9
	} else {
		g.arp_mod = 1.0 + math.Pow(cfg.ArpFreqMult, 2.0)*10.0
	}
	g.arp_time = 0
	g.arp_limit = int(math.Pow(1.0-cfg.ArpChangeSpeed, 2.0)*20000 + 32)
	if cfg.ArpChangeSpeed == 1.0 {
		g.arp_limit = 0
	}

}

func New(cfg *Config) *Generator {
	g := &Generator{
		cfg: *cfg,
	}
	return g
}

const masterVolume = 0.05

func (g *Generator) Generate() []float64 {
	g.init()
	g.initForRepeat()

	g.rep_time = 0
	var buffer []float64

	for {
		g.rep_time++
		if g.rep_limit != 0 && g.rep_time >= g.rep_limit {
			g.rep_time = 0
			g.initForRepeat()
		}

		// frequency envelopes/arpeggios
		g.arp_time++
		if g.arp_limit != 0 && g.arp_time >= g.arp_limit {
			g.arp_limit = 0
			g.fperiod *= g.arp_mod
		}
		g.fslide += g.fdslide
		g.fperiod *= g.fslide
		if g.fperiod > g.fmaxperiod {
			g.fperiod = g.fmaxperiod
			if g.cfg.FreqMinCutoff > 0.0 {
				break
			}
		}
		rfperiod := g.fperiod
		if g.vib_amp > 0.0 {
			g.vib_phase += g.vib_speed
			rfperiod = g.fperiod * (1.0 + math.Sin(g.vib_phase)*g.vib_amp)
		}
		g.period = int(rfperiod)
		if g.period < 8 {
			g.period = 8
		}
		g.square_duty += g.square_slide
		if g.square_duty < 0.0 {
			g.square_duty = 0.0
		}
		if g.square_duty > 0.5 {
			g.square_duty = 0.5
		}

		// volume envelope
		g.env_time++
		if g.env_time > g.env_length[g.env_stage] {
			g.env_time = 0
			g.env_stage++
			if g.env_stage == 3 {
				break
			}
		}
		switch g.env_stage {
		case 0:
			g.env_vol = float64(g.env_time) / float64(g.env_length[0])
		case 1:
			g.env_vol = 1.0 + math.Pow(1.0-float64(g.env_time)/float64(g.env_length[1]), 1.0)*2.0*g.cfg.EnvelopeSustainPunch
		case 2:
			g.env_vol = 1.0 - float64(g.env_time)/float64(g.env_length[2])
		}

		// phaser step
		g.fphase += g.fdphase
		g.iphase = int(math.Abs(g.fphase))
		if g.iphase > 1023 {
			g.iphase = 1023
		}

		if g.flthp_d != 0.0 {
			g.flthp *= g.flthp_d
		}
		if g.flthp < 0.00001 {
			g.flthp = 0.00001
		}
		if g.flthp > 0.1 {
			g.flthp = 0.1
		}

		ssample := 0.0
		for si := 0; si < 8; si++ { // 8x supersampling
			sample := 0.0
			g.phase++
			if g.phase >= g.period {
				g.phase %= g.period
				if g.cfg.Waveform == WaveformNoise {
					for i := 0; i < 32; i++ {
						g.noise_buffer[i] = frnd(2.0) - 1.0
					}
				}
			}

			// base waveform
			fp := float64(g.phase) / float64(g.period)
			switch g.cfg.Waveform {
			case WaveformSquare:
				if fp > g.square_duty {
					sample = 0.5
				} else {
					sample = -0.5
				}
			case WaveformSawtooth:
				sample = 1.0 - fp*2
			case WaveformSine:
				sample = math.Sin(fp * 2 * math.Pi)
			case WaveformNoise:
				sample = g.noise_buffer[g.phase*32/g.period]
			}

			// lp filter
			pp := g.fltp
			g.fltw *= g.fltw_d
			if g.fltw < 0.0 {
				g.fltw = 0.0
			}
			if g.fltw > 0.1 {
				g.fltw = 0.1
			}
			if g.cfg.LPCutoffFreq != 1.0 {
				g.fltdp += (sample - g.fltp) * g.fltw
				g.fltdp -= g.fltdp * g.fltdmp
			} else {
				g.fltp = sample
				g.fltdp = 0.0
			}
			g.fltp += g.fltdp

			// hp filter
			g.fltphp += g.fltp - pp
			g.fltphp -= g.fltphp * g.flthp
			sample = g.fltphp

			// phaser
			g.phaser_buffer[g.ipp&1023] = sample
			sample += g.phaser_buffer[(g.ipp-g.iphase+1024)&1023]
			g.ipp = (g.ipp + 1) & 1023
			// final accumulation and envelope application
			ssample += sample * g.env_vol
		}
		ssample = ssample / 8 * masterVolume
		ssample *= 2.0 * g.cfg.Volume

		if ssample > 1.0 {
			ssample = 1.0
		} else if ssample < -1.0 {
			ssample = -1.0
		}
		buffer = append(buffer, ssample)
	}
	return buffer
}

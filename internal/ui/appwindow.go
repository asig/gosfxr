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
package ui

import (
	"fmt"
	"io/ioutil"
	"log"
	"path/filepath"
	"time"

	"github.com/gotk3/gotk3/gdk"
	"github.com/gotk3/gotk3/glib"
	"github.com/gotk3/gotk3/gtk"
	"github.com/veandco/go-sdl2/mix"

	"github.com/asig/gosfxr/internal/generator"
	"github.com/asig/gosfxr/internal/resources"
	"github.com/asig/gosfxr/internal/wav"
)

type AppWindow struct {
	generatorConfig *generator.Config
	generatedSample []float64

	gtkWindow *gtk.ApplicationWindow

	// Controls
	btnWaveform             map[generator.Waveform]*gtk.RadioButton
	adjVolume               *gtk.Adjustment
	adjEnvelopeAttack       *gtk.Adjustment
	adjEnvelopeSustain      *gtk.Adjustment
	adjEnvelopeSustainPunch *gtk.Adjustment
	adjEnvelopeDecay        *gtk.Adjustment
	adFreqStart             *gtk.Adjustment
	adFreqMinCutoff         *gtk.Adjustment
	adFreqSlide             *gtk.Adjustment
	adFreqDeltaSlide        *gtk.Adjustment
	adjVibratoDepth         *gtk.Adjustment
	adjVibratoSpeed         *gtk.Adjustment
	adjArpFreqMult          *gtk.Adjustment
	adjArpChangeSpeed       *gtk.Adjustment
	adjDutyCycleCycle       *gtk.Adjustment
	adjDutyCycleSweep       *gtk.Adjustment
	adjRepeatRate           *gtk.Adjustment
	adjPhaserOffset         *gtk.Adjustment
	adjPhaserSweep          *gtk.Adjustment
	adjLPCutoffFreq         *gtk.Adjustment
	adjLPCutoffSweep        *gtk.Adjustment
	adjLPResonance          *gtk.Adjustment
	adjHPCutoffFreq         *gtk.Adjustment
	adjHPCutoffSweep        *gtk.Adjustment
	comboExportFreq         *gtk.ComboBox
	comboExportBits         *gtk.ComboBox
	imgGeneratedSample      *gtk.Image
	statusbar               *gtk.Statusbar
	nextStatusMsgId         int

	updating bool
}

func createStatusbar() *gtk.Statusbar {
	statusbar, _ := gtk.StatusbarNew()
	return statusbar
}

func getObj(builder *gtk.Builder, name string) glib.IObject {
	obj, err := builder.GetObject(name)
	if err != nil {
		log.Fatal(err)
	}
	return obj
}

func loadPixbufFromResource(name string) *gdk.Pixbuf {
	pixbuf, err := gdk.PixbufNewFromDataOnly(resources.Find(name))
	if err != nil {
		panic(err)
	}
	return pixbuf
}

func loadImageFromPixbuf(pixbuf *gdk.Pixbuf) *gtk.Image {
	img, err := gtk.ImageNewFromPixbuf(pixbuf)
	if err != nil {
		panic(err)
	}
	//pixbuf.Unref()
	return img
}

func setImage(builder *gtk.Builder, ctrlName, imgName string) {
	getObj(builder, ctrlName).(*gtk.Image).SetFromPixbuf(loadPixbufFromResource(imgName))
}

func setBtnImage(builder *gtk.Builder, btnName, imgName string) {
	getObj(builder, btnName).(*gtk.Button).SetImage(loadImageFromPixbuf(loadPixbufFromResource(imgName)))
}

func setRadioButtonImage(btn *gtk.RadioButton, imgName string) {
	btn.SetImage(loadImageFromPixbuf(loadPixbufFromResource(imgName)))
}

func (a *AppWindow) applyPreset(presetFunc func()) {
	presetFunc()
	a.updateControls()
	a.play()
}

func (a *AppWindow) toggleWave(btn *gtk.RadioButton, wf generator.Waveform) {
	if !btn.GetActive() {
		return
	}
	a.generatorConfig.Waveform = wf
	a.updateControls()
}

func NewAppWindow(a *gtk.Application, cfg *generator.Config) *AppWindow {
	appWindow := &AppWindow{
		generatorConfig: cfg,
	}

	builder, _ := gtk.BuilderNew()
	builder.AddFromString(uiXMLString)

	builder.ConnectSignals(map[string]interface{}{
		"btn_waveform_sine_toggled_cb":     func(btn *gtk.RadioButton) { appWindow.toggleWave(btn, generator.WaveformSine) },
		"btn_waveform_square_toggled_cb":   func(btn *gtk.RadioButton) { appWindow.toggleWave(btn, generator.WaveformSquare) },
		"btn_waveform_sawtooth_toggled_cb": func(btn *gtk.RadioButton) { appWindow.toggleWave(btn, generator.WaveformSawtooth) },
		"btn_waveform_noise_toggled_cb":    func(btn *gtk.RadioButton) { appWindow.toggleWave(btn, generator.WaveformNoise) },

		// Presets
		"btn_pickup_clicked_cb":    func() { appWindow.applyPreset(appWindow.generatorConfig.PresetPickup) },
		"btn_laser_clicked_cb":     func() { appWindow.applyPreset(appWindow.generatorConfig.PresetLaser) },
		"btn_explosion_clicked_cb": func() { appWindow.applyPreset(appWindow.generatorConfig.PresetExplosion) },
		"btn_powerup_clicked_cb":   func() { appWindow.applyPreset(appWindow.generatorConfig.PresetPowerup) },
		"btn_hit_clicked_cb":       func() { appWindow.applyPreset(appWindow.generatorConfig.PresetHit) },
		"btn_jump_clicked_cb":      func() { appWindow.applyPreset(appWindow.generatorConfig.PresetJump) },
		"btn_blip_clicked_cb":      func() { appWindow.applyPreset(appWindow.generatorConfig.PresetBlip) },

		// Automated adjustments
		"btn_mutate_clicked_cb":    func() { appWindow.generatorConfig.Mutate(); appWindow.updateControls() },
		"btn_randomize_clicked_cb": func() { appWindow.generatorConfig.Randomize(); appWindow.updateControls() },

		// Other buttons
		"btn_play_clicked_cb":   func() { appWindow.play() },
		"btn_load_clicked_cb":   func() { appWindow.load() },
		"btn_save_clicked_cb":   func() { appWindow.save() },
		"btn_export_clicked_cb": func() { appWindow.export() },

		// Sliders
		"adj_arp_changespeed_value_changed_cb":        func(adj *gtk.Adjustment) { appWindow.generatorConfig.ArpChangeSpeed = adj.GetValue(); appWindow.updateControls() },
		"adj_arp_freqmult_value_changed_cb":           func(adj *gtk.Adjustment) { appWindow.generatorConfig.ArpFreqMult = adj.GetValue(); appWindow.updateControls() },
		"adj_dutycycle_cycle_value_changed_cb":        func(adj *gtk.Adjustment) { appWindow.generatorConfig.DutyCycle = adj.GetValue(); appWindow.updateControls() },
		"adj_dutycycle_sweep_value_changed_cb":        func(adj *gtk.Adjustment) { appWindow.generatorConfig.DutyCycleSweep = adj.GetValue(); appWindow.updateControls() },
		"adj_envelope_attack_value_changed_cb":        func(adj *gtk.Adjustment) { appWindow.generatorConfig.EnvelopeAttack = adj.GetValue(); appWindow.updateControls() },
		"adj_envelope_decay_value_changed_cb":         func(adj *gtk.Adjustment) { appWindow.generatorConfig.EnvelopeDecay = adj.GetValue(); appWindow.updateControls() },
		"adj_envelope_sustain_value_changed_cb":       func(adj *gtk.Adjustment) { appWindow.generatorConfig.EnvelopeSustain = adj.GetValue(); appWindow.updateControls() },
		"adj_envelope_sustain_punch_value_changed_cb": func(adj *gtk.Adjustment) { appWindow.generatorConfig.EnvelopeSustainPunch = adj.GetValue(); appWindow.updateControls() },
		"adj_freq_delta_slide_value_changed_cb":       func(adj *gtk.Adjustment) { appWindow.generatorConfig.FreqDeltaSlide = adj.GetValue(); appWindow.updateControls() },
		"adj_freq_min_cutoff_value_changed_cb":        func(adj *gtk.Adjustment) { appWindow.generatorConfig.FreqMinCutoff = adj.GetValue(); appWindow.updateControls() },
		"adj_freq_slide_value_changed_cb":             func(adj *gtk.Adjustment) { appWindow.generatorConfig.FreqSlide = adj.GetValue(); appWindow.updateControls() },
		"adj_freq_start_value_changed_cb":             func(adj *gtk.Adjustment) { appWindow.generatorConfig.FreqStart = adj.GetValue(); appWindow.updateControls() },
		"adj_hp_cutoff_freq_value_changed_cb":         func(adj *gtk.Adjustment) { appWindow.generatorConfig.HPCutoffFreq = adj.GetValue(); appWindow.updateControls() },
		"adj_hp_cutoff_sweep_value_changed_cb":        func(adj *gtk.Adjustment) { appWindow.generatorConfig.HPCutoffSweep = adj.GetValue(); appWindow.updateControls() },
		"adj_lp_cutoff_freq_value_changed_cb":         func(adj *gtk.Adjustment) { appWindow.generatorConfig.LPCutoffFreq = adj.GetValue(); appWindow.updateControls() },
		"adj_lp_cutoff_sweep_value_changed_cb":        func(adj *gtk.Adjustment) { appWindow.generatorConfig.LPCutoffSweep = adj.GetValue(); appWindow.updateControls() },
		"adj_lp_resonance_value_changed_cb":           func(adj *gtk.Adjustment) { appWindow.generatorConfig.LPResonance = adj.GetValue(); appWindow.updateControls() },
		"adj_phaser_offset_value_changed_cb":          func(adj *gtk.Adjustment) { appWindow.generatorConfig.PhaserOffset = adj.GetValue(); appWindow.updateControls() },
		"adj_phaser_sweep_value_changed_cb":           func(adj *gtk.Adjustment) { appWindow.generatorConfig.PhaserSweep = adj.GetValue(); appWindow.updateControls() },
		"adj_repeat_rate_value_changed_cb":            func(adj *gtk.Adjustment) { appWindow.generatorConfig.RepeatRate = adj.GetValue(); appWindow.updateControls() },
		"adj_vibrato_depth_value_changed_cb":          func(adj *gtk.Adjustment) { appWindow.generatorConfig.VibDepth = adj.GetValue(); appWindow.updateControls() },
		"adj_vibrato_speed_value_changed_cb":          func(adj *gtk.Adjustment) { appWindow.generatorConfig.VibSpeed = adj.GetValue(); appWindow.updateControls() },
		"adj_volume_value_changed_cb":                 func(adj *gtk.Adjustment) { appWindow.generatorConfig.Volume = adj.GetValue(); appWindow.updateControls() },
	})

	appWindow.gtkWindow = getObj(builder, "application_window").(*gtk.ApplicationWindow)

	appWindow.imgGeneratedSample = getObj(builder, "img_generated_sample").(*gtk.Image)

	// Controls
	appWindow.btnWaveform = map[generator.Waveform]*gtk.RadioButton{
		generator.WaveformSquare:   getObj(builder, "btn_waveform_square").(*gtk.RadioButton),
		generator.WaveformSawtooth: getObj(builder, "btn_waveform_sawtooth").(*gtk.RadioButton),
		generator.WaveformSine:     getObj(builder, "btn_waveform_sine").(*gtk.RadioButton),
		generator.WaveformNoise:    getObj(builder, "btn_waveform_noise").(*gtk.RadioButton),
	}
	appWindow.adjVolume = getObj(builder, "adj_volume").(*gtk.Adjustment)
	appWindow.adjEnvelopeAttack = getObj(builder, "adj_envelope_attack").(*gtk.Adjustment)
	appWindow.adjEnvelopeSustain = getObj(builder, "adj_envelope_sustain").(*gtk.Adjustment)
	appWindow.adjEnvelopeSustainPunch = getObj(builder, "adj_envelope_sustain_punch").(*gtk.Adjustment)
	appWindow.adjEnvelopeDecay = getObj(builder, "adj_envelope_decay").(*gtk.Adjustment)
	appWindow.adFreqStart = getObj(builder, "adj_freq_start").(*gtk.Adjustment)
	appWindow.adFreqMinCutoff = getObj(builder, "adj_freq_min_cutoff").(*gtk.Adjustment)
	appWindow.adFreqSlide = getObj(builder, "adj_freq_slide").(*gtk.Adjustment)
	appWindow.adFreqDeltaSlide = getObj(builder, "adj_freq_delta_slide").(*gtk.Adjustment)
	appWindow.adjVibratoDepth = getObj(builder, "adj_vibrato_depth").(*gtk.Adjustment)
	appWindow.adjVibratoSpeed = getObj(builder, "adj_vibrato_speed").(*gtk.Adjustment)
	appWindow.adjArpFreqMult = getObj(builder, "adj_arp_freqmult").(*gtk.Adjustment)
	appWindow.adjArpChangeSpeed = getObj(builder, "adj_arp_changespeed").(*gtk.Adjustment)
	appWindow.adjDutyCycleCycle = getObj(builder, "adj_dutycycle_cycle").(*gtk.Adjustment)
	appWindow.adjDutyCycleSweep = getObj(builder, "adj_dutycycle_sweep").(*gtk.Adjustment)
	appWindow.adjRepeatRate = getObj(builder, "adj_repeat_rate").(*gtk.Adjustment)
	appWindow.adjPhaserOffset = getObj(builder, "adj_phaser_offset").(*gtk.Adjustment)
	appWindow.adjPhaserSweep = getObj(builder, "adj_phaser_sweep").(*gtk.Adjustment)
	appWindow.adjLPCutoffFreq = getObj(builder, "adj_lp_cutoff_freq").(*gtk.Adjustment)
	appWindow.adjLPCutoffSweep = getObj(builder, "adj_lp_cutoff_sweep").(*gtk.Adjustment)
	appWindow.adjLPResonance = getObj(builder, "adj_lp_resonance").(*gtk.Adjustment)
	appWindow.adjHPCutoffFreq = getObj(builder, "adj_hp_cutoff_freq").(*gtk.Adjustment)
	appWindow.adjHPCutoffSweep = getObj(builder, "adj_hp_cutoff_sweep").(*gtk.Adjustment)
	appWindow.comboExportFreq = getObj(builder, "combo_frequency").(*gtk.ComboBox)
	appWindow.comboExportBits = getObj(builder, "combo_bits").(*gtk.ComboBox)
	appWindow.statusbar = getObj(builder, "statusbar").(*gtk.Statusbar)

	// Set images
	appWindow.btnWaveform[generator.WaveformSquare].SetImage(loadImageFromPixbuf(loadPixbufFromResource("resources/waveforms/waveform_square.png")))
	appWindow.btnWaveform[generator.WaveformSawtooth].SetImage(loadImageFromPixbuf(loadPixbufFromResource("resources/waveforms/waveform_sawtooth.png")))
	appWindow.btnWaveform[generator.WaveformSine].SetImage(loadImageFromPixbuf(loadPixbufFromResource("resources/waveforms/waveform_sine.png")))
	appWindow.btnWaveform[generator.WaveformNoise].SetImage(loadImageFromPixbuf(loadPixbufFromResource("resources/waveforms/waveform_noise.png")))

	setBtnImage(builder, "btn_play", "resources/16px/play-solid.png")
	setBtnImage(builder, "btn_load", "resources/16px/upload-solid.png")
	setBtnImage(builder, "btn_save", "resources/16px/download-solid.png")
	setBtnImage(builder, "btn_export", "resources/16px/file-export-solid.png")

	setImage(builder, "img_volume", "resources/16px/volume-high-solid.png")

	appWindow.comboExportFreq.SetActive(0)
	appWindow.comboExportBits.SetActive(1)

	pb, err := gdk.PixbufNewFromBytesOnly(resources.Find("resources/icons/icon.svg"))
	if err != nil {
		panic(err)
	}
	appWindow.gtkWindow.SetIcon(pb)
	appWindow.gtkWindow.SetTitle("gosfxr")
	appWindow.gtkWindow.SetDefaultSize(800, 800)

	return appWindow
}

func (a *AppWindow) Show() {
	a.gtkWindow.ShowAll()
	a.updateControls()
}

func (a *AppWindow) GtkWindow() *gtk.ApplicationWindow {
	return a.gtkWindow
}

func (a *AppWindow) setStatus(msg string) {
	a.statusbar.Push(0, msg)
	a.nextStatusMsgId++
	curStatusId := a.nextStatusMsgId
	go func() {
		time.Sleep(5 * time.Second)
		if curStatusId == a.nextStatusMsgId {
			// Still the same message, remove if. Otherwise, somebody else will clear it
			a.statusbar.RemoveAll(0)
		}
	}()
}

func (a *AppWindow) play() {
	wavData := wav.Generate(a.generatedSample, 16, 44100)
	chunk, _ := mix.QuickLoadWAV(wavData)
	chunk.Play(-1, 0)
}

func makeFilter(name string, patterns ...string) *gtk.FileFilter {
	filter, _ := gtk.FileFilterNew()
	filter.SetName(name)
	for _, p := range patterns {
		filter.AddPattern(p)
	}
	return filter
}

func fixExtensions(filename, ext string) string {
	if len(filepath.Ext(filename)) > 0 {
		// Already an extensions, we're good.
		return filename
	}
	return filename + ext
}

func (a *AppWindow) fileDialog(title string, action gtk.FileChooserAction, okBtnTitle string, filter *gtk.FileFilter) (string, bool) {
	dlg, _ := gtk.FileChooserDialogNewWith2Buttons(title, a.gtkWindow, action, okBtnTitle, gtk.RESPONSE_ACCEPT, "Cancel", gtk.RESPONSE_CANCEL)
	defer dlg.Destroy()

	dlg.AddFilter(filter)
	dlg.AddFilter(makeFilter("All files", "*.*"))
	res := dlg.Run()
	filename := dlg.GetFilename()
	return filename, res == gtk.RESPONSE_ACCEPT
}

func (a *AppWindow) load() {
	filename, ok := a.fileDialog("Load configuration", gtk.FILE_CHOOSER_ACTION_OPEN, "Load", makeFilter("Configs", "*.json"))
	if !ok {
		return
	}

	content, _ := ioutil.ReadFile(filename)
	a.generatorConfig.InitFromJson(content)
	a.updateControls()
	a.setStatus(fmt.Sprintf("Configuration read from %s.", filename))

}

func (a *AppWindow) save() {
	filename, ok := a.fileDialog("Save configuration", gtk.FILE_CHOOSER_ACTION_SAVE, "Save", makeFilter("Configs", "*.json"))
	if !ok {
		return
	}
	filename = fixExtensions(filename, ".json")
	ioutil.WriteFile(filename, a.generatorConfig.ToJson(), 0644)
	a.setStatus(fmt.Sprintf("Configuration written to %s.", filename))
}

func getComboInt(combo *gtk.ComboBox) int {
	iter, _ := combo.GetActiveIter()
	tm, _ := combo.GetModel()
	val, _ := tm.ToTreeModel().GetValue(iter, 0)
	gv, _ := val.GoValue()
	return gv.(int)
}

func (a *AppWindow) export() {
	filename, ok := a.fileDialog("Export to WAV", gtk.FILE_CHOOSER_ACTION_SAVE, "Export", makeFilter("WAV files", "*.wav"))
	if !ok {
		return
	}
	filename = fixExtensions(filename, ".wav")

	bits := getComboInt(a.comboExportBits)
	freq := getComboInt(a.comboExportFreq)
	wav := wav.Generate(a.generatedSample, bits, freq)
	ioutil.WriteFile(filename, wav, 0644)
	a.setStatus(fmt.Sprintf("WAV exported to %s.", filename))
}

func (a *AppWindow) updateGeneratedSampleImage(sample []float64) {
	alloc := a.imgGeneratedSample.GetAllocation()
	w := alloc.GetWidth()
	h := alloc.GetHeight()

	imgdata := make([]byte, w*h*3, w*h*3)
	stride := w * 3
	for y := 0; y < h; y++ {
		for x := 0; x < w; x++ {
			imgdata[y*stride+3*x+0] = byte(x & 0xff)
			imgdata[y*stride+3*x+1] = byte(y & 0xff)
			imgdata[y*stride+3*x+2] = 0
		}
	}

	for idx, _ := range imgdata {
		imgdata[idx] = 255
	}
	bufstep := float64(len(sample)) / float64(w)
	for x := 0; x < w; x++ {
		// Interpolate all values
		p1 := int(float64(x) * bufstep)
		p2 := int(float64(x+1) * bufstep)
		d := p2 - p1
		s := 0.0
		for p1 < p2 {
			s = s + sample[p1]
			p1++
		}
		s /= float64(d)

		y := int((s + 1) * float64(h) / 2.0)
		imgdata[y*stride+3*x+0] = 0
		imgdata[y*stride+3*x+1] = 0
		imgdata[y*stride+3*x+2] = 0
	}
	pixbuf, err := gdk.PixbufNewFromData(imgdata, gdk.COLORSPACE_RGB, false, 8, w, h, stride)
	if err != nil {
		panic(err)
	}

	a.imgGeneratedSample.SetFromPixbuf(pixbuf)
}

func (a *AppWindow) updateControls() {
	if a.updating {
		return
	}
	a.updating = true

	w := a.generatorConfig.Waveform
	for wf, b := range a.btnWaveform {
		b.SetActive(wf == w)
	}

	a.adjVolume.SetValue(a.generatorConfig.Volume)
	a.adjEnvelopeAttack.SetValue(a.generatorConfig.EnvelopeAttack)
	a.adjEnvelopeSustain.SetValue(a.generatorConfig.EnvelopeSustain)
	a.adjEnvelopeSustainPunch.SetValue(a.generatorConfig.EnvelopeSustainPunch)
	a.adjEnvelopeDecay.SetValue(a.generatorConfig.EnvelopeDecay)
	a.adFreqStart.SetValue(a.generatorConfig.FreqStart)
	a.adFreqMinCutoff.SetValue(a.generatorConfig.FreqMinCutoff)
	a.adFreqSlide.SetValue(a.generatorConfig.FreqSlide)
	a.adFreqDeltaSlide.SetValue(a.generatorConfig.FreqDeltaSlide)
	a.adjVibratoDepth.SetValue(a.generatorConfig.VibDepth)
	a.adjVibratoSpeed.SetValue(a.generatorConfig.VibSpeed)
	a.adjArpFreqMult.SetValue(a.generatorConfig.ArpFreqMult)
	a.adjArpChangeSpeed.SetValue(a.generatorConfig.ArpChangeSpeed)
	a.adjDutyCycleCycle.SetValue(a.generatorConfig.DutyCycle)
	a.adjDutyCycleSweep.SetValue(a.generatorConfig.DutyCycleSweep)
	a.adjRepeatRate.SetValue(a.generatorConfig.RepeatRate)
	a.adjPhaserOffset.SetValue(a.generatorConfig.PhaserOffset)
	a.adjPhaserSweep.SetValue(a.generatorConfig.PhaserSweep)
	a.adjLPCutoffFreq.SetValue(a.generatorConfig.LPCutoffFreq)
	a.adjLPCutoffSweep.SetValue(a.generatorConfig.LPCutoffSweep)
	a.adjLPResonance.SetValue(a.generatorConfig.LPResonance)
	a.adjHPCutoffFreq.SetValue(a.generatorConfig.HPCutoffFreq)
	a.adjHPCutoffSweep.SetValue(a.generatorConfig.HPCutoffSweep)

	sample := generator.New(a.generatorConfig).Generate()
	a.generatedSample = sample
	a.updateGeneratedSampleImage(sample)

	a.updating = false
}

#
# Copyright (c) 2021 Andreas Signer <asigner@gmail.com>
#
# This file is part of gosfxr.
#
# gosfxr is free software: you can redistribute it and/or
# modify it under the terms of the GNU General Public License as
# published by the Free Software Foundation, either version 3 of the
# License, or (at your option) any later version.
#
# gosfxr is distributed in the hope that it will be useful,
# but WITHOUT ANY WARRANTY; without even the implied warranty of
# MERCHANTABILITY or FITNESS FOR A PARTICULAR PURPOSE.  See the
# GNU General Public License for more details.
#
# You should have received a copy of the GNU General Public License
# along with gosfxr.  If not, see <http://www.gnu.org/licenses/>.
#

.PHONY: clean

RESOURCES = \
    resources/icons/icon.svg \

GENERATED_RESOURCES = \
    resources/16px/download-solid.png \
    resources/16px/file-export-solid.png \
    resources/16px/play-solid.png \
    resources/16px/upload-solid.png \
    resources/16px/volume-high-solid.png \
    resources/32px/download-solid.png \
    resources/32px/file-export-solid.png \
    resources/32px/play-solid.png \
    resources/32px/upload-solid.png \
    resources/32px/volume-high-solid.png \
    resources/48px/download-solid.png \
    resources/48px/file-export-solid.png \
    resources/48px/play-solid.png \
    resources/48px/upload-solid.png \
    resources/48px/volume-high-solid.png \
    resources/64px/download-solid.png \
    resources/64px/file-export-solid.png \
    resources/64px/play-solid.png \
    resources/64px/upload-solid.png \
    resources/64px/volume-high-solid.png \
    resources/waveforms/waveform_sine.png \
    resources/waveforms/waveform_sawtooth.png \
    resources/waveforms/waveform_square.png \
    resources/waveforms/waveform_noise.png \


all: internal/resources/resources.go internal/ui/ui_resources.go
	go build

clean:
	@rm -f \
      gosfxr \
      bin2go \
      internal/ui/ui_resources.go \
      internal/resources/resources.go \
      ${GENERATED_RESOURCES}

run:	all
	./gosfxr

internal/ui/ui_resources.go: gosfxr.ui
	echo "package ui;" > internal/ui/ui_resources.go
	echo "const uiXMLString = \`" >> internal/ui/ui_resources.go
	cat gosfxr.ui >> internal/ui/ui_resources.go
	echo "\`" >> internal/ui/ui_resources.go

internal/resources/resources.go: bin2go $(RESOURCES) ${GENERATED_RESOURCES}
	@mkdir -p internal/resources
	./bin2go -pkg resources -out internal/resources/resources.go $(RESOURCES) ${GENERATED_RESOURCES}

bin2go:	tools/bin2go/bin2go.go
	go build tools/bin2go/bin2go.go

resources/16px/%.png: resources/%.svg
	@mkdir -p resources/16px
	inkscape -z -w 16 -h 16 -o $@ $<

resources/32px/%.png: resources/%.svg
	@mkdir -p resources/32px
	inkscape -z -w 32 -h 32 -o $@ $<

resources/48px/%.png: resources/%.svg
	@mkdir -p resources/48px
	inkscape -z -w 48 -h 48 -o $@ $<

resources/64px/%.png: resources/%.svg
	@mkdir -p resources/64px
	inkscape -z -w 64 -h 64 -o $@ $<

resources/waveforms/%.png: resources/%.svg
	@mkdir -p resources/waveforms
	inkscape -z -w 92 -h 24 -o $@ $<


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
package main

import (
	"github.com/veandco/go-sdl2/mix"
	"github.com/veandco/go-sdl2/sdl"

	"github.com/asig/gosfxr/internal/app"
)

func main() {
	err := sdl.Init(sdl.INIT_AUDIO)
	if err != nil {
		panic(err)
	}
	defer sdl.Quit()

	err = mix.Init(0) // We'll only play samples
	if err != nil {
		panic(err)
	}
	defer mix.Quit()

	err = mix.OpenAudio(mix.DEFAULT_FREQUENCY, mix.DEFAULT_FORMAT, mix.DEFAULT_CHANNELS, 4096)
	if err != nil {
		panic(err)
	}
	defer mix.CloseAudio()

	application := app.New()
	application.Run()
}


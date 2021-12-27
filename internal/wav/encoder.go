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
package wav

import (
	"bytes"
	"encoding/binary"
	"fmt"
)

func toUint32(val uint32) []byte {
	buf := make([]byte, 4)
	binary.LittleEndian.PutUint32(buf, val)
	return buf
}

func toUint16(val uint16) []byte {
	buf := make([]byte, 2)
	binary.LittleEndian.PutUint16(buf, val)
	return buf
}

func Generate(data []float64, bits int, freq int) []byte {
	// We're cheap and only support 44100 Hz and 22050 Hz, with 8 or 16 bits.
	if freq != 44100 && freq != 22050 {
		panic(fmt.Errorf("Unsupported frequency %d", freq))
	}
	if bits != 8 && bits != 16 {
		panic(fmt.Errorf("Unsupported bit depth %d", bits))
	}

	var resampledData []float64
	if freq == 44100 {
		// No resampling needed
		resampledData = data
	} else {
		resampledData = make([]float64,len(data)/2);
		for i := 0; i < len(data)/2; i++ {
			resampledData[i] = (data[2*i+0] + data[2*i + 1])/2
		}
	}

	var wavData []byte

	if bits == 8 {
		wavData = make([]byte, len(resampledData))
		for i := 0; i < len(resampledData); i++ {
			wavData[i] = byte((resampledData[i]+1)/2*255)
		}
	} else {
		wavData = make([]byte, 2*len(resampledData))
		for i := 0; i < len(resampledData); i++ {
			val := toUint16(uint16(resampledData[i] * 32767))
			wavData[2*i+0] = val[0]
			wavData[2*i+1] = val[1]
		}
	}


	buf := bytes.NewBuffer(make([]byte,0,len(data)+100))

	// Write WAV header
	buf.Write(toUint32(0x46464952)) // "RIFF"
	buf.Write(toUint32(uint32(4+24+8+len(wavData))))
	buf.Write(toUint32(0x45564157)) // "WAVE"

	buf.Write(toUint32(0x20746d66)) // "fmt "
	buf.Write(toUint32(16)  )       // size remaining header
	buf.Write(toUint16(1)    )      // PCM format
	buf.Write(toUint16(1)     )     // # of channels
	buf.Write(toUint32(uint32(freq)        )) // SampleRate
	buf.Write(toUint32(uint32(freq*1*int(bits/8)))) // ByteRate
	buf.Write(toUint16(uint16(1*int(bits/8)))) // BlockAlign
	buf.Write(toUint16(uint16(bits)))

	buf.Write(toUint32(0x61746164)) // "data"
	buf.Write(toUint32(uint32(len(wavData))))
	buf.Write(wavData)

	return buf.Bytes()
}

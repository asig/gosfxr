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
	"flag"
	"fmt"
	"io/ioutil"
	"os"
	"strings"
)
var(
	flagOut = flag.String("out", "", "Destination file")
	flagPkg = flag.String("pkg", "data", "Package to be used.")

	out *os.File
	resources = make(map[string][]byte)
)



func handleFile(f string) error {
	data, err := ioutil.ReadFile(f)
	if err != nil {
		return err
	}
	resources[f] = data
	return nil
}


func toByteArray(bytes []byte) []string {
	var lines []string
	linelen := 18
	start := 0
	for start < len(bytes) {
		var s []string
		for i := 0; i < linelen; i++ {
			if start + i < len(bytes) {
				s = append(s, fmt.Sprintf("0x%02x", bytes[start+i]))
			}
		}
		lines = append(lines, strings.Join(s, ", "))
		start = start + linelen
	}
	return lines
}

func main() {
	flag.Parse()

	if *flagOut == "" {
		flag.Usage()
		return
	}

	for _, f := range flag.Args() {
		err := handleFile(f)
		if err != nil {
			fmt.Printf("Error while processing file %q: %s\n", f, err)
		}
	}

	var err error
	out, err = os.Create(*flagOut)
	if err != nil {
		fmt.Printf("Can't create output file %q: %s", *flagOut, err)
		return
	}

	fmt.Fprintf(out, "package %s\n", *flagPkg)
	fmt.Fprintf(out, "\n")
	fmt.Fprintf(out, "import \"log\"\n")
	fmt.Fprintf(out, "\n")
	fmt.Fprintf(out, "var resources = map[string][]byte {\n")
	for key, val := range resources {
		fmt.Fprintf(out, "\t%q: {\n", key)
		for _, l := range toByteArray(val) {
			fmt.Fprintf(out, "\t\t%s,\n", l)
		}
		fmt.Fprintf(out, "\t},\n")
	}
	fmt.Fprintf(out, "}\n")

	fmt.Fprint(out, `
func Find(name string) []byte {
	if res, ok := resources[name]; ok {
		return res
	}
	log.Fatalf("Resource %s not found", name)
	return nil
}
`)
}


/* This strips out all of the ATTLIST headings and all of the
   nonstandard attributes associated with KanjiVG for the sake of the
   people whose XML parsers don't process the KanjiVG SVG files
   correctly. */

package main

import (
	"kvg"
	"log"
	"os"
	"path/filepath"
	"regexp"
)

// This is the heading for our files.

var StrippedHeading = kvg.Heading

// This removes all the ATTLIST things from the header.

var attlistRe = regexp.MustCompile(`(?ms)\s*(\[\s*<\!ATTLIST.*?>\s*)+\]`)

// This should be an option in theory but for now we'll use this
// directory.

var stripDir = "/home/ben/software/kanjivg/stripped"

func strip(file string) {
	svg := kvg.ReadKanjiFileOrDie(file)
	base := svg.BaseGroup()
	paths := base.GetPaths()
	groups := base.GetGroups()
	for i, g := range groups {
		g.Element = ""
		g.Variant = false
		g.Partial = false
		g.Original = ""
		g.Part = ""
		g.Number = ""
		g.TradForm = ""
		g.RadicalForm = ""
		g.Position = ""
		g.Radical = ""
		g.Phon = ""
		groups[i] = g
	}
	for j, p := range paths {
		p.Type = ""
		paths[j] = p
	}
	file = filepath.Base(file)
	svg.WriteKanjiFile(stripDir + "/" + file)
}

func main() {
	StrippedHeading = attlistRe.ReplaceAllString(kvg.Heading, "")
	_, err := os.Stat(stripDir)
	if os.IsNotExist(err) {
		err := os.Mkdir(stripDir, 0755)
		if err != nil {
			log.Fatalf("Could not make %s: %s", stripDir, err)
		}
	}
	kvg.Heading = StrippedHeading
	kvg.ExamineAllFilesSimple(strip)
}

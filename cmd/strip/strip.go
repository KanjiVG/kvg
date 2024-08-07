package main

import (
	"kvg"
	"log"
	"os"
	"path/filepath"
	"regexp"
)

var StrippedHeading = kvg.Heading

var attlistRe = regexp.MustCompile(`(?ms)\s*(\[\s*<\!ATTLIST.*?>\s*)+\]`)

//var attlistRe = regexp.MustCompile(`(?ms)\[(<!ATTLIST.*?>\s*)+\]`)
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

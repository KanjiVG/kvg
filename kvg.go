// Read, write, and manipulate KanjiVG files
package kvg

import (
	"encoding/xml"
	"fmt"
	"io/fs"
	"log"
	"os"
	"path/filepath"
	"regexp"
	"strconv"
	"strings"
	"unicode"
)

// This matches most of the variant endings.
var Variant = regexp.MustCompile(`-(Kaisho|MidFst|HzLst|VtLst|HzFstLeRi|HzFstRiLe|TenLst|Hyougai|Jinmei|HzFst|VtFstRiLe|LeFst|Vt6|VtFstRiLe|HzFst|HzFstVtLst|MdLst|VtFst|Vt4|Ten3|DgLst|Insatsu|MdFst|MdFst2|Dg3|TenFst|RiLe|NoDot)`)

// The base directory
var KVDir = "/home/ben/software/kanjivg/kanji"

// A path, in other words a stroke of the kanji.
type Path struct {
	XMLName xml.Name `xml:"path"`
	ID      string   `xml:"id,attr"`
	Type    string   `xml:"kvg:type,attr,omitempty"`
	D       string   `xml:"d,attr"`
	// The node in the tree of either groups or paths which
	// corresponds to this path.
	Parent *Child `xml:"-"`
	Class  string `xml:"class,attr,omitempty"`
}

// Text holder, this contains the stroke numbers.
type Text struct {
	XMLName   xml.Name `xml:"text"`
	Transform string   `xml:"transform,attr,omitempty"`
	Content   []byte   `xml:",chardata"`
	Parent    *Child   `xml:"-"`
	Class     string   `xml:"class,attr,omitempty"`
}

// Either a group or a path element.
type Child struct {
	Path    Path
	Group   Group
	Text    Text
	IsGroup bool
	IsText  bool
	// The parent of this element in the tree of groups.
	Parent *Group `xml:"-"`
}

// A group.
type Group struct {
	XMLName     xml.Name `xml:"g"`
	ID          string   `xml:"id,attr,omitempty"`
	Element     string   `xml:"kvg:element,attr,omitempty"`
	Part        string   `xml:"kvg:part,attr,omitempty"`
	Variant     bool     `xml:"kvg:variant,attr,omitempty"`
	Number      string   `xml:"kvg:number,attr,omitempty"`
	Original    string   `xml:"kvg:original,attr,omitempty"`
	Partial     bool     `xml:"kvg:partial,attr,omitempty"`
	TradForm    string   `xml:"kvg:tradForm,attr,omitempty"`
	Position    string   `xml:"kvg:position,attr,omitempty"`
	Radical     string   `xml:"kvg:radical,attr,omitempty"`
	Phon        string   `xml:"kvg:phon,attr,omitempty"`
	RadicalForm string   `xml:"kvg:radicalForm,attr,omitempty"`
	Style       string   `xml:"style,attr,omitempty"`
	Children    []Child
	Parent      *Child `xml:"-"`
}

// An entire file.
type SVG struct {
	XMLName xml.Name `xml:"svg"`
	XMLNS   string   `xml:"xmlns,attr"`
	Width   string   `xml:"width,attr"`
	Height  string   `xml:"height,attr"`
	ViewBox string   `xml:"viewBox,attr,omitempty"`
	Groups  []Group  `xml:"g"`
}

// This is the heading as repeated in each file.
var Heading = `<?xml version="1.0" encoding="UTF-8"?>
<!--
Copyright (C) 2009/2010/2011 Ulrich Apel.
This work is distributed under the conditions of the Creative Commons
Attribution-Share Alike 3.0 Licence. This means you are free:
* to Share - to copy, distribute and transmit the work
* to Remix - to adapt the work

Under the following conditions:
* Attribution. You must attribute the work by stating your use of KanjiVG in
  your own copyright header and linking to KanjiVG's website
  (http://kanjivg.tagaini.net)
* Share Alike. If you alter, transform, or build upon this work, you may
  distribute the resulting work only under the same or similar license to this
  one.

See http://creativecommons.org/licenses/by-sa/3.0/ for more details.
-->
<!DOCTYPE svg PUBLIC "-//W3C//DTD SVG 1.0//EN" "http://www.w3.org/TR/2001/REC-SVG-20010904/DTD/svg10.dtd" [
<!ATTLIST g
xmlns:kvg CDATA #FIXED "http://kanjivg.tagaini.net"
kvg:element CDATA #IMPLIED
kvg:variant CDATA #IMPLIED
kvg:partial CDATA #IMPLIED
kvg:original CDATA #IMPLIED
kvg:part CDATA #IMPLIED
kvg:number CDATA #IMPLIED
kvg:tradForm CDATA #IMPLIED
kvg:radicalForm CDATA #IMPLIED
kvg:position CDATA #IMPLIED
kvg:radical CDATA #IMPLIED
kvg:phon CDATA #IMPLIED >
<!ATTLIST path
xmlns:kvg CDATA #FIXED "http://kanjivg.tagaini.net"
kvg:type CDATA #IMPLIED >
]>
`

var noRadHan = map[rune]bool{
	0x3005: true,
}

// Do we expect that k will have a valid radical?
func ExpectRadical(k rune) bool {
	if !unicode.In(k, unicode.Han) {
		return false
	}
	if noRadHan[k] {
		return false
	}
	return true
}

// Change the formatting of xmlout to that used by KanjiVG project,
// and add the common heading material. The KanjiVG files are indented
// using tabs, but with some of the indentation inconsistent.
func fixXML(xmlout []byte) []byte {
	outFixed := strings.Replace(string(xmlout), "></path>", "/>", -1)
	outFixed = strings.Replace(outFixed, "\t<g", "<g", -1)
	outFixed = strings.Replace(outFixed, "\t<g", "<g", -1)
	outFixed = strings.Replace(outFixed, "\t\t<path", "<path", -1)
	outFixed = strings.Replace(outFixed, "\t<text", "<text", -1)
	outFixed = strings.Replace(outFixed, "\t</g>", "</g>", -1)
	outFixed = strings.Replace(outFixed, "\t</g>", "</g>", -1)
	outstring := Heading + outFixed + "\n"
	return []byte(outstring)
}

// Make kanjivg into the XML of the KanjiVG files.
func (kanjivg *SVG) MakeXML() (output []byte) {
	return MakeXML(kanjivg)
}

// Make kanjivg into the XML of the KanjiVG files.
func MakeXML(kanjivg *SVG) (output []byte) {
	kanjivg.RenumberXML()
	output, err := xml.MarshalIndent(*kanjivg, "", "\t")
	if err != nil {
		log.Fatalf("Error marshalling: %s\n", err)
	}
	output = fixXML(output)
	return output
}

// Write kanjivg out as a file.
func (kanjivg *SVG) WriteKanjiFile(file string) {
	WriteKanjiFile(file, kanjivg)
}

// Write kanjivg to file.
func WriteKanjiFile(file string, kanjivg *SVG) {
	err := os.WriteFile(file, MakeXML(kanjivg), 0644)
	if err != nil {
		log.Fatalf("Error writing %s: %s\n", file, err)
	}
}

// Special marshaller for "child" elements, since a g may contain
// either a path or a text or another g element.
func (c Child) MarshalXML(e *xml.Encoder, start xml.StartElement) error {
	if c.IsGroup {
		start.Name = xml.Name{Local: "g"}
		return e.EncodeElement(c.Group, start)
	}
	if c.IsText {
		start.Name = xml.Name{Local: "text"}
		return e.EncodeElement(c.Text, start)
	}
	start.Name = xml.Name{Local: "path"}
	return e.EncodeElement(c.Path, start)
}

// Special unmarshaller. For some reason the kvg:type parts were not
// being picked up by the default parser, so I wrote this in order to
// work around that. There must be something I have missed about how
// the default unmarshal routine works.
func (p *Path) UnmarshalXML(d *xml.Decoder, start xml.StartElement) (err error) {
	for _, attr := range start.Attr {
		switch attr.Name.Local {
		case "d":
			p.D = attr.Value
		case "id":
			p.ID = attr.Value
		case "type":
			p.Type = attr.Value
		}
	}
	for {
		/* There is no content in the paths. */
		token, err := d.Token()
		if err != nil {
			return err
		}
		switch el := token.(type) {
		case xml.EndElement:
			if el == start.End() {
				return nil
			}
		}
	}
}

// Unmarshaller for a group.
func (g *Group) UnmarshalXML(d *xml.Decoder, start xml.StartElement) (err error) {
	for _, attr := range start.Attr {
		switch attr.Name.Local {
		case "element":
			g.Element = attr.Value
		case "id":
			g.ID = attr.Value
		case "number":
			g.Number = attr.Value
		case "original":
			g.Original = attr.Value
		case "part":
			g.Part = attr.Value
		case "partial":
			if attr.Value == "true" {
				g.Partial = true
			}
		case "phon":
			g.Phon = attr.Value
		case "position":
			g.Position = attr.Value
		case "radical":
			g.Radical = attr.Value
		case "radicalForm":
			g.RadicalForm = attr.Value
		case "style":
			g.Style = attr.Value
		case "tradForm":
			g.TradForm = attr.Value
		case "variant":
			if attr.Value == "true" {
				g.Variant = true
			}
		default:
			fmt.Printf("Unhandled -> %s\n", attr.Name.Local)
		}
	}
	for {
		token, err := d.Token()
		if err != nil {
			return err
		}
		var c Child
		c.Parent = g
		fail := false
		switch el := token.(type) {
		case xml.StartElement:
			switch el.Name.Local {
			case "g":
				err = d.DecodeElement(&c.Group, &el)
				if err != nil {
					return err
				}
				c.IsGroup = true
				c.IsText = false
				c.Group.Parent = &c
			case "path":
				err = d.DecodeElement(&c.Path, &el)
				if err != nil {
					return err
				}
				c.IsGroup = false
				c.IsText = false
				c.Path.Parent = &c
			case "text":
				err = d.DecodeElement(&c.Text, &el)
				if err != nil {
					return err
				}
				c.IsText = true
				c.IsGroup = false
				c.Text.Parent = &c
			default:
				fmt.Printf("Unhandled -> %s\n", el.Name.Local)
				fail = true
			}
			if !fail {
				g.Children = append(g.Children, c)
			}
		case xml.EndElement:
			if el == start.End() {
				return nil
			}
		}
	}
}

// Given a kanji in "contents", parse it into kanjivg. The error value
// comes from the XML unmarshalling.
func ParseKanji(contents []byte) (kanjivg SVG, oerr error) {
	oerr = xml.Unmarshal(contents, &kanjivg)
	if oerr != nil {
		return kanjivg, oerr
	}
	return kanjivg, nil
}

// The "parent" or "base" group of an SVG. This is a pointer to a
// value within kvg itself. See also Grab for an easy function which
// gets both the SVG and the base group from a file.
func (kvg *SVG) BaseGroup() (group *Group) {
	return &kvg.Groups[0].Children[0].Group
}

// Given a kanji file, read it and put the contents into kanjivg.
func ReadKanjiFile(file string) (kanjivg SVG, oerr error) {
	contents, oerr := os.ReadFile(file)
	if oerr != nil {
		return kanjivg, oerr
	}
	kanjivg, oerr = ParseKanji(contents)
	if oerr != nil {
		return kanjivg, oerr
	}
	return kanjivg, nil
}

// Renumber the groups and strokes recursively.
func renumber(child *Child, base string, nPathPtr, nGroupPtr *int64) {
	if child.IsGroup {
		*nGroupPtr++
		(*child).Group.ID = fmt.Sprintf("%s-g%d", base, *nGroupPtr)
		for i := range child.Group.Children {
			renumber(&child.Group.Children[i], base, nPathPtr, nGroupPtr)
		}
		return
	}
	*nPathPtr++
	(*child).Path.ID = fmt.Sprintf("%s-s%d", base, *nPathPtr)
}

// Get the value of the base element of kvg. The second return value
// is the basegroup ID without the "kvg:" prefix.
func (kvg *SVG) Base() (base, baseNoKVG string) {
	baseGroup := kvg.BaseGroup()
	base = baseGroup.ID
	baseNoKVG = strings.Replace(base, "kvg:", "", 1)
	return base, baseNoKVG
}

// Change the base of kvg to "base", e.g. from "kvg:01234" to
// "kvg:01234-MonkeyShines", and renumber all the groups and paths
// appropriately. By KanjiVG convention, the base always starts with
// the four characters "kvg:".  If the base of kvg is set using this,
// when the file is converted to XML, all of the id values in the
// output will use this base value, so there is no need to set it for
// each element.
func (kvg *SVG) SetBase(base string) {
	if base[0:4] != "kvg:" {
		log.Fatalf("Base name '%s' does not start with 'kvg:'", base)
	}
	baseGroup := kvg.BaseGroup()
	baseGroup.ID = base
	tail := base[4:]
	kvg.Groups[0].ID = "kvg:StrokePaths_" + tail
	/* The documentation says that the stroke numbers are optional, so
	   allow for the possibility that they may not exist. */
	if len(kvg.Groups) > 1 {
		kvg.Groups[1].ID = "kvg:StrokeNumbers_" + tail
	}
	var nPath int64
	var nGroup int64
	for i := range baseGroup.Children {
		renumber(&baseGroup.Children[i], base, &nPath, &nGroup)
	}
}

// Renumber the labels of the "text" group. The numerical labels given
// to the stroke numbers are simply their position within the "text"
// group, so the user does not need to keep track of the original
// numbers within the file.
func (kvg *SVG) RenumberLabels() {
	labels := kvg.Groups[1]
	for i := range labels.Children {
		c := &labels.Children[i]
		if !c.IsText {
			log.Printf("Error: non-text child in label %d\n", i+1)
			continue
		}
		c.Text.Content = []byte(fmt.Sprintf("%d", i+1))
	}
}

// Renumber an XML file read into "kvg" and also check its "style"
// elements.
func (svg *SVG) RenumberXML() {
	var nPath int64
	var nGroup int64
	baseGroup := svg.BaseGroup()
	base := baseGroup.ID
	for i := range baseGroup.Children {
		renumber(&baseGroup.Children[i], base, &nPath, &nGroup)
	}
	svg.RenumberLabels()
	svg.SetStyle()
}

// Set the "style" element to the normal KVG one.
func (svg *SVG) SetStyle() {
	svg.Groups[0].Style = "fill:none;stroke:#000000;stroke-width:3;stroke-linecap:round;stroke-linejoin:round;"
	svg.Groups[1].Style = "font-size:8;fill:#808080"
	removeChildStyles(&svg.Groups[0])
}

func removeChildStyles(g *Group) {
	c := (*g).Children
	for _, child := range c {
		if child.IsGroup {
			child.Group.Style = ""
			removeChildStyles(&child.Group)
		}
	}
}

// Read, renumber, and then write out a kanji file.
func RenumberFile(file string) {
	kvg := ReadKanjiFileOrDie(file)
	kvg.RenumberXML()
	WriteKanjiFile(file, &kvg)
}

// Helper for FindMultiElement
func multiSearch(depth int, pos []*Group, gp *Group, funky string, locs *[][]*Group, debug bool) {
	pos = append(pos, gp)
	if debug {
		fmt.Print("\n*")
		for i := 0; i < depth; i++ {
			fmt.Print("\t")
			fmt.Printf("%s", gp.ID)
		}
		defer fmt.Print("!\n")
	}
	for i, g := range gp.Children {
		if !g.IsGroup {
			continue
		}
		gPtr := &gp.Children[i].Group
		if g.Group.Element != funky {
			//			fmt.Printf("Not found in %s", g.Group.ID)
			multiSearch(depth+1, pos, gPtr, funky, locs, debug)
			continue
		}
		if debug {
			fmt.Printf("Found as group '%s'", g.Group.ID)
		}
		pos = append(pos, gPtr)
		copyPos := make([]*Group, len(pos))
		for i, p := range pos {
			copyPos[i] = p
		}
		*locs = append(*locs, copyPos)
	}
}

// WARNING: Please do not use this at the moment. It is due to be
// changed to a better interface shortly. 2024-09-19 Given a group gp,
// find all instances of subgroups with the kvg:element type of
// "funky", and return their locations as a chain of groups in
// "locs". To find a single instance, use FindElement instead of
// this. Annoyingly this puts things in the opposite order to
// FindElement and also does not return the true/false value of that.
func (gp *Group) FindMultiElement(funky string) (locs [][]*Group) {
	log.Printf("Don't use FindMultiElement")
	pos := make([]*Group, 0)
	locs = make([][]*Group, 0)
	multiSearch(0, pos, gp, funky, &locs, false)
	return locs
}

// See the documentation for FindElement(group,funky)
func (gp *Group) FindElement(funky string) (found bool, loc []*Group) {
	return FindElement(gp, funky)
}

// Find a group with the element "funky". This is usually enough since
// we usually only want to find a single subgroup. If found, return
// true and the location of it as a list of pointers to the parent
// groups it's in. The zeroth element of loc is the group which
// contains the element, the first is that element's parent, and so
// on. If the element is not found, return value is false and an empty
// slice.
func FindElement(gp *Group, funky string) (found bool, loc []*Group) {
	if gp.Element == funky {
		found = true
		loc = []*Group{gp}
		return true, loc
	}
	for i, g := range gp.Children {
		if g.IsGroup {
			gPtr := &gp.Children[i].Group
			found, loc = FindElement(gPtr, funky)
			if found {
				loc = append(loc, gp)
				return true, loc
			}
		}
	}
	return false, loc
}

func printLoc(loc []*Group) {
	for _, g := range loc {
		fmt.Printf("%s - %s\n", g.ID, g.Element)
	}
}

// Given a group g, get all the paths belonging to it, flattened into
// an array structure containing pointers to the actual paths.
func (g *Group) GetPaths() (paths []*Path) {
	return getPaths(g)
}

// Recursive helper for GetPaths
func getPaths(g *Group) (paths []*Path) {
	for i := range g.Children {
		c := &g.Children[i]
		if c.IsGroup {
			gpaths := getPaths(&c.Group)
			paths = append(paths, gpaths...)
			continue
		}
		paths = append(paths, &c.Path)
	}
	return paths
}

// Remove the KVDir prefix from a file name.
func TFile(file string) string {
	return strings.TrimPrefix(file, KVDir+"/")
}

// The structure which contains the radical information for a kanji. A
// single radical may consist of multiple groups, and there are four
// different types of radicals, hence the complicated structure. This
// uses pointers to the actual groups within your base structure.
type Radical struct {
	General, Tradit, Nelson, JIS []*Group
}

// Set this to true if you want SearchRadical to print when it finds a
// duplicate radical. Generally you do not want it to do this, since
// many of the radicals of the characters consist of more than one
// group of strokes.
var PrintDouble = false

// Find the radicals in g. This would usually be called on the base
// group of a character. The return values point to values within g
// itself, for the sake of modifying them.
func (g *Group) SearchRadical(radPtr *Radical) {
	for i := range g.Children {
		child := &g.Children[i]
		if !child.IsGroup {
			continue
		}
		child.Group.SearchRadical(radPtr)
	}
	rad := g.Radical
	if len(rad) == 0 {
		return
	}
	switch rad {
	case "general":
		if PrintDouble && len(radPtr.General) > 0 {
			fmt.Printf("Double %s for General.\n",
				radPtr.General[0].ID)
		}
		(*radPtr).General = append((*radPtr).General, g)
	case "nelson":
		if PrintDouble && len(radPtr.Nelson) > 0 {
			fmt.Printf("Double %s for Nelson.\n",
				radPtr.Nelson[0].ID)
		}
		(*radPtr).Nelson = append((*radPtr).Nelson, g)
	case "jis":
		if PrintDouble && len(radPtr.JIS) > 0 {
			fmt.Printf("Double %s for nelson.\n",
				radPtr.JIS[0].ID)
		}
		(*radPtr).JIS = append((*radPtr).JIS, g)
	case "tradit":
		if PrintDouble && len(radPtr.Tradit) > 0 {
			fmt.Printf("Double %s for tradit.\n",
				radPtr.Tradit[0].ID)
		}
		(*radPtr).Tradit = append((*radPtr).Tradit, g)
	default:
		log.Printf("Unknown value %s for kvg:radical.\n", rad)
	}
}

var hexID = "([0-9a-f]{5})"
var groupIDRe = regexp.MustCompile("kvg:" + hexID + ".*-g([0-9]+)")
var pathIDRe = regexp.MustCompile("kvg:" + hexID + ".*-s([0-9]+)")
var fileIDRe = regexp.MustCompile(".*" + hexID + "(?:-.+)?\\.svg$")
var filePartRe = regexp.MustCompile("(?:.*/)?(" + hexID + "(?:-([A-Za-z][A-Za-z0-9]*))?)\\.svg$")

// Hexadecimal to number. We already know hexID is valid from the
// regex validation, so we can fail fatally if this fails.
func HexIDToNum(hexID string) (num int64) {
	num, err := strconv.ParseInt(hexID, 16, 64)
	die(err, "Error parsing hex number")
	return num
}

// Given a KanjiVG file name fileName, return the hexadecimal id
// number, the kanji as a number, and the extension.
func FileToParts(fileName string) (id string, num int64, extension string) {
	match := filePartRe.FindStringSubmatch(fileName)
	if len(match) == 0 {
		return "", 0, ""
	}
	num = HexIDToNum(match[2])
	return match[1], num, match[3]
}

var Backup = regexp.MustCompile(`/\.#|/#|~$`)

// Given a group "group", get all its subgroups as a flat list. See
// Subgroups for another similar function.
func (base *Group) GetGroups() (groups []*Group) {
	groups = make([]*Group, 0)
	for i, c := range base.Children {
		if !c.IsGroup {
			continue
		}
		cgroups := base.Children[i].Group.GetGroups()
		groups = append(groups, cgroups...)
	}
	groups = append(groups, base)
	return groups
}

// Given a group "base", return its child groups as a map from the
// element to the group. See GetGroups for a simpler list return
// function.
func (base *Group) Subgroups() (elgr map[string][]*Group) {
	el := base.Element
	elgr = make(map[string][]*Group)
	for i, c := range base.Children {
		if !c.IsGroup {
			continue
		}
		celgr := base.Children[i].Group.Subgroups()
		for s, gr := range celgr {
			elgr[s] = append(elgr[s], gr...)
		}
	}
	elgr[el] = append(elgr[el], base)
	return elgr
}

// Function for operating on a file we have read in
type SVGFunc func(kanjivg SVG)

// Function for operating on a file which we have read in, with the
// file name
type SVGFileFunc func(file string, kanjivg SVG)

// Examine a file specified by path, and call fn on the contents if
// valid to do so.
func ExamineKanjiFile(path string, d fs.DirEntry, fn SVGFunc) (oerr error) {
	if d.IsDir() {
		return nil
	}
	if Backup.MatchString(path) {
		return nil
	}
	kanjivg, err := ReadKanjiFile(path)
	if err != nil {
		return err
	}
	fn(kanjivg)
	return nil
}

// Examine all the files in KVDir. Before using this, set KVDir to the
// value on your system. It will go through each file, read its
// contents, and call the function fn you specify on each of them. See
// also ExamineAllFilesSimple, which doesn't read the file contents
// first. This is usually better if you want to just check some files,
// since you can check the Unicode ID of the character using
// FileToNum, and decide whether to read it all in, rather than
// reading everything for all files.
func ExamineAllFiles(fn SVGFileFunc) {
	filepath.WalkDir(KVDir, func(path string, d fs.DirEntry, err error) error {
		if d.IsDir() {
			return nil
		}
		kanjivg := ReadKanjiFileOrDie(path)
		fn(path, kanjivg)
		return nil
	})
}

// Examine all the files in KVDir. Before using this, set KVDir to the
// value on your system. It will go through each file and call fn on
// them. It doesn't read the contents, unlike ExamineAllFiles.
func ExamineAllFilesSimple(fn func(file string)) {
	filepath.WalkDir(KVDir, func(path string, d fs.DirEntry, err error) error {
		if d.IsDir() {
			return nil
		}
		if Backup.MatchString(path) {
			return nil
		}
		fn(path)
		return nil
	})
}

// Get just the Unicode number from a kanjivg file name.
func FileToNum(fileName string) (num int64) {
	match := fileIDRe.FindStringSubmatch(fileName)
	if len(match) == 0 {
		fmt.Printf("No match in %s\n", fileName)
		os.Exit(1)
	}
	return HexIDToNum(match[1])
}

// Die if there is an error otherwise do nothing.
func die(err error, format string, a ...any) {
	if err == nil {
		return
	}
	fmt.Fprintf(os.Stderr, format, a...)
	fmt.Fprintf(os.Stderr, ": %s\n", err)
	os.Exit(1)
}

// The following is for non-web accessing of the kanji files.
//
// For web server applications like the viewer, we want to be able to
// read the kanji file and not die if the file name is faulty, but for
// non-web applications we want to stop processing completely if an
// error occurs.
func ReadKanjiFileOrDie(fileName string) (svg SVG) {
	svg, err := ReadKanjiFile(fileName)
	die(err, "Error unmarshalling '%s'", fileName)
	return svg
}

// Read a KanjiVG file, and return both the SVG and its base group. In
// practice it's a very common pattern to want to read a file, then
// usually one wants to read the base group, and when writing it back
// out again, one wants the full SVG, so this is a handy function in
// practice.
func Grab(fileName string) (svgPtr *SVG, base *Group) {
	svg := ReadKanjiFileOrDie(fileName)
	base = svg.BaseGroup()
	return &svg, base
}

// Given a group g, try to find an element with type (stroke shape)
// t. If an element is found, loc is an array of children with element
// 0 the element of type t and the remaining elements its successive
// parents. found is true or false depending on whether the element is
// found.
func FindType(g *Group, t string) (found bool, loc []*Child) {
	for i, c := range g.Children {
		if c.IsGroup {
			found, loc = FindType(&g.Children[i].Group, t)
			if found {
				loc = append(loc, &g.Children[i])
				return true, loc
			}
			continue
		}
		gc := &g.Children[i]
		if gc.Path.Type == t {
			return true, []*Child{gc}
		}
	}
	return false, loc
}

func (c *Child) Dump() (s string) {
	s += fmt.Sprintf("IsGroup: %t\n", c.IsGroup)
	s += fmt.Sprintf("IsText: %t\n", c.IsText)
	s += fmt.Sprintf("Group: %s\n", c.Group.Dump())
	s += fmt.Sprintf("Text: %s\n", c.Text.Content)
	return s
}

// Recursively-called function to make the string.
func (g *Group) dump(depth int) (s string) {
	indent := strings.Repeat("  ", depth)
	s += fmt.Sprintf("%s%s %s\n", indent, g.ID, g.Element)
	for _, c := range g.Children {
		if c.IsGroup {
			s += c.Group.dump(depth + 1)
			continue
		}
		s += fmt.Sprintf("%s  %s %s\n", indent, c.Path.ID, c.Path.Type)
	}
	return s
}

// Convert g into a printable string showing the IDs, group elements,
// and path types, indented by depth.
func (g *Group) Dump() (s string) {
	return g.dump(0)
}

// Return all the paths in the base group of svg.
func (svg *SVG) GetPaths() (paths []*Path) {
	base := svg.BaseGroup()
	return getPaths(base)
}

// Get the numeric part of a path ID.
func PathIDToNum(id string) (num int64) {
	match := pathIDRe.FindStringSubmatch(id)
	return decimalToNum(match[2])
}

// Decimal to number. We already know it is valid from the regex
// validation so we can fail fatally if this fails.
func decimalToNum(Decimal string) (num int64) {
	num, err := strconv.ParseInt(Decimal, 10, 64)
	die(err, "Error parsing decimal")
	return num
}

// Given a group g, return its element
func (g *Group) El() (element string) {
	element = g.Element
	if len(g.Original) > 0 {
		element = g.Original
	}
	return element
}

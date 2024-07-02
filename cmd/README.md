# SUBDIRECTORIES

* __bogusgroup__ is a tool to find groups with no paths in them

* __empty-path__ finds files where the number of strokes does not
match the number of stroke number labels. It also locates instances
of empty paths with no information. As of 2024-06-20 there are no
instances in the repository.

* __kvg-mode.el__ provides an Emacs editing mode which automatically
renumbers all the XML elements for consistency, and indents the
buffer each time the file is saved (C-x C-s). It requires the user
already has go-mode.el installed. It also uses a hard-coded path for
renumber, so it will require end-user editing to be used correctly.

* __skip__ compares SKIP ("System of Kanji Indexing by Patterns")
  against values calculated from the KanjiVG breakdowns.

* __read-write-test__ provides a utility which reads and then
writes back out all the files of kvg, and prints a report on which
files differ from the standard formatting.

* __renumber__ provides a utility which reformats and renumbers the
files provided on the command line. This is used by the Emacs editing
mode.

* __typeshift__ is a tool for shuffling the `kvg:type` values of strokes.


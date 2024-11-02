# FILES IN THIS DIRECTORY

* __skip.go__ is an attempt at computing the SKIP kanji code from the
KanjiVG information. This uses a file skip.json which is taken from
Kanjidic.

* __skip.json__ is the data for `skip.go`.

* __make-skip-json.pl__ is a Perl script which makes the file
`skip.json` from a copy of Kanjidic. You probably don't need to run
this, since `skip.json` is already in the repository, but if you do
need it, please install the prerequisite modules using `cpanm
Data::Kanji::Kanjidic JSON::Create`.

## SKIP system required notice

The SKIP codes are used under the conditions described below:

http://nihongo.monash.edu/edrdg/skipperm.html

Due to the nature of the JSON format, this information could not be
included in the file __skip.json__ itself.

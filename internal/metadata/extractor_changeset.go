package metadata

import (
	"regexp"
	"strings"
)

type ChangeSetExtractor struct {
	re      *regexp.Regexp
	cleaner *regexp.Regexp
	saver   changeSetSaver
}

func (e ChangeSetExtractor) Matches(line string) bool {
	return e.re.MatchString(line)
}
func (e ChangeSetExtractor) Extract(metadata *BackupLogMetadata, line string) {
	// Sample line: "Files:           0 new,     0 changed,     2 unmodified"
	trimmed := strings.TrimSpace(e.re.ReplaceAllString(line, ""))
	splitted := strings.Split(trimmed, ",")
	var changeset BackupLogMetadataChangeset = BackupLogMetadataChangeset{}
	changeset.Added = e.cleaner.ReplaceAllString(splitted[0], "")
	changeset.Changed = e.cleaner.ReplaceAllString(splitted[1], "")
	changeset.Unmodified = e.cleaner.ReplaceAllString(splitted[2], "")
	e.saver.Save(metadata, changeset)
}

type changeSetSaver interface {
	Save(metadata *BackupLogMetadata, changeset BackupLogMetadataChangeset)
}

type fileSaver struct{}

func (f fileSaver) Save(metadata *BackupLogMetadata, changeset BackupLogMetadataChangeset) {
	metadata.Files = changeset
}

type dirsSaver struct{}

func (d dirsSaver) Save(metadata *BackupLogMetadata, changeset BackupLogMetadataChangeset) {
	metadata.Dirs = changeset
}

func NewFilesExtractor() MetadatExtractor {
	return ChangeSetExtractor{
		re:      regexp.MustCompile(`(?i)^Files:`),
		cleaner: regexp.MustCompile(`[^\d]`),
		saver:   fileSaver{},
	}
}
func NewDirsExtractor() MetadatExtractor {
	return ChangeSetExtractor{
		re:      regexp.MustCompile(`(?i)^Dirs:`),
		cleaner: regexp.MustCompile(`[^\d]`),
		saver:   dirsSaver{},
	}
}

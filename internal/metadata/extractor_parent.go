package metadata

import (
	"regexp"
	"strings"
)

type parentSnapshotIDExtractor struct {
	re *regexp.Regexp
}

func (e parentSnapshotIDExtractor) Matches(line string) bool {
	return e.re.MatchString(line)
}
func (e parentSnapshotIDExtractor) Extract(metadata *BackupLogMetadata, line string) {
	// Sample line: "using parent snapshot c65d9310"
	metadata.ParentSnapshotID = strings.TrimSpace(e.re.ReplaceAllString(line, ""))
}

func NewParentSnapshotIDExtractor() MetadatExtractor {
	return parentSnapshotIDExtractor{regexp.MustCompile(`(?i)^using parent snapshot`)}
}

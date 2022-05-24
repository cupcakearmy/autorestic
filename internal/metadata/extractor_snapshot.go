package metadata

import (
	"regexp"
	"strings"
)

type snapshotExtractor struct {
	re *regexp.Regexp
}

func (e snapshotExtractor) Matches(line string) bool {
	return e.re.MatchString(line)
}
func (e snapshotExtractor) Extract(metadata *BackupLogMetadata, line string) {
	// Sample line: "snapshot 917c7691 saved"
	metadata.SnapshotID = strings.Split(line, " ")[1]
}

func NewSnapshotExtractor() MetadatExtractor {
	return snapshotExtractor{regexp.MustCompile(`(?i)^snapshot \w+ saved`)}
}

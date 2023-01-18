package metadata

import (
	"regexp"
	"strings"
)

type addedExtractor struct {
	re *regexp.Regexp
}

func (e addedExtractor) Matches(line string) bool {
	return e.re.MatchString(line)
}
func (e addedExtractor) Extract(metadata *BackupLogMetadata, line string) {
	// Sample line: "Added to the repo: 0 B"
	metadata.AddedSize = strings.TrimSpace(e.re.ReplaceAllString(line, "$2"))
}

func NewAddedExtractor() MetadatExtractor {
	return addedExtractor{regexp.MustCompile(`(?i)^Added to the repo(sitory)?: ([\d\.]+ \w+)( \([\d\.]+[\w\s]+\))?`)}
}

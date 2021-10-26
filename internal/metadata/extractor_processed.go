package metadata

import (
	"regexp"
	"strings"
)

type processedExtractor struct {
	re      *regexp.Regexp
	cleaner *regexp.Regexp
}

func (e processedExtractor) Matches(line string) bool {
	return e.re.MatchString(line)
}
func (e processedExtractor) Extract(metadata *BackupLogMetadata, line string) {
	// Sample line: "processed 2 files, 24 B in 0:00"
	var processed = BackupLogMetadataProcessed{}
	split := strings.Split(line, "in")
	processed.Duration = strings.TrimSpace(split[1])
	split = strings.Split(split[0], ",")
	processed.Files = e.cleaner.ReplaceAllString(split[0], "")
	processed.Size = strings.TrimSpace(split[1])
	metadata.Processed = processed
}

func NewProcessedExtractor() MetadatExtractor {
	return processedExtractor{
		regexp.MustCompile(`(?i)^processed \d* files`),
		regexp.MustCompile(`(?i)[^\d]`),
	}
}

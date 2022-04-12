package metadata

import (
	"strings"
)

type BackupLogMetadataChangeset struct {
	Added      string
	Changed    string
	Unmodified string
}
type BackupLogMetadataProcessed struct {
	Files    string
	Size     string
	Duration string
}
type BackupLogMetadata struct {
	ParentSnapshotID string
	Files            BackupLogMetadataChangeset
	Dirs             BackupLogMetadataChangeset
	AddedSize        string
	Processed        BackupLogMetadataProcessed
	SnapshotID       string
	ExitCode         string
}

type MetadatExtractor interface {
	Matches(line string) bool
	Extract(metadata *BackupLogMetadata, line string)
}

var extractors = []MetadatExtractor{
	NewParentSnapshotIDExtractor(),
	NewFilesExtractor(),
	NewDirsExtractor(),
	NewAddedExtractor(),
	NewProcessedExtractor(),
	NewSnapshotExtractor(),
}

func ExtractMetadataFromBackupLog(log string) BackupLogMetadata {
	var md BackupLogMetadata
	for _, line := range strings.Split(log, "\n") {
		line = strings.TrimSpace(line)
		for _, extractor := range extractors {
			if extractor.Matches(line) {
				extractor.Extract(&md, line)
				continue
			}
		}
	}
	return md
}

func MakeEnvFromMetadata(metadata *BackupLogMetadata) map[string]string {
	env := make(map[string]string)
	var prefix = "AUTORESTIC_"

	env[prefix+"SNAPSHOT_ID"] = metadata.SnapshotID
	env[prefix+"PARENT_SNAPSHOT_ID"] = metadata.ParentSnapshotID
	env[prefix+"FILES_ADDED"] = metadata.Files.Added
	env[prefix+"FILES_CHANGED"] = metadata.Files.Changed
	env[prefix+"FILES_UNMODIFIED"] = metadata.Files.Unmodified
	env[prefix+"DIRS_ADDED"] = metadata.Dirs.Added
	env[prefix+"DIRS_CHANGED"] = metadata.Dirs.Changed
	env[prefix+"DIRS_UNMODIFIED"] = metadata.Dirs.Unmodified
	env[prefix+"ADDED_SIZE"] = metadata.AddedSize
	env[prefix+"PROCESSED_FILES"] = metadata.Processed.Files
	env[prefix+"PROCESSED_SIZE"] = metadata.Processed.Size
	env[prefix+"PROCESSED_DURATION"] = metadata.Processed.Duration
	env[prefix+"EXIT_CODE"] = metadata.ExitCode

	return env
}

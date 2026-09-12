package cmd

import "testing"

func TestPsLineLooksLikeAutorestic(t *testing.T) {
	const self = "12345"

	cases := []struct {
		name string
		line string
		want bool
	}{
		{
			name: "systemctl unit name contains autorestic",
			line: "root      346977  0.0  0.0  17808  6784 pts/1    S+   09:23   0:00 systemctl start autorestic-backup.service",
			want: false,
		},
		{
			name: "sudo systemctl with autorestic in unit",
			line: "root      346974  9.0  0.0  17992  7296 pts/0    S+   09:23   0:00 sudo systemctl start autorestic-backup.service",
			want: false,
		},
		{
			name: "systemd-inhibit wrapping autorestic",
			line: "root      346995  0.0  0.0  12345  6784 ?        Ss   09:23   0:00 systemd-inhibit autorestic --ci --config config.yaml backup --all",
			want: true,
		},
		{
			name: "direct autorestic binary path",
			line: "root      99999  0.0  0.0  12345  6784 ?        Ss   09:23   0:00 /usr/local/bin/autorestic backup --all",
			want: true,
		},
		{
			name: "current process pid excluded",
			line: "root      12345  0.0  0.0  12345  6784 ?        Ss   09:23   0:00 /usr/local/bin/autorestic unlock",
			want: false,
		},
		{
			name: "unrelated process",
			line: "root      11111  0.0  0.0  12345  6784 ?        Ss   09:23   0:00 /usr/bin/restic backup /",
			want: false,
		},
	}

	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			got := psLineLooksLikeAutorestic(tc.line, self)
			if got != tc.want {
				t.Errorf("got %v, want %v for %q", got, tc.want, tc.line)
			}
		})
	}
}

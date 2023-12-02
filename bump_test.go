package main

import (
	"os/exec"
	"testing"
)

func TestUpdate(t *testing.T) {
	testCases := []struct {
		name        string
		pkgVersions []pkgVersion
		want        map[string]string
	}{
		{
			name: "standard update",
			pkgVersions: []pkgVersion{
				{
					Name:    "github.com/google/uuid",
					Version: "v1.4.0",
				},
			},
			want: map[string]string{
				"github.com/google/uuid": "v1.4.0",
			},
		},
		{
			name: "replace",
			pkgVersions: []pkgVersion{
				{
					Name:    "k8s.io/client-go",
					Version: "v0.28.0",
				},
			},
			want: map[string]string{
				"k8s.io/client-go": "v0.28.0",
			},
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			tmpdir := t.TempDir()
			copyFile(t, "testdata/aws-efs-csi-driver/go.mod", tmpdir)

			modFile, err := doUpdate(tc.pkgVersions, tmpdir)
			if err != nil {
				t.Fatal(err)
			}
			for pkg, want := range tc.want {
				if got := getVersion(modFile, pkg); got != want {
					t.Errorf("expected %s, got %s", want, got)
				}
			}
		})
	}
}

func TestUpdateError(t *testing.T) {
	testCases := []struct {
		name        string
		pkgVersions []pkgVersion
	}{
		{
			name: "no downgrade",
			pkgVersions: []pkgVersion{
				{
					Name:    "github.com/google/uuid",
					Version: "v1.0.0",
				},
			},
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			tmpdir := t.TempDir()
			copyFile(t, "testdata/aws-efs-csi-driver/go.mod", tmpdir)

			_, err := doUpdate(tc.pkgVersions, tmpdir)
			if err == nil {
				t.Fatal("expected error, got nil")
			}
		})
	}
}

func TestCommit(t *testing.T) {
	// We use github.com/NVIDIA/go-nvml v0.11.7-0 in our go.mod
	// That corresponds to 53c34bc04d66e9209eff8654bc70563cf380e214
	pkg := "github.com/NVIDIA/go-nvml"

	// An older commit is c3a16a2b07cf2251cbedb76fa68c9292b22bfa06
	olderCommit := "c3a16a2b07cf2251cbedb76fa68c9292b22bfa06"
	olderVersion := "v0.11.6-0"
	// A newer commit is 95ef6acc3271a9894fd02c1071edef1d88527e20
	newerCommit := "95ef6acc3271a9894fd02c1071edef1d"
	newerVersion := "v0.12.0-1"

	testCases := []struct {
		name    string
		version string
		want    map[string]string
	}{
		{
			name:    "pin to older",
			version: olderCommit,
			want: map[string]string{
				pkg: olderVersion,
			},
		},
		{
			name:    "pin to newer",
			version: newerCommit,
			want: map[string]string{
				pkg: newerVersion,
			},
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			tmpdir := t.TempDir()
			copyFile(t, "testdata/aws-efs-csi-driver/go.mod", tmpdir)

			pkgVersions := []pkgVersion{
				{
					Name:    pkg,
					Version: tc.version,
				},
			}
			modFile, err := doUpdate(pkgVersions, tmpdir)
			if err != nil {
				t.Fatal(err)
			}
			for pkg, want := range tc.want {
				if got := getVersion(modFile, pkg); got != want {
					t.Errorf("expected %s, got %s", want, got)
				}
			}
		})
	}

}

func copyFile(t *testing.T, src, dst string) {
	t.Helper()
	_, err := exec.Command("cp", "-r", src, dst).Output()
	if err != nil {
		t.Fatal(err)
	}
}

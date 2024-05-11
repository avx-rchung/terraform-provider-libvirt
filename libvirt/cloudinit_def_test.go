package libvirt

import (
	"os"
	"testing"

	"github.com/diskfs/go-diskfs"
)

func TestCloudInitTerraformKeyOps(t *testing.T) {
	ci := newCloudInitDef()

	volKey := "volume-key"

	terraformID := ci.buildTerraformKey(volKey)
	if terraformID == "" {
		t.Error("key should not be empty")
	}

	actualKey, _ := getCloudInitVolumeKeyFromTerraformID(terraformID)
	if actualKey != volKey {
		t.Error("wrong key returned")
	}
}

func TestCloudInitCreateISO(t *testing.T) {
	ci := newCloudInitDef()
	ci.Name = "test.iso"
	ci.UserData = "test-user-data"
	ci.MetaData = "test-meta-data"
	ci.NetworkConfig = "test-network-config"

	iso, err := ci.createISO()
	if err != nil {
		t.Errorf("Unexpected error: %v", err)
	}
	if iso == "" {
		t.Errorf("Unexpected iso to be empty")
	}
	t.Logf("iso: %s", iso)

	disk, err := diskfs.Open(iso)
	if err != nil {
		t.Fatalf("Failed to open iso: %v", err)
	}
	fs, err := disk.GetFilesystem(0)
	if err != nil {
		t.Fatalf("Failed to get filesystem: %v", err)
	}

	for _, path := range []string{"/user-data", "/meta-data", "/network-config"} {
		f, err := fs.OpenFile(path, os.O_RDONLY)
		if err != nil {
			t.Errorf("Failed to open file %s: %v", path, err)
		}
		f.Close()
	}
}

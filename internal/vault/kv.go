package vault

import "strings"

// KVVersion represents the KV secrets engine version.
type KVVersion int

const (
	KVv1 KVVersion = 1
	KVv2 KVVersion = 2
)

// NormalizePath adjusts a secret path based on the KV engine version.
// For KV v2, it injects "data/" after the mount point if not already present.
func NormalizePath(mount, path string, version KVVersion) string {
	if version == KVv1 {
		return strings.TrimPrefix(path, "/")
	}

	// Ensure mount has no trailing slash
	mount = strings.TrimRight(mount, "/")
	relative := strings.TrimPrefix(path, mount)
	relative = strings.TrimPrefix(relative, "/")

	// Already has data/ prefix
	if strings.HasPrefix(relative, "data/") {
		return mount + "/" + relative
	}

	return mount + "/data/" + relative
}

// SplitMountPath attempts to split a full path into (mount, subpath).
// It uses the first path segment as the mount point.
func SplitMountPath(fullPath string) (mount, subpath string) {
	parts := strings.SplitN(strings.TrimPrefix(fullPath, "/"), "/", 2)
	if len(parts) == 1 {
		return parts[0], ""
	}
	return parts[0], parts[1]
}

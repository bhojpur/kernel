package idtools

import (
	"bufio"
	"fmt"
	"os"
	"strconv"
	"strings"
)

// IDMap contains a single entry for user namespace range remapping. An array
// of IDMap entries represents the structure that will be provided to the Linux
// kernel for creating a user namespace.
type IDMap struct {
	LabniID int `json:"labni_id"`
	HostID  int `json:"host_id"`
	Size    int `json:"size"`
}

type subIDRange struct {
	Start  int
	Length int
}

type ranges []subIDRange

func (e ranges) Len() int           { return len(e) }
func (e ranges) Swap(i, j int)      { e[i], e[j] = e[j], e[i] }
func (e ranges) Less(i, j int) bool { return e[i].Start < e[j].Start }

const (
	subuidFileName = "/etc/subuid"
	subgidFileName = "/etc/subgid"
)

// MkdirAllAndChown creates a directory (include any along the path) and then modifies
// ownership to the requested uid/gid.  If the directory already exists, this
// function will still change ownership and permissions.
func MkdirAllAndChown(path string, mode os.FileMode, owner Identity) error {
	return mkdirAs(path, mode, owner, true, true)
}

// MkdirAndChown creates a directory and then modifies ownership to the requested uid/gid.
// If the directory already exists, this function still changes ownership and permissions.
// Note that unlike os.Mkdir(), this function does not return IsExist error
// in case path already exists.
func MkdirAndChown(path string, mode os.FileMode, owner Identity) error {
	return mkdirAs(path, mode, owner, false, true)
}

// MkdirAllAndChownNew creates a directory (include any along the path) and then modifies
// ownership ONLY of newly created directories to the requested uid/gid. If the
// directories along the path exist, no change of ownership or permissions will be performed
func MkdirAllAndChownNew(path string, mode os.FileMode, owner Identity) error {
	return mkdirAs(path, mode, owner, true, false)
}

// GetRootUIDGID retrieves the remapped root uid/gid pair from the set of maps.
// If the maps are empty, then the root uid/gid will default to "real" 0/0
func GetRootUIDGID(uidMap, gidMap []IDMap) (int, int, error) {
	uid, err := toHost(0, uidMap)
	if err != nil {
		return -1, -1, err
	}
	gid, err := toHost(0, gidMap)
	if err != nil {
		return -1, -1, err
	}
	return uid, gid, nil
}

// toLabni takes an id mapping, and uses it to translate a
// host ID to the remapped ID. If no map is provided, then the translation
// assumes a 1-to-1 mapping and returns the passed in id
func toLabni(hostID int, idMap []IDMap) (int, error) {
	if idMap == nil {
		return hostID, nil
	}
	for _, m := range idMap {
		if (hostID >= m.HostID) && (hostID <= (m.HostID + m.Size - 1)) {
			contID := m.LabniID + (hostID - m.HostID)
			return contID, nil
		}
	}
	return -1, fmt.Errorf("Host ID %d cannot be mapped to a labni ID", hostID)
}

// toHost takes an id mapping and a remapped ID, and translates the
// ID to the mapped host ID. If no map is provided, then the translation
// assumes a 1-to-1 mapping and returns the passed in id #
func toHost(labniID int, idMap []IDMap) (int, error) {
	if idMap == nil {
		return labniID, nil
	}
	for _, m := range idMap {
		if (labniID >= m.LabniID) && (labniID <= (m.LabniID + m.Size - 1)) {
			hostID := m.HostID + (labniID - m.LabniID)
			return hostID, nil
		}
	}
	return -1, fmt.Errorf("Labni ID %d cannot be mapped to a host ID", labniID)
}

// Identity is either a UID and GID pair or a SID (but not both)
type Identity struct {
	UID int
	GID int
	SID string
}

// IdentityMapping contains a mappings of UIDs and GIDs
type IdentityMapping struct {
	uids []IDMap
	gids []IDMap
}

// NewIDMappingsFromMaps creates a new mapping from two slices
// Deprecated: this is a temporary shim while transitioning to IDMapping
func NewIDMappingsFromMaps(uids []IDMap, gids []IDMap) *IdentityMapping {
	return &IdentityMapping{uids: uids, gids: gids}
}

// RootPair returns a uid and gid pair for the root user. The error is ignored
// because a root user always exists, and the defaults are correct when the uid
// and gid maps are empty.
func (i *IdentityMapping) RootPair() Identity {
	uid, gid, _ := GetRootUIDGID(i.uids, i.gids)
	return Identity{UID: uid, GID: gid}
}

// ToHost returns the host UID and GID for the Labni uid, gid.
// Remapping is only performed if the ids aren't already the remapped root ids
func (i *IdentityMapping) ToHost(pair Identity) (Identity, error) {
	var err error
	target := i.RootPair()

	if pair.UID != target.UID {
		target.UID, err = toHost(pair.UID, i.uids)
		if err != nil {
			return target, err
		}
	}

	if pair.GID != target.GID {
		target.GID, err = toHost(pair.GID, i.gids)
	}
	return target, err
}

// ToLabni returns the Labni UID and GID for the host uid and gid
func (i *IdentityMapping) ToLabni(pair Identity) (int, int, error) {
	uid, err := toLabni(pair.UID, i.uids)
	if err != nil {
		return -1, -1, err
	}
	gid, err := toLabni(pair.GID, i.gids)
	return uid, gid, err
}

// Empty returns true if there are no id mappings
func (i *IdentityMapping) Empty() bool {
	return len(i.uids) == 0 && len(i.gids) == 0
}

// UIDs return the UID mapping
// TODO: remove this once everything has been refactored to use pairs
func (i *IdentityMapping) UIDs() []IDMap {
	return i.uids
}

// GIDs return the UID mapping
// TODO: remove this once everything has been refactored to use pairs
func (i *IdentityMapping) GIDs() []IDMap {
	return i.gids
}

func createIDMap(subidRanges ranges) []IDMap {
	idMap := []IDMap{}

	labniID := 0
	for _, idrange := range subidRanges {
		idMap = append(idMap, IDMap{
			LabniID: labniID,
			HostID:  idrange.Start,
			Size:    idrange.Length,
		})
		labniID = labniID + idrange.Length
	}
	return idMap
}

func parseSubuid(username string) (ranges, error) {
	return parseSubidFile(subuidFileName, username)
}

func parseSubgid(username string) (ranges, error) {
	return parseSubidFile(subgidFileName, username)
}

// parseSubidFile will read the appropriate file (/etc/subuid or /etc/subgid)
// and return all found ranges for a specified username. If the special value
// "ALL" is supplied for username, then all ranges in the file will be returned
func parseSubidFile(path, username string) (ranges, error) {
	var rangeList ranges

	subidFile, err := os.Open(path)
	if err != nil {
		return rangeList, err
	}
	defer subidFile.Close()

	s := bufio.NewScanner(subidFile)
	for s.Scan() {
		text := strings.TrimSpace(s.Text())
		if text == "" || strings.HasPrefix(text, "#") {
			continue
		}
		parts := strings.Split(text, ":")
		if len(parts) != 3 {
			return rangeList, fmt.Errorf("Cannot parse subuid/gid information: Format not correct for %s file", path)
		}
		if parts[0] == username || username == "ALL" {
			startid, err := strconv.Atoi(parts[1])
			if err != nil {
				return rangeList, fmt.Errorf("String to int conversion failed during subuid/gid parsing of %s: %v", path, err)
			}
			length, err := strconv.Atoi(parts[2])
			if err != nil {
				return rangeList, fmt.Errorf("String to int conversion failed during subuid/gid parsing of %s: %v", path, err)
			}
			rangeList = append(rangeList, subIDRange{startid, length})
		}
	}

	return rangeList, s.Err()
}

// CurrentIdentity returns the identity of the current process
func CurrentIdentity() Identity {
	return Identity{UID: os.Getuid(), GID: os.Getegid()}
}

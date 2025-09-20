// /service/hash/hasher.go
package hash

import (
    "crypto/sha256"
    "encoding/hex"
    "io"
)

func SHA256FromReader(r io.Reader) (string, int64, error) {
    hasher := sha256.New()
    n, err := io.Copy(hasher, r)
    if err != nil {
        return "", 0, err
    }
    return hex.EncodeToString(hasher.Sum(nil)), n, nil
}

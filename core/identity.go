package core

import (
	"bytes"
	"crypto/rand"
	"crypto/sha1"
	"encoding/base64"
	"fmt"
	"net"
)

type Identity struct {
	// Must be public for marshalling
	Name string
	Hash string
}

// TODO Better
// TODO find a way to have unique identifier for each client and channel

func NewIdentity(name string) *Identity {
	return &Identity{
		Name: name,
		Hash: genIdHash(),
	}
}

func genIdHash() (id string) {
	interfaces, err := net.Interfaces()
	if err != nil {
		id = randToken(128)
	} else {
		var buffer bytes.Buffer

		for _, inter := range interfaces {
			h := inter.HardwareAddr.String()
			n := inter.Name
			if len(h) == 0 || len(n) == 0 {
				continue
			}
			buffer.WriteString(fmt.Sprintf("%s%s_", n, h))
		}
		buffer.WriteString(randToken(32))
		hasher := sha1.New()
		hasher.Write(buffer.Bytes())
		id = base64.URLEncoding.EncodeToString(hasher.Sum(nil))
	}
	return
}

func randToken(size int) string {
	b := make([]byte, size)
	rand.Read(b)
	return fmt.Sprintf("%x", b)
}

func (id *Identity) String() string {
	return fmt.Sprintf("(%s,%s)", id.Name, id.Hash)
}

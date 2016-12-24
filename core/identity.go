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
	Name string
	Hash string
}

// TODO Better

func newIdentity(name string) *Identity {
	return &Identity{
		Name: name,
		Hash: genIdHash(),
	}
}

func genIdHash() (id string) {
	interfaces, err := net.Interfaces()
	if err != nil {
		id = randToken(64)
	} else {
		var buffer bytes.Buffer

		for _, inter := range interfaces {
			h := inter.HardwareAddr.String()
			n := inter.Name
			if len(h) == 0 || len(n) == 0 {
				continue
			}
			buffer.WriteString(n)
			buffer.WriteString(h)
			buffer.WriteString("_")
		}
		buffer.WriteString(randToken(16))
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

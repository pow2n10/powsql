package mysql

import (
	"crypto/sha1"
	"encoding/binary"
	"errors"
)

func CalcPassword(scramble, password []byte) []byte {
	if len(password) == 0 {
		return nil
	}

	// stage1Hash = SHA1(password)
	crypt := sha1.New()
	crypt.Write(password)
	stage1 := crypt.Sum(nil)

	// scrambleHash = SHA1(scramble + SHA1(stage1Hash))
	// inner Hash
	crypt.Reset()
	crypt.Write(stage1)
	hash := crypt.Sum(nil)

	// outer Hash
	crypt.Reset()
	crypt.Write(scramble)
	crypt.Write(hash)
	scramble = crypt.Sum(nil)

	// token = scrambleHash XOR stage1Hash
	for i := range scramble {
		scramble[i] ^= stage1[i]
	}
	return scramble
}

func ReadUint64(data []byte) (uint64, int, error) {

	if len(data) < 1 {
		return 0, 0, errors.New("empty slice")
	}

	if data[0] < 0xfb {
		return uint64(data[0]), 1, nil
	}

	if data[1] == 0xfc {
		if len(data) < 2 {
			return 0, 0, errors.New("empty slice")
		}
		return binary.LittleEndian.Uint64(data[0:2]), 2, nil
	}

	if data[1] == 0xfd {
		if len(data) < 3 {
			return 0, 0, errors.New("empty slice")
		}
		return binary.LittleEndian.Uint64(data[0:3]), 3, nil
	}
	if data[1] == 0xfe {
		if len(data) < 8 {
			return 0, 0, errors.New("empty slice")
		}
		return binary.LittleEndian.Uint64(data[0:8]), 8, nil
	}

	return 0, 0, errors.New("unknow data type")

}

func ReadInt64(data []byte) (int64, int, error) {
	n, m, e := ReadUint64(data)
	return int64(n), m, e
}

package xid

import (
	"crypto/md5"
	"crypto/rand"
	"encoding/binary"
	"encoding/hex"
	"fmt"
	"io"
	"os"
	"sync/atomic"
	"time"
	"bytes"
	"errors"
)

// 从 mgo.bson 复制

var machineId = readMachineId()
var processId = os.Getpid()
var objectIdCounter uint32 = readRandomUint32()

type XID string

func (this XID) Hex() string {
	return hex.EncodeToString([]byte(this))
}

func (this XID) String() string {
	return fmt.Sprintf(`XID("%x")`, string(this))
}

func (this XID) MarshalJSON() ([]byte, error) {
	return []byte(fmt.Sprintf(`"%x"`, string(this))), nil
}

var nullBytes = []byte("null")

func (this *XID) UnmarshalJSON(data []byte) error {
	if len(data) == 2 && data[0] == '"' && data[1] == '"' || bytes.Equal(data, nullBytes) {
		*this = ""
		return nil
	}
	if (len(data) != 26 && len(data) != 14) || data[0] != '"' || (data[13] != '"' && data[25] != '"') {
		return errors.New(fmt.Sprintf("invalid XID in JSON: %s", string(data)))
	}
	if len(data) == 26 {
		var buf [12]byte
		_, err := hex.Decode(buf[:], data[1:25])
		if err != nil {
			return errors.New(fmt.Sprintf("invalid XID in JSON: %s (%s)", string(data), err))
		}
		*this = XID(string(buf[:]))
	} else if len(data) == 14 {
		var buf [6]byte
		_, err := hex.Decode(buf[:], data[1:13])
		if err != nil {
			return errors.New(fmt.Sprintf("invalid XID in JSON: %s (%s)", string(data), err))
		}
		*this = XID(string(buf[:]))
	}

	return nil
}

func (this XID) Valid() bool {
	return len(this) == 12
}

func (this XID) Time() time.Time {
	var sec = int64(binary.BigEndian.Uint32(this.byteSlice(0, 4)))
	return time.Unix(sec, 0)
}

func (this XID) Machine() []byte {
	return this.byteSlice(4, 7)
}

func (this XID) Counter() int32 {
	b := this.byteSlice(9, 12)
	return int32(uint32(b[0])<<16 | uint32(b[1])<<8 | uint32(b[2]))
}

func (this XID) Pid() uint16 {
	return binary.BigEndian.Uint16(this.byteSlice(7, 9))
}

func (this XID) byteSlice(start, end int) []byte {
	if len(this) != 12 {
		panic(fmt.Sprintf("invalid ObjectId: %q", string(this)))
	}
	return []byte(string(this)[start:end])
}

func XIDHex(v string) XID {
	var d, err = hex.DecodeString(v)
	if err != nil || (len(d) != 12 && len(d) != 6) {
		panic(fmt.Sprintf("invalid input to XID Hex: %q", v))
	}
	return XID(d)
}

func IsXIDHex(s string) bool {
	if len(s) != 24 && len(s) != 12 {
		return false
	}
	_, err := hex.DecodeString(s)
	return err == nil
}

func NewXID() XID {
	return NewXIDWithTime(time.Now())
}

func NewXIDWithTime(t time.Time) XID {
	var b [12]byte
	binary.BigEndian.PutUint32(b[:4], uint32(t.Unix()))
	b[4] = machineId[0]
	b[5] = machineId[1]
	b[6] = machineId[2]

	b[7] = byte(processId >> 8)
	b[8] = byte(processId)

	i := atomic.AddUint32(&objectIdCounter, 1)
	b[9] = byte(i >> 16)
	b[10] = byte(i >> 8)
	b[11] = byte(i)
	return XID(string(b[:]))
}

func readMachineId() []byte {
	var sum [3]byte
	id := sum[:]
	hostname, err1 := os.Hostname()
	if err1 != nil {
		_, err2 := io.ReadFull(rand.Reader, id)
		if err2 != nil {
			panic(fmt.Errorf("cannot get hostname: %v; %v", err1, err2))
		}
		return id
	}
	hw := md5.New()
	hw.Write([]byte(hostname))
	copy(id, hw.Sum(nil))
	return id
}

func readRandomUint32() uint32 {
	var b [4]byte
	_, err := io.ReadFull(rand.Reader, b[:])
	if err != nil {
		panic(fmt.Errorf("cannot read random object id: %v", err))
	}
	return uint32((uint32(b[0]) << 0) | (uint32(b[1]) << 8) | (uint32(b[2]) << 16) | (uint32(b[3]) << 24))
}

func NewMID() XID {
	return NewMIDWithTime(time.Now())
}

func NewMIDWithTime(t time.Time) XID {
	var tUnix = uint32(t.Unix())
	var b [6]byte
	b[0] = byte(tUnix >> 16)
	b[1] = byte(tUnix >> 8)
	b[2] = byte(tUnix)

	i := atomic.AddUint32(&objectIdCounter, 1)
	b[3] = byte(i >> 16)
	b[4] = byte(i >> 8)
	b[5] = byte(i)
	return XID(string(b[:]))
}
package objects

import (
	"encoding/binary"
	"math"

	"github.com/steve-nzr/goff/internal/domain/customtypes"
)

type FPReader struct {
	data   []byte
	offset int
}

func (fp *FPReader) Initialize(data []byte) *FPReader {
	fp.data = data
	return fp
}

func (fp *FPReader) ReadString() string {
	ln := fp.ReadInt()
	return (string)(fp.ReadBytes(ln))
}

func (fp *FPReader) ReadInt() int {
	return (int)(fp.ReadUInt32())
}

func (fp *FPReader) ReadInt32() int32 {
	return (int32)(fp.ReadUInt32())
}

func (fp *FPReader) ReadUInt32() uint32 {
	i := binary.LittleEndian.Uint32(fp.data[fp.offset:])
	fp.offset += 4
	return i
}

func (fp *FPReader) ReadInt16() int16 {
	return (int16)(fp.ReadUInt32())
}

func (fp *FPReader) ReadUInt16() uint16 {
	i := binary.LittleEndian.Uint16(fp.data[fp.offset:])
	fp.offset += 2
	return i
}

func (fp *FPReader) ReadByte() byte {
	b := fp.data[fp.offset]
	fp.offset++
	return b
}

func (fp *FPReader) ReadBytes(n int) []byte {
	data := fp.data[fp.offset : fp.offset+n]
	fp.offset += n
	return data
}

type FPWriter struct {
	data   []byte
	offset int
}

func (fp *FPWriter) initializeEmptyWithSize(size int) *FPWriter {
	fp.offset = 0
	fp.data = make([]byte, size)
	return fp
}

func (fp *FPWriter) initializeEmpty() *FPWriter {
	return fp.initializeEmptyWithSize(4096)
}

func (fp *FPWriter) Initialize() *FPWriter {
	return fp.initializeEmpty().
		WriteByte(0x5E). // header
		WriteInt32(0)    // size
}

func (fp *FPWriter) finalizeWithoutLen() []byte {
	return fp.data[:fp.offset]
}

func (fp *FPWriter) finalize() []byte {
	savedOffset := fp.offset
	fp.offset = 1
	fp.WriteInt(savedOffset - 5)
	fp.offset = savedOffset
	return fp.data[:fp.offset]
}

func (fp *FPWriter) WriteString(s string) *FPWriter {
	return fp.WriteInt(len(s)).WriteBytes(([]byte)(s))
}

func (fp *FPWriter) WriteID(id customtypes.ID) *FPWriter {
	return fp.WriteUInt32((uint32)(id))
}

func (fp *FPWriter) WriteInt(i int) *FPWriter {
	return fp.WriteUInt32((uint32)(i))
}

func (fp *FPWriter) WriteInt32(i int32) *FPWriter {
	return fp.WriteUInt32((uint32)(i))
}

func (fp *FPWriter) WriteUInt32(i uint32) *FPWriter {
	binary.LittleEndian.PutUint32(fp.data[fp.offset:], i)
	fp.offset += 4
	return fp
}

func (fp *FPWriter) WriteUInt64(i uint64) *FPWriter {
	binary.LittleEndian.PutUint64(fp.data[fp.offset:], i)
	fp.offset += 8
	return fp
}

func (fp *FPWriter) WriteInt16(i int16) *FPWriter {
	return fp.WriteUInt16((uint16)(i))
}

func (fp *FPWriter) WriteUInt16(i uint16) *FPWriter {
	binary.LittleEndian.PutUint16(fp.data[fp.offset:], i)
	fp.offset += 2
	return fp
}

func (fp *FPWriter) WriteInt8(i int8) *FPWriter {
	return fp.WriteByte((byte)(i))
}

func (fp *FPWriter) WriteByte(i byte) *FPWriter {
	fp.data[fp.offset] = i
	fp.offset++
	return fp
}

func (fp *FPWriter) WriteBytes(data []byte) *FPWriter {
	for i := range data {
		fp.data[fp.offset+i] = data[i]
	}
	fp.offset += len(data)
	return fp
}

func (fp *FPWriter) WriteFloat32(f float32) *FPWriter {
	return fp.WriteUInt32(math.Float32bits(f))
}

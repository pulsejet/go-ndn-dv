// @generated by the gondn_tlv_gen, DO NOT modify manually
package tlv

import (
	"encoding/binary"
	"io"

	enc "github.com/zjkmxy/go-ndn/pkg/encoding"
)

type PacketEncoder struct {
	length uint

	Advertisement_encoder AdvertisementEncoder
	PrefixOpList_encoder  PrefixOpListEncoder
}

type PacketParsingContext struct {
	Advertisement_context AdvertisementParsingContext
	PrefixOpList_context  PrefixOpListParsingContext
}

func (encoder *PacketEncoder) Init(value *Packet) {
	if value.Advertisement != nil {
		encoder.Advertisement_encoder.Init(value.Advertisement)
	}
	if value.PrefixOpList != nil {
		encoder.PrefixOpList_encoder.Init(value.PrefixOpList)
	}

	l := uint(0)
	if value.Advertisement != nil {
		l += 1
		switch x := encoder.Advertisement_encoder.length; {
		case x <= 0xfc:
			l += 1
		case x <= 0xffff:
			l += 3
		case x <= 0xffffffff:
			l += 5
		default:
			l += 9
		}
		l += encoder.Advertisement_encoder.length
	}
	if value.PrefixOpList != nil {
		l += 3
		switch x := encoder.PrefixOpList_encoder.length; {
		case x <= 0xfc:
			l += 1
		case x <= 0xffff:
			l += 3
		case x <= 0xffffffff:
			l += 5
		default:
			l += 9
		}
		l += encoder.PrefixOpList_encoder.length
	}
	encoder.length = l

}

func (context *PacketParsingContext) Init() {
	context.Advertisement_context.Init()
	context.PrefixOpList_context.Init()
}

func (encoder *PacketEncoder) EncodeInto(value *Packet, buf []byte) {

	pos := uint(0)

	if value.Advertisement != nil {
		buf[pos] = byte(201)
		pos += 1
		switch x := encoder.Advertisement_encoder.length; {
		case x <= 0xfc:
			buf[pos] = byte(x)
			pos += 1
		case x <= 0xffff:
			buf[pos] = 0xfd
			binary.BigEndian.PutUint16(buf[pos+1:], uint16(x))
			pos += 3
		case x <= 0xffffffff:
			buf[pos] = 0xfe
			binary.BigEndian.PutUint32(buf[pos+1:], uint32(x))
			pos += 5
		default:
			buf[pos] = 0xff
			binary.BigEndian.PutUint64(buf[pos+1:], uint64(x))
			pos += 9
		}
		if encoder.Advertisement_encoder.length > 0 {
			encoder.Advertisement_encoder.EncodeInto(value.Advertisement, buf[pos:])
			pos += encoder.Advertisement_encoder.length
		}
	}
	if value.PrefixOpList != nil {
		buf[pos] = 253
		binary.BigEndian.PutUint16(buf[pos+1:], uint16(301))
		pos += 3
		switch x := encoder.PrefixOpList_encoder.length; {
		case x <= 0xfc:
			buf[pos] = byte(x)
			pos += 1
		case x <= 0xffff:
			buf[pos] = 0xfd
			binary.BigEndian.PutUint16(buf[pos+1:], uint16(x))
			pos += 3
		case x <= 0xffffffff:
			buf[pos] = 0xfe
			binary.BigEndian.PutUint32(buf[pos+1:], uint32(x))
			pos += 5
		default:
			buf[pos] = 0xff
			binary.BigEndian.PutUint64(buf[pos+1:], uint64(x))
			pos += 9
		}
		if encoder.PrefixOpList_encoder.length > 0 {
			encoder.PrefixOpList_encoder.EncodeInto(value.PrefixOpList, buf[pos:])
			pos += encoder.PrefixOpList_encoder.length
		}
	}
}

func (encoder *PacketEncoder) Encode(value *Packet) enc.Wire {

	wire := make(enc.Wire, 1)
	wire[0] = make([]byte, encoder.length)
	buf := wire[0]
	encoder.EncodeInto(value, buf)

	return wire
}

func (context *PacketParsingContext) Parse(reader enc.ParseReader, ignoreCritical bool) (*Packet, error) {
	if reader == nil {
		return nil, enc.ErrBufferOverflow
	}

	var handled_Advertisement bool = false
	var handled_PrefixOpList bool = false

	progress := -1
	_ = progress

	value := &Packet{}
	var err error
	var startPos int
	for {
		startPos = reader.Pos()
		if startPos >= reader.Length() {
			break
		}
		typ := enc.TLNum(0)
		l := enc.TLNum(0)
		typ, err = enc.ReadTLNum(reader)
		if err != nil {
			return nil, enc.ErrFailToParse{TypeNum: 0, Err: err}
		}
		l, err = enc.ReadTLNum(reader)
		if err != nil {
			return nil, enc.ErrFailToParse{TypeNum: 0, Err: err}
		}

		err = nil
		if handled := false; true {
			switch typ {
			case 201:
				if true {
					handled = true
					handled_Advertisement = true
					value.Advertisement, err = context.Advertisement_context.Parse(reader.Delegate(int(l)), ignoreCritical)
				}
			case 301:
				if true {
					handled = true
					handled_PrefixOpList = true
					value.PrefixOpList, err = context.PrefixOpList_context.Parse(reader.Delegate(int(l)), ignoreCritical)
				}
			default:
				if !ignoreCritical && ((typ <= 31) || ((typ & 1) == 1)) {
					return nil, enc.ErrUnrecognizedField{TypeNum: typ}
				}
				handled = true
				err = reader.Skip(int(l))
			}
			if err == nil && !handled {
			}
			if err != nil {
				return nil, enc.ErrFailToParse{TypeNum: typ, Err: err}
			}
		}
	}

	startPos = reader.Pos()
	err = nil

	if !handled_Advertisement && err == nil {
		value.Advertisement = nil
	}
	if !handled_PrefixOpList && err == nil {
		value.PrefixOpList = nil
	}

	if err != nil {
		return nil, err
	}

	return value, nil
}

func (value *Packet) Encode() enc.Wire {
	encoder := PacketEncoder{}
	encoder.Init(value)
	return encoder.Encode(value)
}

func (value *Packet) Bytes() []byte {
	return value.Encode().Join()
}

func ParsePacket(reader enc.ParseReader, ignoreCritical bool) (*Packet, error) {
	context := PacketParsingContext{}
	context.Init()
	return context.Parse(reader, ignoreCritical)
}

type AdvertisementEncoder struct {
	length uint

	Entries_subencoder []struct {
		Entries_encoder AdvEntryEncoder
	}
}

type AdvertisementParsingContext struct {
	Entries_context AdvEntryParsingContext
}

func (encoder *AdvertisementEncoder) Init(value *Advertisement) {
	{
		Entries_l := len(value.Entries)
		encoder.Entries_subencoder = make([]struct {
			Entries_encoder AdvEntryEncoder
		}, Entries_l)
		for i := 0; i < Entries_l; i++ {
			pseudoEncoder := &encoder.Entries_subencoder[i]
			pseudoValue := struct {
				Entries *AdvEntry
			}{
				Entries: value.Entries[i],
			}
			{
				encoder := pseudoEncoder
				value := &pseudoValue
				if value.Entries != nil {
					encoder.Entries_encoder.Init(value.Entries)
				}
				_ = encoder
				_ = value
			}
		}
	}

	l := uint(0)
	if value.Entries != nil {
		for seq_i, seq_v := range value.Entries {
			pseudoEncoder := &encoder.Entries_subencoder[seq_i]
			pseudoValue := struct {
				Entries *AdvEntry
			}{
				Entries: seq_v,
			}
			{
				encoder := pseudoEncoder
				value := &pseudoValue
				if value.Entries != nil {
					l += 1
					switch x := encoder.Entries_encoder.length; {
					case x <= 0xfc:
						l += 1
					case x <= 0xffff:
						l += 3
					case x <= 0xffffffff:
						l += 5
					default:
						l += 9
					}
					l += encoder.Entries_encoder.length
				}
				_ = encoder
				_ = value
			}
		}
	}
	encoder.length = l

}

func (context *AdvertisementParsingContext) Init() {
	context.Entries_context.Init()
}

func (encoder *AdvertisementEncoder) EncodeInto(value *Advertisement, buf []byte) {

	pos := uint(0)

	if value.Entries != nil {
		for seq_i, seq_v := range value.Entries {
			pseudoEncoder := &encoder.Entries_subencoder[seq_i]
			pseudoValue := struct {
				Entries *AdvEntry
			}{
				Entries: seq_v,
			}
			{
				encoder := pseudoEncoder
				value := &pseudoValue
				if value.Entries != nil {
					buf[pos] = byte(202)
					pos += 1
					switch x := encoder.Entries_encoder.length; {
					case x <= 0xfc:
						buf[pos] = byte(x)
						pos += 1
					case x <= 0xffff:
						buf[pos] = 0xfd
						binary.BigEndian.PutUint16(buf[pos+1:], uint16(x))
						pos += 3
					case x <= 0xffffffff:
						buf[pos] = 0xfe
						binary.BigEndian.PutUint32(buf[pos+1:], uint32(x))
						pos += 5
					default:
						buf[pos] = 0xff
						binary.BigEndian.PutUint64(buf[pos+1:], uint64(x))
						pos += 9
					}
					if encoder.Entries_encoder.length > 0 {
						encoder.Entries_encoder.EncodeInto(value.Entries, buf[pos:])
						pos += encoder.Entries_encoder.length
					}
				}
				_ = encoder
				_ = value
			}
		}
	}
}

func (encoder *AdvertisementEncoder) Encode(value *Advertisement) enc.Wire {

	wire := make(enc.Wire, 1)
	wire[0] = make([]byte, encoder.length)
	buf := wire[0]
	encoder.EncodeInto(value, buf)

	return wire
}

func (context *AdvertisementParsingContext) Parse(reader enc.ParseReader, ignoreCritical bool) (*Advertisement, error) {
	if reader == nil {
		return nil, enc.ErrBufferOverflow
	}

	var handled_Entries bool = false

	progress := -1
	_ = progress

	value := &Advertisement{}
	var err error
	var startPos int
	for {
		startPos = reader.Pos()
		if startPos >= reader.Length() {
			break
		}
		typ := enc.TLNum(0)
		l := enc.TLNum(0)
		typ, err = enc.ReadTLNum(reader)
		if err != nil {
			return nil, enc.ErrFailToParse{TypeNum: 0, Err: err}
		}
		l, err = enc.ReadTLNum(reader)
		if err != nil {
			return nil, enc.ErrFailToParse{TypeNum: 0, Err: err}
		}

		err = nil
		if handled := false; true {
			switch typ {
			case 202:
				if true {
					handled = true
					handled_Entries = true
					if value.Entries == nil {
						value.Entries = make([]*AdvEntry, 0)
					}
					{
						pseudoValue := struct {
							Entries *AdvEntry
						}{}
						{
							value := &pseudoValue
							value.Entries, err = context.Entries_context.Parse(reader.Delegate(int(l)), ignoreCritical)
							_ = value
						}
						value.Entries = append(value.Entries, pseudoValue.Entries)
					}
					progress--
				}
			default:
				if !ignoreCritical && ((typ <= 31) || ((typ & 1) == 1)) {
					return nil, enc.ErrUnrecognizedField{TypeNum: typ}
				}
				handled = true
				err = reader.Skip(int(l))
			}
			if err == nil && !handled {
			}
			if err != nil {
				return nil, enc.ErrFailToParse{TypeNum: typ, Err: err}
			}
		}
	}

	startPos = reader.Pos()
	err = nil

	if !handled_Entries && err == nil {
		// sequence - skip
	}

	if err != nil {
		return nil, err
	}

	return value, nil
}

func (value *Advertisement) Encode() enc.Wire {
	encoder := AdvertisementEncoder{}
	encoder.Init(value)
	return encoder.Encode(value)
}

func (value *Advertisement) Bytes() []byte {
	return value.Encode().Join()
}

func ParseAdvertisement(reader enc.ParseReader, ignoreCritical bool) (*Advertisement, error) {
	context := AdvertisementParsingContext{}
	context.Init()
	return context.Parse(reader, ignoreCritical)
}

type AdvEntryEncoder struct {
	length uint

	Destination_encoder DestinationEncoder
	NextHop_encoder     DestinationEncoder
}

type AdvEntryParsingContext struct {
	Destination_context DestinationParsingContext
	NextHop_context     DestinationParsingContext
}

func (encoder *AdvEntryEncoder) Init(value *AdvEntry) {
	if value.Destination != nil {
		encoder.Destination_encoder.Init(value.Destination)
	}
	if value.NextHop != nil {
		encoder.NextHop_encoder.Init(value.NextHop)
	}

	l := uint(0)
	if value.Destination != nil {
		l += 1
		switch x := encoder.Destination_encoder.length; {
		case x <= 0xfc:
			l += 1
		case x <= 0xffff:
			l += 3
		case x <= 0xffffffff:
			l += 5
		default:
			l += 9
		}
		l += encoder.Destination_encoder.length
	}
	if value.NextHop != nil {
		l += 1
		switch x := encoder.NextHop_encoder.length; {
		case x <= 0xfc:
			l += 1
		case x <= 0xffff:
			l += 3
		case x <= 0xffffffff:
			l += 5
		default:
			l += 9
		}
		l += encoder.NextHop_encoder.length
	}
	l += 1
	switch x := value.Cost; {
	case x <= 0xff:
		l += 2
	case x <= 0xffff:
		l += 3
	case x <= 0xffffffff:
		l += 5
	default:
		l += 9
	}
	l += 1
	switch x := value.OtherCost; {
	case x <= 0xff:
		l += 2
	case x <= 0xffff:
		l += 3
	case x <= 0xffffffff:
		l += 5
	default:
		l += 9
	}
	encoder.length = l

}

func (context *AdvEntryParsingContext) Init() {
	context.Destination_context.Init()
	context.NextHop_context.Init()

}

func (encoder *AdvEntryEncoder) EncodeInto(value *AdvEntry, buf []byte) {

	pos := uint(0)

	if value.Destination != nil {
		buf[pos] = byte(204)
		pos += 1
		switch x := encoder.Destination_encoder.length; {
		case x <= 0xfc:
			buf[pos] = byte(x)
			pos += 1
		case x <= 0xffff:
			buf[pos] = 0xfd
			binary.BigEndian.PutUint16(buf[pos+1:], uint16(x))
			pos += 3
		case x <= 0xffffffff:
			buf[pos] = 0xfe
			binary.BigEndian.PutUint32(buf[pos+1:], uint32(x))
			pos += 5
		default:
			buf[pos] = 0xff
			binary.BigEndian.PutUint64(buf[pos+1:], uint64(x))
			pos += 9
		}
		if encoder.Destination_encoder.length > 0 {
			encoder.Destination_encoder.EncodeInto(value.Destination, buf[pos:])
			pos += encoder.Destination_encoder.length
		}
	}
	if value.NextHop != nil {
		buf[pos] = byte(206)
		pos += 1
		switch x := encoder.NextHop_encoder.length; {
		case x <= 0xfc:
			buf[pos] = byte(x)
			pos += 1
		case x <= 0xffff:
			buf[pos] = 0xfd
			binary.BigEndian.PutUint16(buf[pos+1:], uint16(x))
			pos += 3
		case x <= 0xffffffff:
			buf[pos] = 0xfe
			binary.BigEndian.PutUint32(buf[pos+1:], uint32(x))
			pos += 5
		default:
			buf[pos] = 0xff
			binary.BigEndian.PutUint64(buf[pos+1:], uint64(x))
			pos += 9
		}
		if encoder.NextHop_encoder.length > 0 {
			encoder.NextHop_encoder.EncodeInto(value.NextHop, buf[pos:])
			pos += encoder.NextHop_encoder.length
		}
	}
	buf[pos] = byte(208)
	pos += 1
	switch x := value.Cost; {
	case x <= 0xff:
		buf[pos] = 1
		buf[pos+1] = byte(x)
		pos += 2
	case x <= 0xffff:
		buf[pos] = 2
		binary.BigEndian.PutUint16(buf[pos+1:], uint16(x))
		pos += 3
	case x <= 0xffffffff:
		buf[pos] = 4
		binary.BigEndian.PutUint32(buf[pos+1:], uint32(x))
		pos += 5
	default:
		buf[pos] = 8
		binary.BigEndian.PutUint64(buf[pos+1:], uint64(x))
		pos += 9
	}
	buf[pos] = byte(210)
	pos += 1
	switch x := value.OtherCost; {
	case x <= 0xff:
		buf[pos] = 1
		buf[pos+1] = byte(x)
		pos += 2
	case x <= 0xffff:
		buf[pos] = 2
		binary.BigEndian.PutUint16(buf[pos+1:], uint16(x))
		pos += 3
	case x <= 0xffffffff:
		buf[pos] = 4
		binary.BigEndian.PutUint32(buf[pos+1:], uint32(x))
		pos += 5
	default:
		buf[pos] = 8
		binary.BigEndian.PutUint64(buf[pos+1:], uint64(x))
		pos += 9
	}
}

func (encoder *AdvEntryEncoder) Encode(value *AdvEntry) enc.Wire {

	wire := make(enc.Wire, 1)
	wire[0] = make([]byte, encoder.length)
	buf := wire[0]
	encoder.EncodeInto(value, buf)

	return wire
}

func (context *AdvEntryParsingContext) Parse(reader enc.ParseReader, ignoreCritical bool) (*AdvEntry, error) {
	if reader == nil {
		return nil, enc.ErrBufferOverflow
	}

	var handled_Destination bool = false
	var handled_NextHop bool = false
	var handled_Cost bool = false
	var handled_OtherCost bool = false

	progress := -1
	_ = progress

	value := &AdvEntry{}
	var err error
	var startPos int
	for {
		startPos = reader.Pos()
		if startPos >= reader.Length() {
			break
		}
		typ := enc.TLNum(0)
		l := enc.TLNum(0)
		typ, err = enc.ReadTLNum(reader)
		if err != nil {
			return nil, enc.ErrFailToParse{TypeNum: 0, Err: err}
		}
		l, err = enc.ReadTLNum(reader)
		if err != nil {
			return nil, enc.ErrFailToParse{TypeNum: 0, Err: err}
		}

		err = nil
		if handled := false; true {
			switch typ {
			case 204:
				if true {
					handled = true
					handled_Destination = true
					value.Destination, err = context.Destination_context.Parse(reader.Delegate(int(l)), ignoreCritical)
				}
			case 206:
				if true {
					handled = true
					handled_NextHop = true
					value.NextHop, err = context.NextHop_context.Parse(reader.Delegate(int(l)), ignoreCritical)
				}
			case 208:
				if true {
					handled = true
					handled_Cost = true
					value.Cost = uint64(0)
					{
						for i := 0; i < int(l); i++ {
							x := byte(0)
							x, err = reader.ReadByte()
							if err != nil {
								if err == io.EOF {
									err = io.ErrUnexpectedEOF
								}
								break
							}
							value.Cost = uint64(value.Cost<<8) | uint64(x)
						}
					}
				}
			case 210:
				if true {
					handled = true
					handled_OtherCost = true
					value.OtherCost = uint64(0)
					{
						for i := 0; i < int(l); i++ {
							x := byte(0)
							x, err = reader.ReadByte()
							if err != nil {
								if err == io.EOF {
									err = io.ErrUnexpectedEOF
								}
								break
							}
							value.OtherCost = uint64(value.OtherCost<<8) | uint64(x)
						}
					}
				}
			default:
				if !ignoreCritical && ((typ <= 31) || ((typ & 1) == 1)) {
					return nil, enc.ErrUnrecognizedField{TypeNum: typ}
				}
				handled = true
				err = reader.Skip(int(l))
			}
			if err == nil && !handled {
			}
			if err != nil {
				return nil, enc.ErrFailToParse{TypeNum: typ, Err: err}
			}
		}
	}

	startPos = reader.Pos()
	err = nil

	if !handled_Destination && err == nil {
		value.Destination = nil
	}
	if !handled_NextHop && err == nil {
		value.NextHop = nil
	}
	if !handled_Cost && err == nil {
		err = enc.ErrSkipRequired{Name: "Cost", TypeNum: 208}
	}
	if !handled_OtherCost && err == nil {
		err = enc.ErrSkipRequired{Name: "OtherCost", TypeNum: 210}
	}

	if err != nil {
		return nil, err
	}

	return value, nil
}

func (value *AdvEntry) Encode() enc.Wire {
	encoder := AdvEntryEncoder{}
	encoder.Init(value)
	return encoder.Encode(value)
}

func (value *AdvEntry) Bytes() []byte {
	return value.Encode().Join()
}

func ParseAdvEntry(reader enc.ParseReader, ignoreCritical bool) (*AdvEntry, error) {
	context := AdvEntryParsingContext{}
	context.Init()
	return context.Parse(reader, ignoreCritical)
}

type DestinationEncoder struct {
	length uint

	Name_length uint
}

type DestinationParsingContext struct {
}

func (encoder *DestinationEncoder) Init(value *Destination) {
	if value.Name != nil {
		encoder.Name_length = 0
		for _, c := range value.Name {
			encoder.Name_length += uint(c.EncodingLength())
		}
	}

	l := uint(0)
	if value.Name != nil {
		l += 1
		switch x := encoder.Name_length; {
		case x <= 0xfc:
			l += 1
		case x <= 0xffff:
			l += 3
		case x <= 0xffffffff:
			l += 5
		default:
			l += 9
		}
		l += encoder.Name_length
	}
	encoder.length = l

}

func (context *DestinationParsingContext) Init() {

}

func (encoder *DestinationEncoder) EncodeInto(value *Destination, buf []byte) {

	pos := uint(0)

	if value.Name != nil {
		buf[pos] = byte(7)
		pos += 1
		switch x := encoder.Name_length; {
		case x <= 0xfc:
			buf[pos] = byte(x)
			pos += 1
		case x <= 0xffff:
			buf[pos] = 0xfd
			binary.BigEndian.PutUint16(buf[pos+1:], uint16(x))
			pos += 3
		case x <= 0xffffffff:
			buf[pos] = 0xfe
			binary.BigEndian.PutUint32(buf[pos+1:], uint32(x))
			pos += 5
		default:
			buf[pos] = 0xff
			binary.BigEndian.PutUint64(buf[pos+1:], uint64(x))
			pos += 9
		}
		for _, c := range value.Name {
			pos += uint(c.EncodeInto(buf[pos:]))
		}
	}
}

func (encoder *DestinationEncoder) Encode(value *Destination) enc.Wire {

	wire := make(enc.Wire, 1)
	wire[0] = make([]byte, encoder.length)
	buf := wire[0]
	encoder.EncodeInto(value, buf)

	return wire
}

func (context *DestinationParsingContext) Parse(reader enc.ParseReader, ignoreCritical bool) (*Destination, error) {
	if reader == nil {
		return nil, enc.ErrBufferOverflow
	}

	var handled_Name bool = false

	progress := -1
	_ = progress

	value := &Destination{}
	var err error
	var startPos int
	for {
		startPos = reader.Pos()
		if startPos >= reader.Length() {
			break
		}
		typ := enc.TLNum(0)
		l := enc.TLNum(0)
		typ, err = enc.ReadTLNum(reader)
		if err != nil {
			return nil, enc.ErrFailToParse{TypeNum: 0, Err: err}
		}
		l, err = enc.ReadTLNum(reader)
		if err != nil {
			return nil, enc.ErrFailToParse{TypeNum: 0, Err: err}
		}

		err = nil
		if handled := false; true {
			switch typ {
			case 7:
				if true {
					handled = true
					handled_Name = true
					value.Name = make(enc.Name, l/2+1)
					startName := reader.Pos()
					endName := startName + int(l)
					for j := range value.Name {
						if reader.Pos() >= endName {
							value.Name = value.Name[:j]
							break
						}
						var err1, err3 error
						value.Name[j].Typ, err1 = enc.ReadTLNum(reader)
						l, err2 := enc.ReadTLNum(reader)
						value.Name[j].Val, err3 = reader.ReadBuf(int(l))
						if err1 != nil || err2 != nil || err3 != nil {
							err = io.ErrUnexpectedEOF
							break
						}
					}
					if err == nil && reader.Pos() != endName {
						err = enc.ErrBufferOverflow
					}
				}
			default:
				if !ignoreCritical && ((typ <= 31) || ((typ & 1) == 1)) {
					return nil, enc.ErrUnrecognizedField{TypeNum: typ}
				}
				handled = true
				err = reader.Skip(int(l))
			}
			if err == nil && !handled {
			}
			if err != nil {
				return nil, enc.ErrFailToParse{TypeNum: typ, Err: err}
			}
		}
	}

	startPos = reader.Pos()
	err = nil

	if !handled_Name && err == nil {
		value.Name = nil
	}

	if err != nil {
		return nil, err
	}

	return value, nil
}

func (value *Destination) Encode() enc.Wire {
	encoder := DestinationEncoder{}
	encoder.Init(value)
	return encoder.Encode(value)
}

func (value *Destination) Bytes() []byte {
	return value.Encode().Join()
}

func ParseDestination(reader enc.ParseReader, ignoreCritical bool) (*Destination, error) {
	context := DestinationParsingContext{}
	context.Init()
	return context.Parse(reader, ignoreCritical)
}

type PrefixOpListEncoder struct {
	length uint

	ExitRouter_encoder DestinationEncoder

	PrefixOpAdds_subencoder []struct {
		PrefixOpAdds_encoder PrefixOpAddEncoder
	}
	PrefixOpRemoves_subencoder []struct {
		PrefixOpRemoves_encoder PrefixOpRemoveEncoder
	}
}

type PrefixOpListParsingContext struct {
	ExitRouter_context DestinationParsingContext

	PrefixOpAdds_context    PrefixOpAddParsingContext
	PrefixOpRemoves_context PrefixOpRemoveParsingContext
}

func (encoder *PrefixOpListEncoder) Init(value *PrefixOpList) {
	if value.ExitRouter != nil {
		encoder.ExitRouter_encoder.Init(value.ExitRouter)
	}

	{
		PrefixOpAdds_l := len(value.PrefixOpAdds)
		encoder.PrefixOpAdds_subencoder = make([]struct {
			PrefixOpAdds_encoder PrefixOpAddEncoder
		}, PrefixOpAdds_l)
		for i := 0; i < PrefixOpAdds_l; i++ {
			pseudoEncoder := &encoder.PrefixOpAdds_subencoder[i]
			pseudoValue := struct {
				PrefixOpAdds *PrefixOpAdd
			}{
				PrefixOpAdds: value.PrefixOpAdds[i],
			}
			{
				encoder := pseudoEncoder
				value := &pseudoValue
				if value.PrefixOpAdds != nil {
					encoder.PrefixOpAdds_encoder.Init(value.PrefixOpAdds)
				}
				_ = encoder
				_ = value
			}
		}
	}
	{
		PrefixOpRemoves_l := len(value.PrefixOpRemoves)
		encoder.PrefixOpRemoves_subencoder = make([]struct {
			PrefixOpRemoves_encoder PrefixOpRemoveEncoder
		}, PrefixOpRemoves_l)
		for i := 0; i < PrefixOpRemoves_l; i++ {
			pseudoEncoder := &encoder.PrefixOpRemoves_subencoder[i]
			pseudoValue := struct {
				PrefixOpRemoves *PrefixOpRemove
			}{
				PrefixOpRemoves: value.PrefixOpRemoves[i],
			}
			{
				encoder := pseudoEncoder
				value := &pseudoValue
				if value.PrefixOpRemoves != nil {
					encoder.PrefixOpRemoves_encoder.Init(value.PrefixOpRemoves)
				}
				_ = encoder
				_ = value
			}
		}
	}

	l := uint(0)
	if value.ExitRouter != nil {
		l += 1
		switch x := encoder.ExitRouter_encoder.length; {
		case x <= 0xfc:
			l += 1
		case x <= 0xffff:
			l += 3
		case x <= 0xffffffff:
			l += 5
		default:
			l += 9
		}
		l += encoder.ExitRouter_encoder.length
	}
	if value.PrefixOpReset {
		l += 3
		l += 1
	}
	if value.PrefixOpAdds != nil {
		for seq_i, seq_v := range value.PrefixOpAdds {
			pseudoEncoder := &encoder.PrefixOpAdds_subencoder[seq_i]
			pseudoValue := struct {
				PrefixOpAdds *PrefixOpAdd
			}{
				PrefixOpAdds: seq_v,
			}
			{
				encoder := pseudoEncoder
				value := &pseudoValue
				if value.PrefixOpAdds != nil {
					l += 3
					switch x := encoder.PrefixOpAdds_encoder.length; {
					case x <= 0xfc:
						l += 1
					case x <= 0xffff:
						l += 3
					case x <= 0xffffffff:
						l += 5
					default:
						l += 9
					}
					l += encoder.PrefixOpAdds_encoder.length
				}
				_ = encoder
				_ = value
			}
		}
	}
	if value.PrefixOpRemoves != nil {
		for seq_i, seq_v := range value.PrefixOpRemoves {
			pseudoEncoder := &encoder.PrefixOpRemoves_subencoder[seq_i]
			pseudoValue := struct {
				PrefixOpRemoves *PrefixOpRemove
			}{
				PrefixOpRemoves: seq_v,
			}
			{
				encoder := pseudoEncoder
				value := &pseudoValue
				if value.PrefixOpRemoves != nil {
					l += 3
					switch x := encoder.PrefixOpRemoves_encoder.length; {
					case x <= 0xfc:
						l += 1
					case x <= 0xffff:
						l += 3
					case x <= 0xffffffff:
						l += 5
					default:
						l += 9
					}
					l += encoder.PrefixOpRemoves_encoder.length
				}
				_ = encoder
				_ = value
			}
		}
	}
	encoder.length = l

}

func (context *PrefixOpListParsingContext) Init() {
	context.ExitRouter_context.Init()

	context.PrefixOpAdds_context.Init()
	context.PrefixOpRemoves_context.Init()
}

func (encoder *PrefixOpListEncoder) EncodeInto(value *PrefixOpList, buf []byte) {

	pos := uint(0)

	if value.ExitRouter != nil {
		buf[pos] = byte(204)
		pos += 1
		switch x := encoder.ExitRouter_encoder.length; {
		case x <= 0xfc:
			buf[pos] = byte(x)
			pos += 1
		case x <= 0xffff:
			buf[pos] = 0xfd
			binary.BigEndian.PutUint16(buf[pos+1:], uint16(x))
			pos += 3
		case x <= 0xffffffff:
			buf[pos] = 0xfe
			binary.BigEndian.PutUint32(buf[pos+1:], uint32(x))
			pos += 5
		default:
			buf[pos] = 0xff
			binary.BigEndian.PutUint64(buf[pos+1:], uint64(x))
			pos += 9
		}
		if encoder.ExitRouter_encoder.length > 0 {
			encoder.ExitRouter_encoder.EncodeInto(value.ExitRouter, buf[pos:])
			pos += encoder.ExitRouter_encoder.length
		}
	}
	if value.PrefixOpReset {
		buf[pos] = 253
		binary.BigEndian.PutUint16(buf[pos+1:], uint16(302))
		pos += 3
		buf[pos] = byte(0)
		pos += 1
	}
	if value.PrefixOpAdds != nil {
		for seq_i, seq_v := range value.PrefixOpAdds {
			pseudoEncoder := &encoder.PrefixOpAdds_subencoder[seq_i]
			pseudoValue := struct {
				PrefixOpAdds *PrefixOpAdd
			}{
				PrefixOpAdds: seq_v,
			}
			{
				encoder := pseudoEncoder
				value := &pseudoValue
				if value.PrefixOpAdds != nil {
					buf[pos] = 253
					binary.BigEndian.PutUint16(buf[pos+1:], uint16(304))
					pos += 3
					switch x := encoder.PrefixOpAdds_encoder.length; {
					case x <= 0xfc:
						buf[pos] = byte(x)
						pos += 1
					case x <= 0xffff:
						buf[pos] = 0xfd
						binary.BigEndian.PutUint16(buf[pos+1:], uint16(x))
						pos += 3
					case x <= 0xffffffff:
						buf[pos] = 0xfe
						binary.BigEndian.PutUint32(buf[pos+1:], uint32(x))
						pos += 5
					default:
						buf[pos] = 0xff
						binary.BigEndian.PutUint64(buf[pos+1:], uint64(x))
						pos += 9
					}
					if encoder.PrefixOpAdds_encoder.length > 0 {
						encoder.PrefixOpAdds_encoder.EncodeInto(value.PrefixOpAdds, buf[pos:])
						pos += encoder.PrefixOpAdds_encoder.length
					}
				}
				_ = encoder
				_ = value
			}
		}
	}
	if value.PrefixOpRemoves != nil {
		for seq_i, seq_v := range value.PrefixOpRemoves {
			pseudoEncoder := &encoder.PrefixOpRemoves_subencoder[seq_i]
			pseudoValue := struct {
				PrefixOpRemoves *PrefixOpRemove
			}{
				PrefixOpRemoves: seq_v,
			}
			{
				encoder := pseudoEncoder
				value := &pseudoValue
				if value.PrefixOpRemoves != nil {
					buf[pos] = 253
					binary.BigEndian.PutUint16(buf[pos+1:], uint16(306))
					pos += 3
					switch x := encoder.PrefixOpRemoves_encoder.length; {
					case x <= 0xfc:
						buf[pos] = byte(x)
						pos += 1
					case x <= 0xffff:
						buf[pos] = 0xfd
						binary.BigEndian.PutUint16(buf[pos+1:], uint16(x))
						pos += 3
					case x <= 0xffffffff:
						buf[pos] = 0xfe
						binary.BigEndian.PutUint32(buf[pos+1:], uint32(x))
						pos += 5
					default:
						buf[pos] = 0xff
						binary.BigEndian.PutUint64(buf[pos+1:], uint64(x))
						pos += 9
					}
					if encoder.PrefixOpRemoves_encoder.length > 0 {
						encoder.PrefixOpRemoves_encoder.EncodeInto(value.PrefixOpRemoves, buf[pos:])
						pos += encoder.PrefixOpRemoves_encoder.length
					}
				}
				_ = encoder
				_ = value
			}
		}
	}
}

func (encoder *PrefixOpListEncoder) Encode(value *PrefixOpList) enc.Wire {

	wire := make(enc.Wire, 1)
	wire[0] = make([]byte, encoder.length)
	buf := wire[0]
	encoder.EncodeInto(value, buf)

	return wire
}

func (context *PrefixOpListParsingContext) Parse(reader enc.ParseReader, ignoreCritical bool) (*PrefixOpList, error) {
	if reader == nil {
		return nil, enc.ErrBufferOverflow
	}

	var handled_ExitRouter bool = false
	var handled_PrefixOpReset bool = false
	var handled_PrefixOpAdds bool = false
	var handled_PrefixOpRemoves bool = false

	progress := -1
	_ = progress

	value := &PrefixOpList{}
	var err error
	var startPos int
	for {
		startPos = reader.Pos()
		if startPos >= reader.Length() {
			break
		}
		typ := enc.TLNum(0)
		l := enc.TLNum(0)
		typ, err = enc.ReadTLNum(reader)
		if err != nil {
			return nil, enc.ErrFailToParse{TypeNum: 0, Err: err}
		}
		l, err = enc.ReadTLNum(reader)
		if err != nil {
			return nil, enc.ErrFailToParse{TypeNum: 0, Err: err}
		}

		err = nil
		if handled := false; true {
			switch typ {
			case 204:
				if true {
					handled = true
					handled_ExitRouter = true
					value.ExitRouter, err = context.ExitRouter_context.Parse(reader.Delegate(int(l)), ignoreCritical)
				}
			case 302:
				if true {
					handled = true
					handled_PrefixOpReset = true
					value.PrefixOpReset = true
				}
			case 304:
				if true {
					handled = true
					handled_PrefixOpAdds = true
					if value.PrefixOpAdds == nil {
						value.PrefixOpAdds = make([]*PrefixOpAdd, 0)
					}
					{
						pseudoValue := struct {
							PrefixOpAdds *PrefixOpAdd
						}{}
						{
							value := &pseudoValue
							value.PrefixOpAdds, err = context.PrefixOpAdds_context.Parse(reader.Delegate(int(l)), ignoreCritical)
							_ = value
						}
						value.PrefixOpAdds = append(value.PrefixOpAdds, pseudoValue.PrefixOpAdds)
					}
					progress--
				}
			case 306:
				if true {
					handled = true
					handled_PrefixOpRemoves = true
					if value.PrefixOpRemoves == nil {
						value.PrefixOpRemoves = make([]*PrefixOpRemove, 0)
					}
					{
						pseudoValue := struct {
							PrefixOpRemoves *PrefixOpRemove
						}{}
						{
							value := &pseudoValue
							value.PrefixOpRemoves, err = context.PrefixOpRemoves_context.Parse(reader.Delegate(int(l)), ignoreCritical)
							_ = value
						}
						value.PrefixOpRemoves = append(value.PrefixOpRemoves, pseudoValue.PrefixOpRemoves)
					}
					progress--
				}
			default:
				if !ignoreCritical && ((typ <= 31) || ((typ & 1) == 1)) {
					return nil, enc.ErrUnrecognizedField{TypeNum: typ}
				}
				handled = true
				err = reader.Skip(int(l))
			}
			if err == nil && !handled {
			}
			if err != nil {
				return nil, enc.ErrFailToParse{TypeNum: typ, Err: err}
			}
		}
	}

	startPos = reader.Pos()
	err = nil

	if !handled_ExitRouter && err == nil {
		value.ExitRouter = nil
	}
	if !handled_PrefixOpReset && err == nil {
		value.PrefixOpReset = false
	}
	if !handled_PrefixOpAdds && err == nil {
		// sequence - skip
	}
	if !handled_PrefixOpRemoves && err == nil {
		// sequence - skip
	}

	if err != nil {
		return nil, err
	}

	return value, nil
}

func (value *PrefixOpList) Encode() enc.Wire {
	encoder := PrefixOpListEncoder{}
	encoder.Init(value)
	return encoder.Encode(value)
}

func (value *PrefixOpList) Bytes() []byte {
	return value.Encode().Join()
}

func ParsePrefixOpList(reader enc.ParseReader, ignoreCritical bool) (*PrefixOpList, error) {
	context := PrefixOpListParsingContext{}
	context.Init()
	return context.Parse(reader, ignoreCritical)
}

type PrefixOpAddEncoder struct {
	length uint

	Name_length uint
}

type PrefixOpAddParsingContext struct {
}

func (encoder *PrefixOpAddEncoder) Init(value *PrefixOpAdd) {
	if value.Name != nil {
		encoder.Name_length = 0
		for _, c := range value.Name {
			encoder.Name_length += uint(c.EncodingLength())
		}
	}

	l := uint(0)
	if value.Name != nil {
		l += 1
		switch x := encoder.Name_length; {
		case x <= 0xfc:
			l += 1
		case x <= 0xffff:
			l += 3
		case x <= 0xffffffff:
			l += 5
		default:
			l += 9
		}
		l += encoder.Name_length
	}
	l += 1
	switch x := value.Cost; {
	case x <= 0xff:
		l += 2
	case x <= 0xffff:
		l += 3
	case x <= 0xffffffff:
		l += 5
	default:
		l += 9
	}
	encoder.length = l

}

func (context *PrefixOpAddParsingContext) Init() {

}

func (encoder *PrefixOpAddEncoder) EncodeInto(value *PrefixOpAdd, buf []byte) {

	pos := uint(0)

	if value.Name != nil {
		buf[pos] = byte(7)
		pos += 1
		switch x := encoder.Name_length; {
		case x <= 0xfc:
			buf[pos] = byte(x)
			pos += 1
		case x <= 0xffff:
			buf[pos] = 0xfd
			binary.BigEndian.PutUint16(buf[pos+1:], uint16(x))
			pos += 3
		case x <= 0xffffffff:
			buf[pos] = 0xfe
			binary.BigEndian.PutUint32(buf[pos+1:], uint32(x))
			pos += 5
		default:
			buf[pos] = 0xff
			binary.BigEndian.PutUint64(buf[pos+1:], uint64(x))
			pos += 9
		}
		for _, c := range value.Name {
			pos += uint(c.EncodeInto(buf[pos:]))
		}
	}
	buf[pos] = byte(208)
	pos += 1
	switch x := value.Cost; {
	case x <= 0xff:
		buf[pos] = 1
		buf[pos+1] = byte(x)
		pos += 2
	case x <= 0xffff:
		buf[pos] = 2
		binary.BigEndian.PutUint16(buf[pos+1:], uint16(x))
		pos += 3
	case x <= 0xffffffff:
		buf[pos] = 4
		binary.BigEndian.PutUint32(buf[pos+1:], uint32(x))
		pos += 5
	default:
		buf[pos] = 8
		binary.BigEndian.PutUint64(buf[pos+1:], uint64(x))
		pos += 9
	}
}

func (encoder *PrefixOpAddEncoder) Encode(value *PrefixOpAdd) enc.Wire {

	wire := make(enc.Wire, 1)
	wire[0] = make([]byte, encoder.length)
	buf := wire[0]
	encoder.EncodeInto(value, buf)

	return wire
}

func (context *PrefixOpAddParsingContext) Parse(reader enc.ParseReader, ignoreCritical bool) (*PrefixOpAdd, error) {
	if reader == nil {
		return nil, enc.ErrBufferOverflow
	}

	var handled_Name bool = false
	var handled_Cost bool = false

	progress := -1
	_ = progress

	value := &PrefixOpAdd{}
	var err error
	var startPos int
	for {
		startPos = reader.Pos()
		if startPos >= reader.Length() {
			break
		}
		typ := enc.TLNum(0)
		l := enc.TLNum(0)
		typ, err = enc.ReadTLNum(reader)
		if err != nil {
			return nil, enc.ErrFailToParse{TypeNum: 0, Err: err}
		}
		l, err = enc.ReadTLNum(reader)
		if err != nil {
			return nil, enc.ErrFailToParse{TypeNum: 0, Err: err}
		}

		err = nil
		if handled := false; true {
			switch typ {
			case 7:
				if true {
					handled = true
					handled_Name = true
					value.Name = make(enc.Name, l/2+1)
					startName := reader.Pos()
					endName := startName + int(l)
					for j := range value.Name {
						if reader.Pos() >= endName {
							value.Name = value.Name[:j]
							break
						}
						var err1, err3 error
						value.Name[j].Typ, err1 = enc.ReadTLNum(reader)
						l, err2 := enc.ReadTLNum(reader)
						value.Name[j].Val, err3 = reader.ReadBuf(int(l))
						if err1 != nil || err2 != nil || err3 != nil {
							err = io.ErrUnexpectedEOF
							break
						}
					}
					if err == nil && reader.Pos() != endName {
						err = enc.ErrBufferOverflow
					}
				}
			case 208:
				if true {
					handled = true
					handled_Cost = true
					value.Cost = uint64(0)
					{
						for i := 0; i < int(l); i++ {
							x := byte(0)
							x, err = reader.ReadByte()
							if err != nil {
								if err == io.EOF {
									err = io.ErrUnexpectedEOF
								}
								break
							}
							value.Cost = uint64(value.Cost<<8) | uint64(x)
						}
					}
				}
			default:
				if !ignoreCritical && ((typ <= 31) || ((typ & 1) == 1)) {
					return nil, enc.ErrUnrecognizedField{TypeNum: typ}
				}
				handled = true
				err = reader.Skip(int(l))
			}
			if err == nil && !handled {
			}
			if err != nil {
				return nil, enc.ErrFailToParse{TypeNum: typ, Err: err}
			}
		}
	}

	startPos = reader.Pos()
	err = nil

	if !handled_Name && err == nil {
		value.Name = nil
	}
	if !handled_Cost && err == nil {
		err = enc.ErrSkipRequired{Name: "Cost", TypeNum: 208}
	}

	if err != nil {
		return nil, err
	}

	return value, nil
}

func (value *PrefixOpAdd) Encode() enc.Wire {
	encoder := PrefixOpAddEncoder{}
	encoder.Init(value)
	return encoder.Encode(value)
}

func (value *PrefixOpAdd) Bytes() []byte {
	return value.Encode().Join()
}

func ParsePrefixOpAdd(reader enc.ParseReader, ignoreCritical bool) (*PrefixOpAdd, error) {
	context := PrefixOpAddParsingContext{}
	context.Init()
	return context.Parse(reader, ignoreCritical)
}

type PrefixOpRemoveEncoder struct {
	length uint

	Name_length uint
}

type PrefixOpRemoveParsingContext struct {
}

func (encoder *PrefixOpRemoveEncoder) Init(value *PrefixOpRemove) {
	if value.Name != nil {
		encoder.Name_length = 0
		for _, c := range value.Name {
			encoder.Name_length += uint(c.EncodingLength())
		}
	}

	l := uint(0)
	if value.Name != nil {
		l += 1
		switch x := encoder.Name_length; {
		case x <= 0xfc:
			l += 1
		case x <= 0xffff:
			l += 3
		case x <= 0xffffffff:
			l += 5
		default:
			l += 9
		}
		l += encoder.Name_length
	}
	encoder.length = l

}

func (context *PrefixOpRemoveParsingContext) Init() {

}

func (encoder *PrefixOpRemoveEncoder) EncodeInto(value *PrefixOpRemove, buf []byte) {

	pos := uint(0)

	if value.Name != nil {
		buf[pos] = byte(7)
		pos += 1
		switch x := encoder.Name_length; {
		case x <= 0xfc:
			buf[pos] = byte(x)
			pos += 1
		case x <= 0xffff:
			buf[pos] = 0xfd
			binary.BigEndian.PutUint16(buf[pos+1:], uint16(x))
			pos += 3
		case x <= 0xffffffff:
			buf[pos] = 0xfe
			binary.BigEndian.PutUint32(buf[pos+1:], uint32(x))
			pos += 5
		default:
			buf[pos] = 0xff
			binary.BigEndian.PutUint64(buf[pos+1:], uint64(x))
			pos += 9
		}
		for _, c := range value.Name {
			pos += uint(c.EncodeInto(buf[pos:]))
		}
	}
}

func (encoder *PrefixOpRemoveEncoder) Encode(value *PrefixOpRemove) enc.Wire {

	wire := make(enc.Wire, 1)
	wire[0] = make([]byte, encoder.length)
	buf := wire[0]
	encoder.EncodeInto(value, buf)

	return wire
}

func (context *PrefixOpRemoveParsingContext) Parse(reader enc.ParseReader, ignoreCritical bool) (*PrefixOpRemove, error) {
	if reader == nil {
		return nil, enc.ErrBufferOverflow
	}

	var handled_Name bool = false

	progress := -1
	_ = progress

	value := &PrefixOpRemove{}
	var err error
	var startPos int
	for {
		startPos = reader.Pos()
		if startPos >= reader.Length() {
			break
		}
		typ := enc.TLNum(0)
		l := enc.TLNum(0)
		typ, err = enc.ReadTLNum(reader)
		if err != nil {
			return nil, enc.ErrFailToParse{TypeNum: 0, Err: err}
		}
		l, err = enc.ReadTLNum(reader)
		if err != nil {
			return nil, enc.ErrFailToParse{TypeNum: 0, Err: err}
		}

		err = nil
		if handled := false; true {
			switch typ {
			case 7:
				if true {
					handled = true
					handled_Name = true
					value.Name = make(enc.Name, l/2+1)
					startName := reader.Pos()
					endName := startName + int(l)
					for j := range value.Name {
						if reader.Pos() >= endName {
							value.Name = value.Name[:j]
							break
						}
						var err1, err3 error
						value.Name[j].Typ, err1 = enc.ReadTLNum(reader)
						l, err2 := enc.ReadTLNum(reader)
						value.Name[j].Val, err3 = reader.ReadBuf(int(l))
						if err1 != nil || err2 != nil || err3 != nil {
							err = io.ErrUnexpectedEOF
							break
						}
					}
					if err == nil && reader.Pos() != endName {
						err = enc.ErrBufferOverflow
					}
				}
			default:
				if !ignoreCritical && ((typ <= 31) || ((typ & 1) == 1)) {
					return nil, enc.ErrUnrecognizedField{TypeNum: typ}
				}
				handled = true
				err = reader.Skip(int(l))
			}
			if err == nil && !handled {
			}
			if err != nil {
				return nil, enc.ErrFailToParse{TypeNum: typ, Err: err}
			}
		}
	}

	startPos = reader.Pos()
	err = nil

	if !handled_Name && err == nil {
		value.Name = nil
	}

	if err != nil {
		return nil, err
	}

	return value, nil
}

func (value *PrefixOpRemove) Encode() enc.Wire {
	encoder := PrefixOpRemoveEncoder{}
	encoder.Init(value)
	return encoder.Encode(value)
}

func (value *PrefixOpRemove) Bytes() []byte {
	return value.Encode().Join()
}

func ParsePrefixOpRemove(reader enc.ParseReader, ignoreCritical bool) (*PrefixOpRemove, error) {
	context := PrefixOpRemoveParsingContext{}
	context.Init()
	return context.Parse(reader, ignoreCritical)
}

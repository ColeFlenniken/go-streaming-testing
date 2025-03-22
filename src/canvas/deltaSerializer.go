package canvas

type CanvasDelta struct {
	X     uint
	Y     uint
	Color byte
}

func DeltaDeserialize(changes []byte) []CanvasDelta {
	var numChanges int = int(changes[0])<<16 | int(changes[1])<<8 | int(changes[2])
	var output []CanvasDelta = make([]CanvasDelta, numChanges)
	var place BitCursor = BitCursor{byteN: 3, bitN: 0}

	for i := 0; i < numChanges; i++ {
		output[i] = deltaDeseralizeSingle(changes, &place)
	}
	return output
}
func deltaDeseralizeSingle(changes []byte, place *BitCursor) CanvasDelta {
	var delta = CanvasDelta{}
	for i := 0; i < 12; i++ {
		delta.Y <<= 1
		delta.Y |= uint(changes[place.byteN]>>(7-byte(place.bitN))) & 1
		place.Increment()
	}
	for i := 0; i < 12; i++ {
		delta.X <<= 1
		delta.X |= uint(changes[place.byteN]>>(7-byte(place.bitN))) & 1
		place.Increment()
	}

	for i := 0; i < 3; i++ {
		delta.Color <<= 1
		delta.Color |= (changes[place.byteN] >> (7 - byte(place.bitN))) & 1
		place.Increment()
	}

	return delta
}

/*
Serialized Format:
Header:
|     24 bit      |
|Number of deltas |
For each delta
| 12 bit | 12 bit | 3 bit | -> 27 bit per delta
|    y   |    x   | color |

Color Values:
000 : Red - #FF0000
001 : Orange - #FFA500
010 : Yellow - #FFFF00
011 : Green - #008000
100 : Blue - #0000FF
101 : Purple - #800080
110 : Pink - #FFC0CB
111 : White - #FFFFFF
*/

func DeltaSerialize(changes []CanvasDelta) []byte {
	var output []byte = make([]byte, 3+(27*len(changes)+7)/8)
	var numChanges int = len(changes)
	output[0] = byte(numChanges >> 16)
	output[1] = byte(numChanges >> 8)
	output[2] = byte(numChanges)
	var place BitCursor = BitCursor{byteN: 3, bitN: 0}
	for i := 0; i < numChanges; i++ {
		PackDelta(output, &place, changes[i])
	}
	return output
}
func PackDelta(buff []byte, place *BitCursor, delta CanvasDelta) {
	//This seems quite inefficient. Perhaps can get the next x bits where
	// x is the min of bits left in byte and number of bits left in x/y/color in the delta
	for i := 0; i < 12; i++ {
		buff[place.byteN] = buff[place.byteN]<<1 | (byte(delta.Y>>(11-i)) & 1)
		place.Increment()
	}
	for i := 0; i < 12; i++ {
		buff[place.byteN] = buff[place.byteN]<<1 | (byte(delta.X>>(11-i)) & 1)
		place.Increment()
	}
	for i := 0; i < 3; i++ {
		buff[place.byteN] = buff[place.byteN]<<1 | (byte(delta.Color>>(2-i)) & 1)
		place.Increment()
	}
	buff[len(buff)-1] <<= 8 - byte(place.bitN)

}

package api

type CanvasDelta struct {
	x     uint
	y     uint
	color byte
}
type Canvas struct {
	width  uint
	height uint
	pixels []byte
}
type BitCursor struct {
	byteN int
	bitN  int
}

func (place *BitCursor) Increment() {
	place.bitN += 1
	if place.bitN == 8 {
		place.byteN += 1
		place.bitN = 0
	}
}

/*
Serialized Format:
| 12 bit | 12 bit | 3 bit | <- multiply color by the number of pixels.
| Height | Width  | color |

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
func Deserialize(data []byte) Canvas {
	var height uint = (uint(data[0]) << 4) | uint(data[1]>>4)
	var width uint = (uint(data[1]&0x0F) << 8) | uint(data[2])
	var pixels []byte = make([]byte, height*width)

	var bitCount uint = width * height * 3

	for i := uint(0); i < bitCount; i += 3 {
		pixels[i/3] = ((data[(i+24)/8]>>(i%8))&1)<<2 |
			((data[(i+24+1)/8]>>((i+1)%8))&1)<<1 |
			(data[(i+24+2)/8]>>((i+2)%8))&1
	}

	return Canvas{width: width, height: height, pixels: pixels}
}

func Serialize(canvas Canvas) []byte {
	var bitsNeeded = canvas.height*canvas.width*3 + 24
	var output []byte = make([]byte, (bitsNeeded+7)/8)
	output[0] = byte(canvas.height >> 4)
	output[1] = byte(canvas.height&0b00001111)<<4 | byte(canvas.width>>4)
	output[2] = byte(canvas.width&0b00001111) << 4
	//for each pixel add 3 bit color code
	var place BitCursor = BitCursor{byteN: 3, bitN: 0}
	for i := 0; i < len(canvas.pixels); i++ {
		output[place.byteN] = (output[place.byteN] << 1) | (canvas.pixels[i] >> 2)
		place.Increment()
		output[place.byteN] = (output[place.byteN] << 1) | (canvas.pixels[i] >> 1)
		place.Increment()
		output[place.byteN] = (output[place.byteN] << 1) | (canvas.pixels[i])
		place.Increment()
	}
	return output
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
		buff[place.byteN] = buff[place.byteN]<<1 | byte(delta.y>>(11-i))
		place.Increment()
	}
	for i := 0; i < 12; i++ {
		buff[place.byteN] = buff[place.byteN]<<1 | byte(delta.x>>(11-i))
		place.Increment()
	}
	for i := 0; i < 3; i++ {
		buff[place.byteN] = buff[place.byteN]<<1 | byte(delta.color>>(2-i))
		place.Increment()
	}

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
		delta.y <<= 1
		delta.y |= uint(changes[place.byteN] >> (7 - byte(place.bitN)))
		place.Increment()
	}
	for i := 0; i < 12; i++ {
		delta.x <<= 1
		delta.x |= uint(changes[place.byteN] >> (7 - byte(place.bitN)))
		place.Increment()
	}
	for i := 0; i < 3; i++ {
		delta.color <<= 1
		delta.color |= changes[place.byteN]>>7 - byte(place.bitN)
		place.Increment()
	}
	return delta
}

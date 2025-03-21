package api

type Canvas struct {
	width  uint
	height uint
	pixels []byte
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
	var byteN = 3 // start at 3 to account for size header
	var bitN = 0
	for i := 0; i < len(canvas.pixels); i++ {
		output[byteN] = (output[byteN] << 1) | (canvas.pixels[i] >> 2)
		Increment(&byteN, &bitN)
		output[byteN] = (output[byteN] << 1) | (canvas.pixels[i] >> 1)
		Increment(&byteN, &bitN)
		output[byteN] = (output[byteN] << 1) | (canvas.pixels[i])
		Increment(&byteN, &bitN)
	}
	return output
}

func Increment(byteN *int, bitN *int) {
	*bitN += 1
	if *bitN == 8 {
		*byteN += 1
		*bitN = 0
	}
}

type CanvasDelta struct {
	x     uint
	y     uint
	color byte
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
	var byteN int = 3
	var bitN int = 0
	for i := 0; i < numChanges; i++ {
		byteN, bitN = PackDelta(&output, byteN, bitN, changes[i])
	}
	return output
}

func PackDelta(buff *[]byte, byteN int, bitN int, delta CanvasDelta) (int, int) {
	//This seems quite inefficient. Perhaps can get the next x bits where
	// x is the min of bits left in byte and number of bits left in x/y/color in the delta
	for i := 0; i < 12; i++ {
		(*buff)[byteN] = (*buff)[byteN]<<1 | byte(delta.y>>(11-i))
		Increment(&byteN, &bitN)
	}
	for i := 0; i < 12; i++ {
		(*buff)[byteN] = (*buff)[byteN]<<1 | byte(delta.x>>(11-i))
		Increment(&byteN, &bitN)
	}
	for i := 0; i < 3; i++ {
		(*buff)[byteN] = (*buff)[byteN]<<1 | byte(delta.color>>(2-i))
		Increment(&byteN, &bitN)
	}
	return byteN, bitN
}

func DeltaDeserialize(changes []byte) []CanvasDelta {
	var numChanges int = int(changes[0])<<16 | int(changes[1])<<8 | int(changes[2])
	var output []CanvasDelta = make([]CanvasDelta, numChanges)
	var byteN int = 3
	var bitN int = 0
	for i := 0; i < numChanges; i++ {
		output[i] = deltaDeseralizeSingle(changes, &byteN, &bitN)
	}
	return output
}
func deltaDeseralizeSingle(changes []byte, byteN *int, bitN *int) CanvasDelta {
	var delta = CanvasDelta{}
	for i := 0; i < 12; i++ {
		delta.y <<= 1
		delta.y |= uint(changes[*byteN]>>7 - byte(*bitN))
		Increment(byteN, bitN)
	}
	for i := 0; i < 12; i++ {
		delta.x <<= 1
		delta.x |= uint(changes[*byteN]>>7 - byte(*bitN))
		Increment(byteN, bitN)
	}
	for i := 0; i < 3; i++ {
		delta.color <<= 1
		delta.color |= changes[*byteN]>>7 - byte(*bitN)
		Increment(byteN, bitN)
	}
}

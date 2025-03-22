package canvas

type Canvas struct {
	Width  uint
	Height uint
	Pixels []byte
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
	//fmt.Println("width: " + fmt.Sprint(width) + " " + "height:  " + fmt.Sprint(height))
	var bitCount uint = width * height * 3

	for i := uint(0); i < bitCount; i += 3 {
		pixels[i/3] = ((data[(i+24)/8]>>(i%8))&1)<<2 |
			((data[(i+24+1)/8]>>((i+1)%8))&1)<<1 |
			(data[(i+24+2)/8]>>((i+2)%8))&1
	}

	return Canvas{Width: width, Height: height, Pixels: pixels}
}

func Serialize(canvas Canvas) []byte {
	var bitsNeeded = canvas.Height*canvas.Width*3 + 24
	var output []byte = make([]byte, (bitsNeeded+7)/8)
	output[0] = byte(canvas.Height >> 4)
	output[1] = byte(canvas.Height&0b00001111)<<4 | byte(canvas.Width>>8)
	output[2] = byte(canvas.Width)
	//for each pixel add 3 bit color code
	var place BitCursor = BitCursor{byteN: 3, bitN: 0}
	for i := 0; i < len(canvas.Pixels); i++ {
		output[place.byteN] = (output[place.byteN] << 1) | (canvas.Pixels[i] >> 2)
		place.Increment()
		output[place.byteN] = (output[place.byteN] << 1) | (canvas.Pixels[i] >> 1)
		place.Increment()
		output[place.byteN] = (output[place.byteN] << 1) | (canvas.Pixels[i])
		place.Increment()
	}
	return output
}

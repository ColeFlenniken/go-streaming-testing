function deserialize(data) {
    const height = (data[0] << 4) | (data[1] >> 4);
    const width = ((data[1] & 0x0F) << 8) | data[2];
    const pixels = new Array(height * width);
    let byteN = 3;
    let bitN = 0;

    for (let i = 0; i < pixels.length; i++) {
        pixels[i] = ((data[byteN] >> (7 - bitN)) & 1) << 2;
        bitN++;
        if (bitN > 7) {
            bitN = 0;
            byteN++;
        }

        pixels[i] |= ((data[byteN] >> (7 - bitN)) & 1) << 1;
        bitN++;
        if (bitN > 7) {
            bitN = 0;
            byteN++;
        }

        pixels[i] |= (data[byteN] >> (7 - bitN)) & 1;
        bitN++;
        if (bitN > 7) {
            bitN = 0;
            byteN++;
        }
    }
    alert(pixels)
    return { Width: width, Height: height, Pixels: pixels };
}




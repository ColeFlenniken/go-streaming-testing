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


window.addEventListener('load', async ()=>{ 
    document.addEventListener('mousedown', startPainting); 
    document.addEventListener('mouseup', stopPainting); 
    document.addEventListener('mousemove', sketch); 
}); 

const c = document.getElementById("myCanvas");
const ctx = c.getContext("2d");
let coord = {x:0 , y:0};  


let paint = false; 
    
function getPosition(event){ 
    coord.x = event.clientX - c.offsetLeft; 
    coord.y = event.clientY - c.offsetTop; 
} 
    

function startPainting(event){ 
    paint = true; 
    getPosition(event); 
} 
function stopPainting(){ 
    paint = false; 
} 
    
function sketch(event){ 
    if (!paint) return; 
    ctx.beginPath(); 
    ctx.lineWidth = 2; 
    ctx.lineCap = 'round'; 
    ctx.strokeStyle = 'green'; 
    ctx.moveTo(coord.x, coord.y);
    getPosition(event); 
    ctx.lineTo(coord.x , coord.y); 
    ctx.stroke(); 
}



                    



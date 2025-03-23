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
    setCanvasData();    
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
//this will need to change
async function setCanvasData() {
    const url = "/getImageData";
    try{
    const response = await fetch(url)
    if(!response.ok){
        throw new Error("HTTP error " + response.status);

    }
    const json = await response.json();
    const imageData = ctx.getImageData(0, 0, c.width, c.height).data;
    
    let data = new Uint8ClampedArray(json.pixels)
    for (let i = 0; i < data.length; i++) {
        if(data[i] == 0 && imageData[i] != 0){
        data[i] = imageData[i];
        }

    }
    const imageData2 = new ImageData(data, c.width, c.height);
    ctx.putImageData(imageData2, 0, 0);
    }catch (error){
    console.log(error);
    }
    
}
//this will need to change 
function getCanvasData() {
    const imageData = ctx.getImageData(0, 0, c.width, c.height);
    let pixelarr = Array.from(imageData.data);
    fetch('/saveCanvas', {
    method: 'POST',
    headers: {
        'Content-Type': 'application/json',
    },
    body: JSON.stringify({
        pixels: pixelarr
    })
})
.then(response => {
    if (!response.ok) {
        throw new Error('Network response was not ok');
    }
    console.log('Canvas data saved successfully');
})
.catch(error => {
    console.error('Error saving canvas data:', error);
});
}
                    



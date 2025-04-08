/*
000 : Red - #FF0000
001 : Orange - #FFA500
010 : Yellow - #FFFF00
011 : Green - #008000
100 : Blue - #0000FF
101 : Purple - #800080
110 : Pink - #FFC0CB
111 : White - #FFFFFF
*/
const colors = [
    {  name: 'Red', color: '#FF0000' },
    {  name: 'Orange', color: '#FFA500' },
    {   name: 'Yellow', color: '#FFFF00' },
    {   name: 'Green', color: '#008000' },
    {  name: 'Blue', color: '#0000FF' },
    {   name: 'Purple', color: '#800080' },
    {  name: 'Black', color: '#000000' },
    {   name: 'White', color: '#FFFFFF' },
  ];
var colorBar = document.getElementById("colorBar");
var canvas = document.getElementById("myCanvas");
var ctx = canvas.getContext("2d");
ctx.globalAlpha = 1;
ctx.globalCompositeOperation = "destination-over";
var oldCanvas = ctx.getImageData(0, 0, canvas.width, canvas.height);

// Initialize default strokeStyle
ctx.strokeStyle = "black";

// Create color buttons dynamically
colors.forEach(color => {
    const button = document.createElement("div");
    button.classList.add("colorButton");
    button.style.backgroundColor = color.color;
    button.addEventListener("click", () => {
        ctx.fillStyle = color.color; // Update canvas strokeStyle
        console.log(`Stroke style changed to ${color}`);
    });
    colorBar.appendChild(button);
});

 
let changes = [];

 
var latestChangeId = -1;
// Set up variables to track mouse state
let isDrawing = false;



 
let lastX = null;
let lastY = null;

// Start drawing when mouse is pressed
canvas.addEventListener("mousedown", (event) => {
  isDrawing = true;
  const { x, y } = getMousePosition(event);
  lastX = x;
  lastY = y;
  drawPixel(x, y); // Start with a single pixel
});

// Stop drawing when mouse is released
canvas.addEventListener("mouseup", () => {
  isDrawing = false;
  lastX = null;
  lastY = null;
});

// Draw continuously as the mouse moves
canvas.addEventListener("mousemove", (event) => {
  if (isDrawing) {
    const { x, y } = getMousePosition(event);
    drawLine(lastX, lastY, x, y); // Draw a line between the last and current positions
    lastX = x;
    lastY = y; // Update the last position
  }
});

// Function to get mouse position relative to the canvas
function getMousePosition(event) {
  const rect = canvas.getBoundingClientRect();
  return {
    x: event.clientX - rect.left,
    y: event.clientY - rect.top,
  };
}

// Function to draw a pixel
function drawPixel(x, y) {
 
  ctx.fillRect(x, y, 1, 1); // Draw a 1x1 rectangle
}

// Function to draw a line between two points
function drawLine(x1, y1, x2, y2) {
 
  ctx.lineWidth = 2; // Set line width to 1 pixel
  ctx.beginPath();
  ctx.moveTo(x1, y1);
  ctx.lineTo(x2, y2);
  ctx.stroke();
}

window.addEventListener('load', async ()=>{
    // Function to start drawing
 
    setInterval(sendData, 2000);
    setInterval(updateData, 2000);
});


function sendData(){
    currData = ctx.getImageData(0, 0, canvas.width, canvas.height).data;
    deltas = [];
    for (let i = 0; i < currData.length; i += 4) {
        if(currData[i+3] == 0) continue;
        if(currData[i] != 0 || currData[i+1] != 0){
            console.log("FOUND NON WHITE " + currData[i] + " " + currData[i+1] + " " + currData[i+2] + " " + currData[i+3] );
        } 
        if(currData[i] != oldCanvas[i] || currData[i+1] != oldCanvas[i+1] ||currData[i+2] != oldCanvas[i+2] ||currData[i+3] != oldCanvas[i+3] ){
            let red = currData[i].toString(16).padStart(2, '0');
            let green =  currData[i+1].toString(16).padStart(2, '0');
            let blue = currData[i+2].toString(16).padStart(2, '0');
            
            let hexCode = `#${red}${green}${blue}`;
            console.log("code" + hexCode);
            let matchingColor = colors.find(colorObj => colorObj.color === hexCode);
            console.log("colro " + matchingColor.name);
            //TODO black is not a color on server side
            deltas.push({x:(i/4)%1000,y:(i/4)/1000,color:matchingColor.name})

        } else{

        }
    }
    console.log("sending " + deltas.length + " delts");
    if(deltas.length == 0) return;
    fetch("/update", {
        method: "POST",
        body: JSON.stringify(deltas)
      });
}
 
 

async function  updateData(){
    console.log("passed in" + latestChangeId)
    const response = await fetch("/getData",{
        method: "POST",
        body: latestChangeId
    });
    if(!response.ok){
        console.log("BADTHING");
        //TODO MAKE THIS A FULL REFRESH
        return;
    }
    const list = JSON.parse(await response.text());
    console.log("data is " + list);

    newCanvas = ctx.getImageData(0,0,canvas.width, canvas.height);
    list.forEach((item) =>{
        let ndx = item.y*1000 + item.x;
        hex = colors[item.color].color.replace(/^#/, '');
 
        // Parse the red, green, and blue components
        const r = parseInt(hex.substring(0, 2), 16);
        const g = parseInt(hex.substring(2, 4), 16);
        const b = parseInt(hex.substring(4, 6), 16);
        newCanvas.data[ndx] = r;
        newCanvas.data[ndx*4+1] = g;
        newCanvas.data[ndx*4+2] = b;
        console.log("color: rgb" + r + " " + g + " " + b);
        
    });
    ctx.putImageData(newCanvas,0,0);
    oldCanvas = ctx.getImageData(0, 0, canvas.width, canvas.height).data;
}

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
    {  name: 'Pink', color: '#FFC0CB' },
    {   name: 'White', color: '#FFFFFF' },
  ];
const colorBar = document.getElementById("colorBar");
const canvas = document.getElementById("myCanvas");
const ctx = canvas.getContext("2d");
var oldCanvas = ctx.getImageData(0, 0, canvas.width, canvas.height);
// Initialize default strokeStyle
ctx.strokeStyle = "black";

// Create color buttons dynamically
colors.forEach(color => {
    const button = document.createElement("div");
    button.classList.add("colorButton");
    button.style.backgroundColor = color.color;
    button.addEventListener("click", () => {
        ctx.strokeStyle = color.color; // Update canvas strokeStyle
        console.log(`Stroke style changed to ${color}`);
    });
    colorBar.appendChild(button);
});

let coord = {x:0 , y:0};
let changes = [];

let paint = false;

function getPosition(event){
    coord.x = event.clientX - canvas.offsetLeft;
    coord.y = event.clientY - canvas.offsetTop;
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
    ctx.moveTo(coord.x, coord.y);
    getPosition(event);
    ctx.lineTo(coord.x , coord.y);
    ctx.stroke();
}





window.addEventListener('load', async ()=>{
    document.addEventListener('mousedown', startPainting);
    document.addEventListener('mouseup', stopPainting);
    document.addEventListener('mousemove', sketch);
});


function sendData(){
    currData = ctx.getImageData(0, 0, canvas.width, canvas.height);
    deltas = [];
    for (let i = 0; i < currData.data.length; i += 4) {
        if(currData[i] != oldCanvas[i] || currData[i+1] != oldCanvas[i+1] ||currData[i+2] != oldCanvas[i+2] ||currData[i+3] != oldCanvas[i+3] ){
            let red = currData[i].toString(16).padStart(2, '0');
            let green =  currData[i+1].toString(16).padStart(2, '0');
            let blue = currData[i+2].toString(16).padStart(2, '0');
            
            let hexCode = `#${red}${green}${blue}`;
            let matchingColor = colors.find(colorObj => colorObj.color === hexCode).name;
            deltas.push({x:(i/4)%1000,y:(i/4)/1000,color:matchingColor})
        }    
    }
    fetch("/update", {
        method: "POST",
        body: JSON.stringify(deltas),
      });
}
 
 

function updateData(ctx){
    //http req to get the changes
    if(data.type == "delta"){

    }else{
        //recall the onstart grabcanvas function
    }
    //reset tracking
    oldCanvas = ctx.getImageData(0, 0, canvas.width, canvas.height);
}

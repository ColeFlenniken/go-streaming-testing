

const colors = ["red", "blue", "green", "yellow", "orange", "purple", "black", "pink"];
const colorBar = document.getElementById("colorBar");
const canvas = document.getElementById("myCanvas");
const ctx = canvas.getContext("2d");

// Initialize default strokeStyle
ctx.strokeStyle = "black";

// Create color buttons dynamically
colors.forEach(color => {
    const button = document.createElement("div");
    button.classList.add("colorButton");
    button.style.backgroundColor = color;
    button.addEventListener("click", () => {
        ctx.strokeStyle = color; // Update canvas strokeStyle
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
    
}
 
 

function updateData(ctx){
    //http req to get the changes
    if(data.type == "delta"){

    }else{
        //recall the onstart grabcanvas function
    }
}

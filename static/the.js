// Charts dealio
google.setOnLoadCallback(drawChart);       
function drawChart() {
    var data = google.visualization.arrayToDataTable(records["data"]); 

    var options = {
        title: 'Rowdy\'s followers and following',
        curveType: 'none',
        legend: { position: 'bottom' }
    };

    var chart = new google.visualization.LineChart(document.getElementById('chart'));

    chart.draw(data, options);
}


// Hashtags dealio
// Deal with form and hashtag
var form  = document.getElementById("form");
var adder = document.getElementById("empty");
var hashes = document.getElementById("hashes");
var link  = document.getElementById("auth").value;
var index = document.querySelectorAll(".hashtag").length;
var formListener = function(event){
    event.preventDefault();

    // Invalid? Get out of dodge
    if(document.querySelector(".invalid")){
        alert("Fix mistakes. Alpha numeric only")
        return;
    }

    var taglist = [];
    for(var hashtags = document.querySelectorAll(".hashtag"),i=0; i < hashtags.length; i=taglist.push(hashtags[i].value)){}
    document.location = link + "?hashtags=" + taglist.join("+")
}
form.addEventListener("submit", formListener, false);  //Modern browsers

// Set remove hashtags
var setRemove = function(remover){
    remover.addEventListener("click", function(event){
        document.querySelector(".hashtag" + this.dataset["index"]).remove();
        this.remove();
    });
}
for(var remove = document.querySelectorAll(".remove"),i=0; i < remove.length; i++){
    setRemove(remove[i]);
}

var setValidate = function(hashtag){
    hashtag.addEventListener("input", function(event){
        if(this.value.match(/[^A-Za-z0-9]/) || this.value == ""){
            this.classList.add("invalid");
        }else{
            this.classList.remove("invalid");
        }
    });
}

for(var hashtags = document.querySelectorAll(".hashtag"),i=0; i < hashtags.length; i++){
    setValidate(hashtags[i]);
}

// Add hashtags
adder.addEventListener("focus", function(event){
    var div = document.createElement("div");
    div.classList.add("hashgoup");

    var remover = document.createElement("div");
    remover.classList.add("remove");
    remover.dataset["index"] = index + 1
    remover.textContent = " x ";
    setRemove(remover)

    var input = document.createElement("input");                
    input.type = "text";
    input.classList.add("hashtag","invalid","hashtag" + (index + 1));
    setValidate(input);

    index++;
    div.appendChild(input);
    div.appendChild(remover);
    hashes.appendChild(div);

    input.focus();
});
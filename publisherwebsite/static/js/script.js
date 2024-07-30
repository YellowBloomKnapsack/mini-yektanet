/*
JSON data received from AdServer must have the following fields:
{
    "title"
    "image_link"
    "impression_link"
    "click_link"
}
----

*/
// var path = require("path")
// const publisherName = %publisherName%
// const AdServerAPILink = %AdServerAPILink%
const arr = window.location.href.split('/')
const publisherName = arr[arr.length-1]
const AdServerAPILink = "http://localhost:8081"
// const AdServerAPILink = "http://"+arr[arr.length-2].replace("8084","8081")
let publisherId = 1;
// hardcoded for the purposes of testing. the actual scriptTemplate.js does not work like this
switch (publisherName) {
    case "varzesh3": publisherId = 1; break;
    case "digikala": publisherId = 2; break;
    case "zoomit": publisherId = 3; break;
    case "sheypoor": publisherId = 4; break;
    case "filimo": publisherId = 5; break;
}

fetch(AdServerAPILink+"/"+publisherId)
.then((res) => {
    if (!res.ok) {
        throw new Error("unable to load from ad server.")
    }
    return res.json()
}).then((data) => {
    console.log(data)
    const adDiv = document.createElement("div")
    const img = document.createElement("img")
    img.src = data["image_link"]
    const title = document.createTextNode(data["title"])
    adDiv.appendChild(img)
    adDiv.appendChild(title)
    adDiv.onclick = ()=>clickHandler(data)
    adDiv.style="cursor:pointer;border:solid black 20px;"
    document.getElementsByTagName("body")[0].appendChild(adDiv)

    // handling when an element is in view
    let viewed = false
    let options = {
        root: null, // i.e. viewport
        rootMargin: "0px",
        threshold: 0.05,
    }

    const impressionHandler = (entries, observer) => {
        console.log(entries)
        entries.forEach((entry) => {
            console.log(entry)
            if(entry.isIntersecting && !viewed) {
                viewed = true
                fetch(data["impression_link"], {
                    method: "POST",
                    body: JSON.stringify({
                        token: data["impression_token"]
                    })
                })
            }
        })
        
    }

    let observer = new IntersectionObserver(impressionHandler, options)
    observer.observe(adDiv)
})

function clickHandler(data) {
    fetch(data["click_link"], {
        method: "POST",
        body: JSON.stringify({
            token: data["click_token"]
        })
    })
    .then(res=>{
        console.log(res)
        window.open(res.url)
    })
}
const publisherId = parseInt("%d")
const AdServerAPILink = "%s"

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
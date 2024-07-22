/*
{
    "ad_id"
	"image_link"
	"title"
	"impression_link"
	"click_link"
    "impression_token"
    "click_token"
}
----

*/
// const publisherID = %publisherID%
// const AdServerAPILink = %AdServerAPILink%
const publisherID = 1
const AdServerAPILink = "http://localhost:9001/here"

fetch(AdServerAPILink)
.then((res) => {
	if (!res.ok) {
		throw new Error("unable to load from ad server.")
	}
	return res.json()
}).then((data) => {
	const div = document.createElement("div")
	const img = document.createElement("img")
	img.src = data.image_link
	const title = document.createTextNode(data.title)
	div.appendChild(img)
	div.appendChild(title)
    div.onclick = ()=>clickHandler(data)
    div.style="cursor:pointer;border:solid black 20px;"
	document.getElementsByTagName("body")[0].appendChild(div)
    let viewed = false

    const impressionHandler = onVisibilityChange(div, function() {
        if(!viewed) {
            viewed = true
            fetch(AdServerAPILink, {
                method: "POST",
                body: JSON.stringify({
                    token: data["impression_token"],
                    pub_id: publisherID,
                    ad_id: data["ad_id"]
                })
            })
        }
    });

    addEventListener('DOMContentLoaded', handler, false);
    addEventListener('load', handler, false);
    addEventListener('scroll', handler, false);
    addEventListener('resize', handler, false);
})

function clickHandler(data) {
    fetch(data["click_link"], {
        method: "POST",
        body: JSON.stringify({
            token: data["click_token"],
            pub_id: publisherID,
            ad_id: data["ad_id"]
        })
    })
}

const isElementInViewport = (el, partiallyVisible = true) => {
    const { top, left, bottom, right } = el.getBoundingClientRect();
    const { innerHeight, innerWidth } = window;
    return partiallyVisible
    ? ((top > 0 && top < innerHeight) ||
        (bottom > 0 && bottom < innerHeight)) &&
        ((left > 0 && left < innerWidth) || (right > 0 && right < innerWidth))
    : top >= 0 && left >= 0 && bottom <= innerHeight && right <= innerWidth;
};

function onVisibilityChange(el, callback) {
    var old_visible = false;
    return function (e) {
        var visible = isElementInViewport(el);
        if (visible != old_visible) {
            old_visible = visible;
            callback()
        }
    }
}
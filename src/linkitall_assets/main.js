// This script is used by the generated HTML page containing the graph.
//
// Thanks to:
//   - https://anseki.github.io/leader-line/
//     Used for making connections.
// Note: These entities are not associated with the project.

let links = []

// just renaming the function to a simpler one
function id2el(idstr) {
    return document.getElementById(idstr)
}

// Find all the link-source dots and connect them to their target dot.
function connectDots() {
    // Template used to color links and their dots
    const colorTemplate = "hsl({hue}, 40%, 50%)"
    // we will cycle over different values of hue
    let hue = 0
    const linkSources = document.getElementsByClassName("link-source")

    for (let idx=0; idx < linkSources.length; idx++) {
        const source = linkSources[idx]
        const sourceId = source.id
        if (!sourceId.startsWith("D_")) {
            console.log(`ERROR: source ${source} does not start with D_`)
            continue
        }

        const targetId = sourceId.replace(/^D_/, "U_")
        const target = id2el(targetId)
        if (target == null) {
            console.log(`ERROR: no target dot with id ${targetId}`)
            continue
        }

        // Color for this link
        let color = colorTemplate.replace("{hue}", hue.toString())
        hue = (hue + 67) % 360
        source.style.backgroundColor = color
        target.style.backgroundColor = color

        // The main step - connect the dots
        let link = new LeaderLine(source, target)
        link.setOptions({startSocket: 'bottom', endSocket: 'top'})
        link.color = color
        link.size = 2
        links.push(link)
    }
}

function main() {
    connectDots()
}

document.addEventListener("DOMContentLoaded", main)

// When clicking the link box, focus and show the target node.
function showNode(nodeId) {
    const elem = id2el(nodeId)
    elem.scrollIntoView({behavior: "smooth", block: "center", inline: "center"})
    elem.classList.add("highlighted-node")

    // TODO: this looks so wrong. We need a *better* way to highlight a node.
    setTimeout(() => {
        elem.classList.remove("highlighted-node")
    }, 1500)
}

// Use state=true to enable link-view-panel and state=false to hide it.
function setLinkViewPatelState(state) {
    let linkViewOuterElem = id2el("link-view-panel")
    if (state) {
        linkViewOuterElem.style.zIndex = "3"
        linkViewOuterElem.style.display = "block"
        document.body.style.overflow = 'hidden'
    } else if (id2el("link-view-panel").style.display != "none") {
        linkViewOuterElem.style.zIndex = "-1"
        linkViewOuterElem.style.display = "none"
        document.body.style.overflow = 'visible'
    }
}

var currentLinkViewUrl = ""

// Open panel for viewing the target url.
// If we open the same link again, reuse the iframe.
function openNodeLink(evt, url, aux) {
    evt.preventDefault()

    // In case of middle-click or ctrl-click, open link in a new tab
    if (aux || (evt.ctrlKey == true)) {
        window.open(url, "newTab")
        return
    }

    if (url == currentLinkViewUrl) {
        setLinkViewPatelState(true)
        return
    }

    let inner = id2el("link-view-inner")
    inner.textContent = ""

    let height = window.innerHeight - 100
    let width = Math.floor(window.innerWidth * 0.8)
    if (width < 800) {
        width = 800
    }

    // Open link in an iframe and insert it into the inner div
    let iframe = document.createElement("iframe")
    iframe.src = url
    iframe.width = `${width}px`
    iframe.height = `${height}px`
    iframe.frameBorder="0"
    iframe.onload = () => {
        setLinkViewPatelState(true)
        iframe.focus()
        currentLinkViewUrl = url
    }
    inner.appendChild(iframe)
}

function closeLinkViewPanel() {
    setLinkViewPatelState(false)
}

// Handle escape keypress.
// Close the link view panel when pressing escape.
document.onkeydown = function(evt) {
    if(evt.key === "Escape") {
        closeLinkViewPanel()
    }
}



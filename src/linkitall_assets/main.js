// This script is used by the generated HTML page containing the graph.
//
// Thanks to:
//   - https://anseki.github.io/leader-line/
//     Used for making connections.
// Note: These entities are not associated with the project.

let links = []
let imgWidth = "60%"
let _buildConfig = null

// just renaming the function to a simpler one
function id2el(idstr) {
    return document.getElementById(idstr)
}

function getBuildConfig() {
    if (_buildConfig != null) {
        return _buildConfig
    }
    const elem = id2el("buildconfig")
    _buildConfig = JSON.parse(elem.innerHTML)
    return _buildConfig
}

function getBaseUrl(url) {
    return url.split(/[?#]/)[0]
}

function removeAndAddClass(elem, className) {
    elem.classList.remove(className)
    elem.classList.add(className)
}

function isImageFile(filename) {
  let checkExt = ['.jpg', '.jpeg', '.png', '.gif', '.bmp', '.svg']
  let lowerFilename = filename.toLowerCase()
  return checkExt.some(ext => lowerFilename.endsWith(ext))
}

function getLinkOptions(source, target, color) {
    let sourceTop = source.getBoundingClientRect().top
    let targetTop = target.getBoundingClientRect().top

    let startSocket = 'bottom'
    let endSocket = 'top'
    if (targetTop < sourceTop) {
        startSocket = 'top'
        endSocket = 'bottom'
    }

    let options = {
        startSocket,
        endSocket,
        color,
        size: 2
    }

    return options
}

// Find all the link-source dots and connect them to their target dot.
function connectDots() {
    // Template used to color links and their dots
    const colorTemplate = "hsl({hue}, 40%, 50%)"
    // we will cycle over different values of hue
    let hue = 0
    const linkSources = document.getElementsByClassName("link-source")

    for (let idx=0; idx < linkSources.length; idx++) {
        let source = linkSources[idx]
        const sourceId = source.id
        if (!sourceId.startsWith("D_")) {
            console.log(`ERROR: source ${source} does not start with D_`)
            continue
        }

        const targetId = sourceId.replace(/^D_/, "U_")
        let target = id2el(targetId)
        if (target == null) {
            console.log(`ERROR: no target dot with id ${targetId}`)
            continue
        }

        // Color for this link
        let color = colorTemplate.replace("{hue}", hue.toString())
        hue = (hue + 67) % 360
        source.style.backgroundColor = color
        target.style.backgroundColor = color

        let buildConfig = getBuildConfig()
        if (buildConfig.ArrowDirection == "parent2child") {
            // Swap source and target in this case
            let temp = source
            source = target
            target = temp
        }

        let link = new LeaderLine(source, target)
        link.setOptions(getLinkOptions(source, target, color))
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

function getWidthAndHeightForFrame() {
    let height = window.innerHeight - 100
    let width = Math.floor(window.innerWidth * 0.8)
    if (width < 800) {
        width = 800
    }
    return [width, height]
}

function withpx(val) {
    return `${val}px`
}

function openInIframe(url, targetParent) {
    let [width, height] = getWidthAndHeightForFrame()

    // Open link in an iframe and insert it into the inner div
    let iframe = document.createElement("iframe")
    iframe.src = url
    iframe.width = withpx(width)
    iframe.height = withpx(height)
    iframe.frameBorder="0"
    iframe.onload = () => {
        setLinkViewPatelState(true)
        iframe.focus()
    }
    targetParent.appendChild(iframe)
}

function openInDiv(url, targetParent) {
    let [width, height] = getWidthAndHeightForFrame()

    targetParent.innerHTML = `
        <div id="div-frame" class="div-frame">
            <img id="frame-img" src="${url}" width="${imgWidth}"/>
        </div>
    `
    let divFrame = id2el("div-frame")
    divFrame.style.width = withpx(width)
    divFrame.style.height = withpx(height)

    setLinkViewPatelState(true)
    removeAndAddClass(targetParent, "display-load-effect")
}

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

    // clear contents of the link-view-inner
    let inner = id2el("link-view-inner")
    inner.textContent = ""

    let baseUrl = getBaseUrl(url)
    if (isImageFile(baseUrl)) {
        // Images are open witout an iframe, with a simple div
        openInDiv(url, inner)
    } else {
        // Everything else opened with an iframe
        openInIframe(url, inner)
    }
    currentLinkViewUrl = url
}

function closeLinkViewPanel() {
    setLinkViewPatelState(false)
}


function zoomFrameImg(key) {
    let minWidth = 100
    let maxWidth = 4000
    let elem = id2el("frame-img")
    if (elem == null) {
        return
    }
    let width = elem.width
    if ((key === "[") && (width > minWidth)) {
        imgWidth = width * 0.9
        elem.width = imgWidth
    }
    if ((key === "]") && (width < maxWidth)) {
        imgWidth = width * 1 / 0.9
        elem.width = imgWidth
    }
}


// Handle escape keypress.
// Close the link view panel when pressing escape.
document.onkeydown = function(evt) {
    if(evt.key === "Escape") {
        closeLinkViewPanel()
    }

    if((evt.key === "[") || (evt.key === "]")) {
        zoomFrameImg(evt.key)
    }
}

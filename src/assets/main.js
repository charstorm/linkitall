// This script is used by the generate HTML page containing the graph.
//
// Thanks to:
//   - https://anseki.github.io/leader-line/
//     Used for making connections.
// Note: These entities are not associated with the project.

let links = []

// Find all the link-source dots and connect them to their target dot.
function connectDots() {
    // Template used to color links and their dots
    const colorTemplate = "hsl({hue}, 50%, 50%)"
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
        const target = document.getElementById(targetId)
        if (target == null) {
            console.log(`ERROR: no target dot with id ${targetId}`)
            continue
        }

        // Color for this link
        let color = colorTemplate.replace("{hue}", hue.toString())
        hue = (hue + 43) % 360
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

// Focus on the target node
function showNode(nodeId) {
    const elem = document.getElementById(nodeId)
    elem.scrollIntoView({behavior: "smooth", block: "center", inline: "center"})
    elem.classList.add("highlighted-node")

    // TODO: this looks so wrong. We need a *better* way to highligh a node.
    setTimeout(() => {
        elem.classList.remove("highlighted-node")
    }, 1000)
}

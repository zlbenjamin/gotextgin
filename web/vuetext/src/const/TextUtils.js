
// export
// Process content of a text
export function processTextContent(content) {
    if (Object.is(content, null) || Object.is(content, undefined)) {
        return ''
    }

    let pcontent = content
    pcontent = htmlEncode(pcontent)

    // replace \n
    pcontent = pcontent.replaceAll("\n", "<br>")

    // replace url in content
    // Note, there are some bugs TODO
    let urlg = new RegExp("http[s]://[a-zA-Z0-9\.\?/#=_:\\-\\+%]+", "ig")
    pcontent = pcontent.replaceAll(urlg, replaceUrl)

    return pcontent
}

// export
// Process comment of text
export function processComment(content) {
    if (Object.is(content, null) || Object.is(content, undefined)) {
        return ''
    }

    let pcontent = content
    pcontent = htmlEncode(pcontent)

    // replace \n
    pcontent = pcontent.replaceAll("\n", "<br>")

    // replace url in content
    // Note, there are some bugs TODO
    let urlg = new RegExp("http[s]://[a-zA-Z0-9\.\?/#=_:\\-\\+%]+", "ig")
    pcontent = pcontent.replaceAll(urlg, replaceUrl)

    return pcontent
}

// HTML special characters
function htmlEncode(content) {
    content = content.replaceAll("&", "&amp;")
			
    content = content.replaceAll("<", "&lt;")
    content = content.replaceAll(">", "&gt;")
    content = content.replaceAll("\"", "&quot;")
    content = content.replaceAll(" ", "&nbsp;")

    content = content.replaceAll("\t", "&emsp;")

    return content
}

// Replace matched url with <a> element
function replaceUrl(url) {
    return "<a href='" + url + "' target='_blank' class='content-url' >" + decodeURI(url) + "</a>"
}


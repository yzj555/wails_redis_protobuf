<template>
    <pre id="json-renderer" class="showJson json-highlight" contenteditable="true"></pre>
</template>
<script>
export default {
    name: 'JsonView',
    props: {
        data: {
            type: Object,
        },
        textFunc: {
            type: Function,
        },
    },
    data() {
        return {
            viewer: new JSONViewer({ clickableUrls: false })
        }
    },
    // watch: {
    //     data(newVal, oldVal) {
    //         // this.data = newVal
    //         this.jsonRender(newVal)
    //     },
    // },
    mounted() {
        //初始化
        this.viewer = new JSONViewer({ clickableUrls: false });
        this.jsonRender(this.data)
    },
    methods: {
        jsonRender(data) {
            try {
                // const viewer = new JSONViewer({ clickableUrls: false });
                this.viewer.render(data);
                console.log(this.getText())
            } catch (error) {
                console.log('Wrong json format.', error)
            }
        },
        getText() {
            return this.viewer.getText()
        },

    }
}

export class JSONViewer {
    constructor(options = {}) {
        this.options = Object.assign({ rootCollapsable: true, clickableUrls: true, bigNumbers: false }, options);
    }

    isCollapsable(arg) {
        return arg instanceof Object && Object.keys(arg).length > 0;
    }

    isUrl(string) {
        const protocols = ['http', 'https', 'ftp', 'ftps'];
        return protocols.some(protocol => string.startsWith(protocol + '://'));
    }

    htmlEscape(s) { return s.replace(/&/g, '&amp;').replace(/</g, '&lt;').replace(/>/g, '&gt;').replace(/'/g, '&apos;').replace(/"/g, '&quot;'); }

    json2html(json) {
        let html = '';
        if (typeof json === 'string') {
            json = this.htmlEscape(json);
            if (this.options.clickableUrls && this.isUrl(json)) {
                html += `<a href="${json}" class="json-string" target="_blank">${json}</a>`;
            } else {
                json = json.replace(/&quot;/g, '\\&quot;');
                html += `<span class="json-string">"${json}"</span>`;
            }
        } else if (typeof json === 'number' || typeof json === 'bigint') {
            html += `<span class="json-literal">${json}</span>`;
        } else if (typeof json === 'boolean') {
            html += `<span class="json-literal">${json}</span>`;
        } else if (json === null) {
            html += '<span class="json-literal">null</span>';
        } else if (Array.isArray(json)) {
            if (json.length > 0) {
                html += '[<ol class="json-array">';
                for (let i = 0; i < json.length; ++i) {
                    html += '<li>';
                    if (this.isCollapsable(json[i])) {
                        html += '<a class="json-toggle"></a>';
                    }
                    html += this.json2html(json[i]);
                    if (i < json.length - 1) {
                        html += ',';
                    }
                    html += '</li>';
                }
                html += '</ol>]';
            } else {
                html += '[]';
            }
        } else if (typeof json === 'object') {
            if (this.options.bigNumbers && (typeof json.toExponential === 'function' || json.isLosslessNumber)) {
                html += `<span class="json-literal">${json.toString()}</span>`;
            } else {
                const keyCount = Object.keys(json).length;
                if (keyCount > 0) {
                    html += '{<ul class="json-dict">';
                    let count = 0;
                    for (const key in json) {
                        if (Object.prototype.hasOwnProperty.call(json, key)) {
                            const jsonElement = json[key];
                            const escapedKey = this.htmlEscape(key);
                            const keyRepr = `<span class="json-string-key">"${escapedKey}"</span>`;
                            html += '<li>';
                            if (this.isCollapsable(jsonElement)) {
                                html += `<a class="json-toggle">${keyRepr}</a>`;
                            } else {
                                html += keyRepr;
                            }
                            html += ': ' + this.json2html(jsonElement);
                            if (++count < keyCount) {
                                html += ',';
                            }
                            html += '</li>';
                        }
                    }
                    html += '</ul>}';
                } else {
                    html += '{}';
                }
            }
        }
        return html;
    }
    getText() {
        const rootElement = document.getElementById('json-renderer');
        return rootElement.innerText
    }
    render(jsonData) {
        const html = this.json2html(jsonData);
        const rootElement = document.getElementById('json-renderer');

        if (this.options.rootCollapsable && this.isCollapsable(jsonData)) {
            rootElement.innerHTML = '<a class="json-toggle"></a>' + html;
        } else {
            rootElement.innerHTML = html;
        }
        const toggleChildren = (element, collapse) => {
            const childToggles = element.querySelectorAll('.json-toggle');
            childToggles.forEach(toggle => {
                const container = toggle.nextElementSibling;
                if (container && (container.classList.contains('json-dict') || container.classList.contains('json-array'))) {
                    const isCollapsed = toggle.classList.contains('collapsed');
                    if (collapse !== isCollapsed) {
                        toggle.click();
                    }
                }
            });
        };
        document.addEventListener('click', (e) => {
            if (e.target.matches('.json-toggle')) {
                e.preventDefault();
                const target = e.target;
                const listItem = target.closest('li') || target.parentElement;
                target.classList.toggle('collapsed');
                const container = target.nextElementSibling;
                if (container && (container.classList.contains('json-dict') || container.classList.contains('json-array'))) {
                    container.classList.toggle('hidden');
                    if (e.ctrlKey) {
                        toggleChildren(container, target.classList.contains('collapsed'));
                    }
                    const siblings = container.parentNode.children;
                    for (let i = Array.from(siblings).indexOf(container) + 1; i < siblings.length; i++) {
                        if (siblings[i].classList.contains('json-placeholder')) {
                            siblings[i].remove();
                            i--;
                        }
                    }
                    if (container.classList.contains('hidden')) {
                        const count = container.children.length;
                        let placeholder = document.createElement('a');
                        placeholder.className = 'json-placeholder';
                        placeholder.textContent = `${count} ${count > 1 ? 'items' : 'item'}`;
                        if (!container.nextElementSibling?.classList.contains('json-placeholder')) {
                            container.parentNode.insertBefore(placeholder, container.nextSibling);
                        }
                    }
                }
            } else if (e.target.matches('.json-placeholder')) {
                e.preventDefault();
                const toggle = e.target.previousElementSibling.previousElementSibling;
                if (toggle) {
                    toggle.click();
                }
            }
        });
    }
}

</script>

<style>

.showJson{
    background-color: rgba(0, 0, 0, 0.7);
    height: 80%;
    overflow: auto;
    color: aliceblue;
}

pre {
    white-space: pre-wrap;
    word-wrap: break-word;
    overflow-wrap: break-word;
    display: block;
}

#json-renderer {
    padding: 1em 2em;
}

ul.json-dict,
ol.json-array {
    list-style-type: none;
    margin: 0 0 0 1px;
    border-left: 1px dotted #666;
    padding-left: 2em;
}

.json-string-key {
    color: #7bdcfe;
}

.json-string {
    color: #ce9178;
}

.json-literal {
    color: #b5cea8;
    font-weight: bold;
}

a.json-toggle {
    position: relative;
    color: inherit;
    text-decoration: none;
    cursor: pointer;
}

a.json-toggle:focus {
    outline: none;
}

a.json-toggle:before {
    font-size: 0.8em;
    color: #666;
    content: "\25BC";
    position: absolute;
    display: inline-block;
    width: 1em;
    text-align: center;
    line-height: 1em;
    left: -1.3em;
    top: 1px;
}

a.json-toggle:hover:before {
    color: #aaa;
}

a.json-toggle.collapsed:before {
    content: "\25B6";
}

a.json-placeholder {
    color: #aaa;
    padding: 0 1em;
    text-decoration: none;
    cursor: pointer;
}

a.json-placeholder:hover {
    text-decoration: underline;
}

.hidden {
    display: none;
}
</style>
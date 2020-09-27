const http = require('https');
const fs = require('fs').promises;
const url = require('url');
const path = require('path');

async function getFileList(fromFileName) {
    const file = await fs.readFile(fromFileName);

    return file.toString().split('\n').filter(str => str);
}

function getBaseNameFromUrl(uri) {
    return path.basename(url.parse(uri).pathname);
}

async function fetchFile(filename) {
    return new Promise((resolve, reject) => {
        // console.log(`Begin downloading: ${filename}`);

        http.get(filename, res => {
            if (res.statusCode !== 200) {
                return reject(new Error(`Request ${filename} Failed. Status Code: ${res.statusCode}`));
            }

            let body = '';
            res.on('data', data => body += data);
            res.on('end', _ => resolve(body));
            res.on('error', reject);
        });
    })
}

module.exports.getFileList = getFileList;
module.exports.getBaseNameFromUrl = getBaseNameFromUrl;
module.exports.fetchFile = fetchFile;
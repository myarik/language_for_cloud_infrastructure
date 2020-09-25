const http = require('http');
const os = require('os');
const fs = require('fs').promises;
const BPromise = require('bluebird');

const STORAGE_MAPPING = {
    "storage0": "http://212.183.159.230/10MB.zip",
    "storage1": "http://ipv4.download.thinkbroadband.com/10MB.zip",
    "storage2": "http://speedtest.tele2.net/10MB.zip"
}

let cancel = false;

async function fetchStorage(storage) {
    const filename = STORAGE_MAPPING[storage];

    return new Promise((resolve, reject) => {
        console.log(`Begin downloading from storage: ${storage}`);

        let req = http.get(filename, res => {
            if (res.statusCode !== 200) {
                return reject(new Error(`Request ${storage} Failed. Status Code: ${res.statusCode}`));
            }

            let body = '';
            res.on('data', data => {
                body += data;
                if (cancel) {
                    req.abort();
                    reject(new Error('Aborted'));
                }
            });
            res.on('end', _ => resolve(body));
            res.on('error', reject);
        });

        req.on('error', reject);
    })
}

(async _ => {
    const time = Date.now();

    const promises = Object.keys(STORAGE_MAPPING).map(async storage => {
        try {
            const body = await fetchStorage(storage)
            console.log(`${storage} returns the first result.`)
        } catch (e) {
            if (!cancel) {
                console.log(`[!] ${storage}. Task returns an error: ${e.message}`);
            }
            throw e;
        }
    });

    await BPromise.some(promises, 1);
    cancel = true;

    console.log(`Execution time: ${(Date.now() - time) / 1000} second(s)`);
})();
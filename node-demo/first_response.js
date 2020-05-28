const http = require('https');
const os = require('os');
const fs = require('fs').promises;

const HOST_URL = process.env['API_HOST_URL'];

const STORAGE_MAPPING = {
    "storage0": "0cf50f1c99234954b00340471538ce9d.MOV",
    "storage1": "0db9a58b669048dc999eb8f11f7ba424.MOV",
    "storage2": "0d38ceda70b14ccfaf6960514615757f.MOV"
}

let cancel = false;


async function fetchStorage(storage) {
    const filename = STORAGE_MAPPING[storage];

    return new Promise((resolve, reject) => {
        console.log(`Begin downloading from storage: ${storage}`);

        let req = http.get(HOST_URL + filename, res => {
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
    })
}

(async _ => {
    const time = Date.now();

    const promises = Object.keys(STORAGE_MAPPING).map(async storage => {
        try {
            await fetchStorage(storage)
            console.log(`${storage} returns the first result.`)
        } catch (e) {
            if (!cancel) {
                console.log(`Task returns an error: ${e.message}`);
            }
        }
    });

    //Getting the first response
    await Promise.race(promises);
    cancel = true;
    console.log(`Execution time: ${(Date.now() - time) / 1000} second(s)`);
})();
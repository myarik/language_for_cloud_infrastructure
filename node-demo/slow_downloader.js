const http = require('https');
const os = require('os');
const fs = require('fs').promises;
const qrate = require('qrate');

const HOST_URL = process.env['API_HOST_URL'];

const FILES = [
    "0cf50f1c99234954b00340471538ce9d.MOV",
    "0db9a58b669048dc999eb8f11f7ba424.MOV",
    "0d38ceda70b14ccfaf6960514615757f.MOV",
    "0CB55372-0173-49F7-9EAF-6CF1A40382C5.MOV",
    "0f132134b2474cbd858559ed979835a3.MOV"
];

async function fetchFile(filename) {
    return new Promise((resolve, reject) => {
        console.log(`Begin downloading: ${filename}`);

        http.get(HOST_URL + filename, res => {
            if (res.statusCode !== 200) {
                return reject(new Error(`Request ${filename} Failed. Status Code: ${res.statusCode}`));
            }

            let body = '';
            res.on('data', data => body += data);
            res.on('end', _ => {
                console.log(`File ${filename} downloaded`);
                return resolve(body);
            });
            res.on('error', reject);
        });
    })
}

let tasks = FILES.length;
const worker = async ({filename, body}) => {
    // to make storage slow (for 10 seconds)
    // await new Promise(resolve => { setTimeout(() => resolve(), 10000); });

    await fs.writeFile(filename, body);
    console.log(`Finished writing: ${filename}`)
    tasks--;
};

(async _ => {
    const time = Date.now();

    const q = qrate(worker, 3);
    q.drain = _ => {
        !tasks && console.log(`Execution time: ${(Date.now() - time) / 1000} second(s)`);
    };

    await Promise.all(FILES.map(name => {
        const filename = `${os.tmpdir()}/${name}`;
        try {
            fetchFile(name).then(
                body => q.push({filename, body})
            );
        } catch (e) {
            console.log(`Error download: ${e.message}`);
        }
    }));
})();
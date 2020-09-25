const os = require('os');
const fs = require('fs').promises;
const BPromise = require('bluebird');
const utils = require('./utils');

const CONCURRENCY_LEVEL = 3;

(async _ => {
    const time = Date.now();

    const fileNames = await utils.getFileList(process.env.CONTENT_FILE);

    await BPromise.map(fileNames, async url => {
        const baseName = utils.getBaseNameFromUrl(url);
        const filename = `${os.tmpdir()}/${baseName}`;

        try {
            const body = await utils.fetchFile(url)
            await fs.writeFile(filename, body);
            console.log(`Finished writing: ${filename}`)
        } catch (e) {
            console.log(`Error download: ${e.message}`);
        }

    }, {concurrency: CONCURRENCY_LEVEL});

    console.log(`Execution time: ${(Date.now() - time) / 1000} second(s)`);
})();
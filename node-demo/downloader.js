const os = require('os');
const fs = require('fs').promises;
const utils = require('./utils');

(async _ => {
    const time = Date.now();

    const fileNames = await utils.getFileList(process.env.CONTENT_FILE);

    await Promise.all(fileNames.map(async url => {
        const baseName = utils.getBaseNameFromUrl(url);
        const filename = `${os.tmpdir()}/${baseName}`;

        try {
            const body = await utils.fetchFile(url)
            await fs.writeFile(filename, body);
            // console.log(`Finished writing: ${filename}`)
        } catch (e) {
            console.log(`Error download: ${e.message}`);
        }
    }));

    console.log(`Execution time: ${(Date.now() - time) / 1000} second(s)`);
})();
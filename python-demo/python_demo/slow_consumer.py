import asyncio
import os
import sys
import tempfile
import time
import uuid
from typing import List

import aiofiles
import aiohttp
import click
import uvloop
from loguru import logger

logger.remove()
logger.add(sys.stderr, level="INFO")

TIMEOUT = aiohttp.ClientTimeout(total=10)


async def content_producer(
    source_file: str, queue: asyncio.Queue, connector: aiohttp.TCPConnector
) -> None:
    """
    Download a content
    """
    for url in await read_file(source_file):
        logger.debug(f"Begin downloading {url}")
        async with aiohttp.ClientSession(
            connector=connector, timeout=TIMEOUT, connector_owner=False
        ) as session:
            try:
                async with session.get(url) as resp:
                    body = await resp.read()
                    # put the item in the queue
                    await queue.put(body)
            except asyncio.TimeoutError as err:
                logger.error(
                    f"Cannot download a content, url: {url} error: {type(err).__name__}"
                )


async def content_consumer(worker_id: int, queue: asyncio.Queue, tmpdir: str) -> None:
    """
    Save a content to the localstorage
    """
    while True:
        content = await queue.get()
        # Add timeout to see how it works
        if os.environ.get("DEBUG", False):
            await asyncio.sleep(1 + (1 * worker_id))

        filename = os.path.join(tmpdir, f"async_{str(uuid.uuid4())}.mov")

        async with aiofiles.open(filename, "wb") as video_file:
            await video_file.write(content)
            logger.debug(f"[WORKER {worker_id}] Finished writing {filename}")

        queue.task_done()


async def read_file(source_file: str) -> List[str]:
    """
    Read urls
    """
    async with aiofiles.open(source_file, mode="r") as f:
        content = await f.read()
    return [url for url in content.split()]


@click.command()
def main() -> None:
    """
    The scraper gets data from sources and saves them to our local machine.
    This scraper has only three simultaneous consumers to prevent the storage overload

    To use this script, you should set the environment variable CONTENT_FILE.
    Then run `python -m python_demo.slow_consumer`.
    """
    asyncio.set_event_loop_policy(uvloop.EventLoopPolicy())
    asyncio.run(async_main())


async def async_main() -> None:
    s = time.perf_counter()

    is_debug_level = os.environ.get("DEBUG", False)
    if is_debug_level:
        logger.remove()
        logger.add(sys.stderr, level="DEBUG")

    source_file = os.environ.get("CONTENT_FILE", None)
    if not source_file:
        logger.error("cannot find the CONTENT_FILE environment variable")
        return

    queue: asyncio.Queue[bytes] = asyncio.Queue(maxsize=3)
    connector = aiohttp.TCPConnector(limit=5)

    with tempfile.TemporaryDirectory() as tmpdir:
        # Create three worker tasks to process the queue concurrently.
        tasks = []
        for index in range(3):
            task = asyncio.create_task(content_consumer(index, queue, tmpdir))
            tasks.append(task)

        # Wait until the tasks is fully processed
        await asyncio.gather(content_producer(source_file, queue, connector))
        # Wait until the queue is fully processed.
        await queue.join()

        # Cancel our worker tasks.
        for task in tasks:
            task.cancel()
        # Wait until all worker tasks are cancelled.
        await asyncio.gather(*tasks, return_exceptions=True)

    await connector.close()
    elapsed = time.perf_counter() - s
    logger.info(f"Execution time: {elapsed:0.2f} seconds.")


if __name__ == "__main__":
    main()

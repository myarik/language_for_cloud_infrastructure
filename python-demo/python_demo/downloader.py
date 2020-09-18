import asyncio
import os
import sys
import tempfile
import time
import uuid
from typing import Optional, List

import aiofiles
import aiohttp
import click
import uvloop
from aiohttp import TCPConnector
from attr import dataclass
from loguru import logger

logger.remove()
logger.add(sys.stderr, level="INFO")


@dataclass
class Result:
    content: Optional[bytes] = None
    error: Optional[Exception] = None


async def download_video(url: str, connector: TCPConnector) -> Result:
    """
    Download a content
    """
    logger.debug(f"Begin downloading {url}")
    timeout = aiohttp.ClientTimeout(total=10)
    async with aiohttp.ClientSession(
        connector=connector, timeout=timeout, connector_owner=False
    ) as session:
        try:
            async with session.get(url) as resp:
                return Result(content=await resp.read())
        except asyncio.TimeoutError as err:
            return Result(error=err)


async def write_to_file(tmpdirname: str, content: bytes) -> None:
    """
    Save a content to the localstorage
    """
    filename = os.path.join(tmpdirname, f"async_{str(uuid.uuid4())}.mov")
    async with aiofiles.open(filename, "wb") as video_file:
        await video_file.write(content)
        logger.debug(f"Finished writing {filename}")


async def web_scrape_task(url: str, tmpdirname: str, connector: TCPConnector) -> None:
    resp = await download_video(url, connector)
    if resp.error is None:
        await write_to_file(tmpdirname, resp.content)
    else:
        logger.error(
            f"Cannot download a content, url: {url} error: {type(resp.error).__name__}"
        )


@click.command()
def main() -> None:
    """
    This is the simple web scraper.
    The scraper gets data from sources and saves them to our local machine

    To use this script, you should set the environment variable API_HOST_URL.
    Then run `python -m python_demo.downloader`.
    """
    asyncio.set_event_loop_policy(uvloop.EventLoopPolicy())
    asyncio.run(async_main())


async def read_file(source_file: str) -> List[str]:
    async with aiofiles.open(source_file, mode="r") as f:
        content = await f.read()
    return [url for url in content.split()]


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

    conn = aiohttp.TCPConnector(limit=5)
    with tempfile.TemporaryDirectory() as tmpdirname:
        await asyncio.gather(
            *[
                web_scrape_task(url, tmpdirname, conn)
                for url in await read_file(source_file)
            ]
        )

    # Close connection
    await conn.close()
    elapsed = time.perf_counter() - s
    logger.info(f"Execution time: {elapsed:0.2f} seconds.")


if __name__ == "__main__":
    main()

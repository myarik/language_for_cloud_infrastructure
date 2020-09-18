import asyncio
import os
import sys
import time
from typing import Optional

import aiohttp
import click
import uvloop
from attr import dataclass
from loguru import logger

logger.remove()
logger.add(sys.stderr, level="INFO")

STORAGE_MAPPING = {
    "storage0": "http://212.183.159.230/10MB.zip",
    "storage1": "http://ipv4.download.thinkbroadband.com/10MB.zip",
    "storage2": "http://speedtest.tele2.net/10MB.zip",
}


@dataclass
class Result:
    storage: str
    content: Optional[bytes] = None
    error: Optional[Exception] = None


async def replica_storage(storage: str) -> Result:
    """
    Download a content
    """
    logger.debug(f"Begin downloading {storage}")
    async with aiohttp.ClientSession() as session:
        async with session.get(STORAGE_MAPPING[storage]) as resp:
            return Result(storage, content=await resp.read())


@click.command()
def main() -> None:
    """
    The service sends requests to multiple replicas, and use the first response

    To use this script, you should set the environment variable API_HOST_URL.
    Then run `python -m python_demo.first_response`.
    """
    asyncio.set_event_loop_policy(uvloop.EventLoopPolicy())
    asyncio.run(async_main())


async def async_main() -> None:
    """
    Send requests to multiple replicas, and use the first response.
    """
    s = time.perf_counter()

    is_debug_level = os.environ.get("DEBUG", False)
    if is_debug_level:
        logger.remove()
        logger.add(sys.stderr, level="DEBUG")

    tasks = [
        asyncio.create_task(replica_storage(f"storage{index}")) for index in range(3)
    ]
    done, pending = await asyncio.wait(
        tasks, timeout=5.0, return_when=asyncio.FIRST_COMPLETED
    )

    # The first response
    for task in done:
        try:
            result = task.result()
            logger.info(f"{result.storage.title()} returns the first result")
        except Exception as e:
            logger.error(f"Task returns an error {repr(e)}")
        break
    else:
        logger.error("Timeout error")

    # Cancel pending tasks
    for task in pending:
        task.cancel()
    # Wait until all pending tasks are cancelled.
    await asyncio.gather(*pending, return_exceptions=True)

    elapsed = time.perf_counter() - s
    logger.info(f"Execution time: {elapsed:0.2f} seconds.")


if __name__ == "__main__":
    main()

import asyncio
import os
import time
from typing import Optional

import aiohttp
import click
import uvloop
from attr import dataclass

HOST_URL = os.environ.get("API_HOST_URL")

STORAGE_MAPPING = {
    "storage0": "0cf50f1c99234954b00340471538ce9d.MOV",
    "storage1": "0db9a58b669048dc999eb8f11f7ba424.MOV",
    "storage2": "0d38ceda70b14ccfaf6960514615757f.MOV",
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
    click.secho(f"Begin downloading from the {storage}", fg="yellow")
    url = f"{HOST_URL}{STORAGE_MAPPING[storage]}"
    async with aiohttp.ClientSession() as session:
        async with session.get(url) as resp:
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
            click.secho(
                f"{result.storage.title()} returns the first result", fg="green"
            )
        except Exception as e:
            click.secho(f"Task returns an error {repr(e)}", fg="red")
        break
    else:
        click.secho(f"Timeout error", fg="red")

    # Cancel pending tasks
    for task in pending:
        task.cancel()
    # Wait until all pending tasks are cancelled.
    await asyncio.gather(*pending, return_exceptions=True)

    elapsed = time.perf_counter() - s
    click.secho(f"Execution time: {elapsed:0.2f} seconds.", fg="bright_blue")


if __name__ == "__main__":
    main()

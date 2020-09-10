import asyncio
import os
import tempfile
import time
import uuid

import aiohttp
import click
import uvloop

HOST_URL = os.environ.get("API_HOST_URL")

FILES = [
    "0cf50f1c99234954b00340471538ce9d.MOV",
    "0db9a58b669048dc999eb8f11f7ba424.MOV",
    "0d38ceda70b14ccfaf6960514615757f.MOV",
    "0CB55372-0173-49F7-9EAF-6CF1A40382C5.MOV",
    "0f132134b2474cbd858559ed979835a3.MOV",
]

TIMEOUT = aiohttp.ClientTimeout(total=10)


async def download_video(
    file_name: str, queue: asyncio.Queue, connector: aiohttp.TCPConnector
) -> None:
    """
    Download a content
    """
    click.secho(f"Begin downloading {file_name}", fg="yellow")
    url = f"{HOST_URL}{file_name}"
    async with aiohttp.ClientSession(
        connector=connector, timeout=TIMEOUT, connector_owner=False
    ) as session:
        try:
            async with session.get(url) as resp:
                body = await resp.read()
                # put the item in the queue
                await queue.put(body)
                click.secho(f"File {file_name} downloaded", fg="yellow")
        except asyncio.TimeoutError:
            click.secho(f"Error download a {file_name}", fg="red")


async def write_to_file(queue: asyncio.Queue, tmpdir: str) -> None:
    """
    Save a content to the localstorage
    """
    while True:
        content = await queue.get()
        filename = os.path.join(tmpdir, f"async_{str(uuid.uuid4())}.mov")
        with open(filename, "wb") as video_file:
            video_file.write(content)
            click.secho(f"Finished writing {filename}", fg="green")
        queue.task_done()


@click.command()
def main() -> None:
    """
    The scraper gets data from sources and saves them to our local machine.
    This scraper has only three simultaneous consumers to prevent the storage overload

    To use this script, you should set the environment variable API_HOST_URL.
    Then run `python -m python_demo.slow_consumer`.
    """
    asyncio.set_event_loop_policy(uvloop.EventLoopPolicy())
    asyncio.run(async_main())


async def async_main() -> None:
    s = time.perf_counter()
    queue: asyncio.Queue[bytes] = asyncio.Queue(maxsize=100)
    connector = aiohttp.TCPConnector(limit=5)

    with tempfile.TemporaryDirectory() as tmpdir:
        # Create three worker tasks to process the queue concurrently.
        tasks = []
        for _ in range(3):
            task = asyncio.create_task(write_to_file(queue, tmpdir))
            tasks.append(task)

        # Wait until the tasks is fully processed.s
        await asyncio.gather(
            *[download_video(file_name, queue, connector) for file_name in FILES]
        )
        # Wait until the queue is fully processed.
        await queue.join()

        # Cancel our worker tasks.
        for task in tasks:
            task.cancel()
        # Wait until all worker tasks are cancelled.
        await asyncio.gather(*tasks, return_exceptions=True)

    await connector.close()
    elapsed = time.perf_counter() - s
    click.secho(f"Execution time: {elapsed:0.2f} seconds.", fg="bright_blue")


if __name__ == "__main__":
    main()

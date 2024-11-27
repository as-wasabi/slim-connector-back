from typing_extensions import Coroutine


async def await_all(*tasks: Coroutine):
    for task in tasks:
        await task

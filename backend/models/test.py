import anyio
from rich import get_console


console = get_console()


async def run_script() -> None:
    """Load data from a fixture."""

    async with engine.begin() as conn:
        await conn.run_sync(UUIDBase.metadata.create_all)
    # 1) create a new Author record.
    console.print("1) Adding a new record")
    author = await create_author()
    author_id = author.id
    # 2) Let's update the Author record.
    console.print("2) Updating a record.")
    author.dod = datetime.strptime("1940-12-21", "%Y-%m-%d").date()
    await update_author(author)
    # 3) Let's delete the record we just created.
    console.print("3) Removing a record.")
    await remove_author(author_id)
    # 4) Let's verify the record no longer exists.
    console.print("4) Select one or none.")
    _should_be_none = await get_author_if_exists(author_id)


if __name__ == "__main__":

    anyio.run(run_script)

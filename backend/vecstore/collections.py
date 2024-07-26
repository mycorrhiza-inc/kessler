
from models import utils


async def ensure_collections_table():
    async with utils.sqlalchemy_config.get_engine().begin() as conn:
        statement = """
                CREATE TABLE IF NOT EXISTS milvus_collections (
                        id UUID PRIMARY KEY,
                        name VARCHAR(255) NOT NULL
                );
                """
        await conn.execute(statement)
        await conn.commit()


async def insert_collection(collection_name):
    async with utils.sqlalchemy_config.get_engine().begin() as conn:
        select_statement = """
                SELECT id FROM milvus_collections WHERE name = :collection_name;
                """
        result = await conn.execute(select_statement, collection_name=collection_name)
        existing_collection = await result.fetchone()

        if existing_collection is None:
            insert_statement = """
                INSERT INTO milvus_collections (id, name)
                VALUES (uuid_generate_v4(), :collection_name);
                """
            await conn.execute(insert_statement, collection_name=collection_name)
            await conn.commit()


async def ensure_file_milvus_collection_table():
    async with utils.sqlalchemy_config.get_engine().begin() as conn:
        statement = """
                CREATE TABLE IF NOT EXISTS file_milvus_collection (
                        file_id UUID REFERENCES files(id),
                        collection_id UUID REFERENCES collections(id),
                        PRIMARY KEY (file_id, collection_id)
                );
                """
        await conn.execute(statement)
        await conn.commit()
